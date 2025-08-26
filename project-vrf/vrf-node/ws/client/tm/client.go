package tm

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/metrics"

	"github.com/gorilla/websocket"
	"github.com/tendermint/tendermint/libs/service"
	tmsync "github.com/tendermint/tendermint/libs/sync"
	"github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

const (
	defaultMaxReconnectAttempts = 25
	defaultWriteWait            = 0
	defaultReadWait             = 0
	defaultPingPeriod           = 0
	defaultMaxSleepTime         = 3600 * time.Second
)

type WSClient struct { // nolint: maligned
	conn                 *websocket.Conn
	Address              string // IP:PORT or /path/to/socket
	Endpoint             string // /websocket/url/endpoint
	Dialer               func(string, string) (net.Conn, error)
	PubKey               string
	PriKey               *ecdsa.PrivateKey
	RequestsCh           chan types.RPCRequest
	onReconnect          func()
	send                 chan types.RPCResponse // user requests
	backlog              chan types.RPCResponse // stores a single user request received during a conn failure
	reconnectAfter       chan error             // reconnect requests
	readRoutineQuit      chan struct{}          // a way for readRoutine to close writeRoutine
	maxReconnectAttempts int
	protocol             string
	wg                   sync.WaitGroup
	mtx                  tmsync.RWMutex
	sentLastPingAt       time.Time
	reconnecting         bool
	nextReqID            int
	writeWait            time.Duration
	readWait             time.Duration
	pingPeriod           time.Duration
	service.BaseService
	PingPongLatencyTimer *metrics.Timer
}

func NewWS(remoteAddr, endpoint string, options ...func(*WSClient)) (*WSClient, error) {
	parsedURL, err := newParsedURL(remoteAddr)
	if err != nil {
		return nil, err
	}
	if parsedURL.Scheme != protoWSS {
		parsedURL.Scheme = protoWS
	}

	dialFn, err := makeHTTPDialer(remoteAddr)
	if err != nil {
		return nil, err
	}

	c := &WSClient{
		Address:              parsedURL.GetTrimmedHostWithPath(),
		Dialer:               dialFn,
		Endpoint:             endpoint,
		PingPongLatencyTimer: metrics.NewTimer(),

		maxReconnectAttempts: defaultMaxReconnectAttempts,
		readWait:             defaultReadWait,
		writeWait:            defaultWriteWait,
		pingPeriod:           defaultPingPeriod,
		protocol:             parsedURL.Scheme,
	}
	c.BaseService = *service.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c, nil
}

func MaxReconnectAttempts(max int) func(*WSClient) {
	return func(c *WSClient) {
		c.maxReconnectAttempts = max
	}
}

func ReadWait(readWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.readWait = readWait
	}
}

func WriteWait(writeWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.writeWait = writeWait
	}
}

func PingPeriod(pingPeriod time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.pingPeriod = pingPeriod
	}
}

func OnReconnect(cb func()) func(*WSClient) {
	return func(c *WSClient) {
		c.onReconnect = cb
	}
}

func (c *WSClient) String() string {
	return fmt.Sprintf("WSClient{%s (%s)}", c.Address, c.Endpoint)
}

func (c *WSClient) OnStart() error {
	err := c.dial()
	if err != nil {
		return err
	}

	c.RequestsCh = make(chan types.RPCRequest)

	c.send = make(chan types.RPCResponse)
	c.reconnectAfter = make(chan error, 1)
	c.backlog = make(chan types.RPCResponse, 1)

	c.startReadWriteRoutines()
	go c.reconnectRoutine()

	return nil
}

func (c *WSClient) Stop() error {
	if err := c.BaseService.Stop(); err != nil {
		return err
	}
	c.wg.Wait()
	close(c.RequestsCh)

	return nil
}

func (c *WSClient) IsReconnecting() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.reconnecting
}

func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsReconnecting()
}

func (c *WSClient) Send(ctx context.Context, request types.RPCResponse) error {
	select {
	case c.send <- request:
		c.Logger.Info("sent a request", "reqId", request.ID)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *WSClient) nextRequestID() types.JSONRPCIntID {
	c.mtx.Lock()
	id := c.nextReqID
	c.nextReqID++
	c.mtx.Unlock()
	return types.JSONRPCIntID(id)
}

func (c *WSClient) dial() error {
	dialer := &websocket.Dialer{
		NetDial: c.Dialer,
		Proxy:   http.ProxyFromEnvironment,
	}
	rHeader := http.Header{}
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)

	digestBz := crypto.Keccak256Hash([]byte(timeStr)).Bytes()
	sig, err := crypto.Sign(digestBz, c.PriKey)
	if err != nil {
		return err
	}

	rHeader.Set("pubKey", c.PubKey)
	rHeader.Set("sig", hex.EncodeToString(sig))
	rHeader.Set("time", timeStr)
	conn, _, err := dialer.Dial(c.protocol+"://"+c.Address+c.Endpoint, rHeader) // nolint:bodyclose
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *WSClient) reconnect() error {
	attempt := 0

	c.mtx.Lock()
	c.reconnecting = true
	c.mtx.Unlock()
	defer func() {
		c.mtx.Lock()
		c.reconnecting = false
		c.mtx.Unlock()
	}()

	for {
		backoffDuration := (1 << uint(attempt)) * time.Second
		if backoffDuration > defaultMaxSleepTime {
			backoffDuration = defaultMaxSleepTime
		}
		c.Logger.Info("reconnecting", "attempt", attempt+1, "backoff_duration", backoffDuration)
		time.Sleep(backoffDuration)

		err := c.dial()
		if err != nil {
			c.Logger.Error("failed to redial", "err", err)
		} else {
			c.Logger.Info("reconnected")
			if c.onReconnect != nil {
				go c.onReconnect()
			}
			return nil
		}

		attempt++

		if attempt > c.maxReconnectAttempts {
			return fmt.Errorf("reached maximum reconnect attempts: %w", err)
		}
	}
}

func (c *WSClient) startReadWriteRoutines() {
	c.wg.Add(2)
	c.readRoutineQuit = make(chan struct{})
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) processBacklog() error {
	select {
	case request := <-c.backlog:
		if c.writeWait > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
				c.Logger.Error("failed to set write deadline", "err", err)
			}
		}
		if err := c.conn.WriteJSON(request); err != nil {
			c.Logger.Error("failed to resend request", "err", err)
			c.reconnectAfter <- err
			c.backlog <- request
			return err
		}
		c.Logger.Info("resend a request", "req", request)
	default:
	}
	return nil
}

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:
			c.wg.Wait()
			if err := c.reconnect(); err != nil {
				c.Logger.Error("failed to reconnect", "err", err, "original_err", originalError)
				if err = c.Stop(); err != nil {
					c.Logger.Error("failed to stop conn", "error", err)
				}

				return
			}
		LOOP:
			for {
				select {
				case <-c.reconnectAfter:
				default:
					break LOOP
				}
			}
			err := c.processBacklog()
			if err == nil {
				c.startReadWriteRoutines()
			}

		case <-c.Quit():
			return
		}
	}
}

func (c *WSClient) writeRoutine() {
	var ticker *time.Ticker
	if c.pingPeriod > 0 {
		ticker = time.NewTicker(c.pingPeriod)
	} else {
		ticker = &time.Ticker{C: make(<-chan time.Time)}
	}

	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.wg.Done()
	}()

	for {
		select {
		case request := <-c.send:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteJSON(request); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.reconnectAfter <- err
				c.backlog <- request
				return
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.reconnectAfter <- err
				return
			}
			c.mtx.Lock()
			c.sentLastPingAt = time.Now()
			c.mtx.Unlock()
			c.Logger.Debug("sent ping")
		case <-c.readRoutineQuit:
			return
		case <-c.Quit():
			if err := c.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

func (c *WSClient) readRoutine() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
		c.wg.Done()
	}()

	c.conn.SetPongHandler(func(string) error {
		c.mtx.RLock()
		t := c.sentLastPingAt
		c.mtx.RUnlock()
		c.PingPongLatencyTimer.UpdateSince(t)

		c.Logger.Debug("got pong")
		return nil
	})

	for {
		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			c.Logger.Error("failed to read response", "err", err)
			close(c.readRoutineQuit)
			c.reconnectAfter <- err
			return
		}

		var request types.RPCRequest
		err = json.Unmarshal(data, &request)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}

		if err = validateResponseID(request.ID); err != nil {
			c.Logger.Error("error in response ID", "id", request.ID, "err", err)
			continue
		}

		c.Logger.Info("got request", "id", request.ID)

		select {
		case <-c.Quit():
		case c.RequestsCh <- request:
		}
	}
}
