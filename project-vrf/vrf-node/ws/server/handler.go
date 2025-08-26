package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/rpc/jsonrpc/types"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	legalTimeStampPeriod       = 5
	messageSignatureLength     = 64
	publicKeyLength            = 66
	defaultWSWriteChanCapacity = 100
	defaultWSWriteWait         = 10 * time.Second
	defaultWSReadWait          = 30 * time.Second
	defaultWSPingPeriod        = (defaultWSReadWait * 9) / 10
)

type WebsocketManager struct {
	websocket.Upgrader
	logger        log.Logger
	wsConnOptions []func(*wsConnection)
	recvChanMap   map[string]chan ResponseMsg
	rcRWLock      *sync.RWMutex
	sendChan      map[string]chan types.RPCRequest // node -> send channel
	aliveNodes    map[string]struct{}              // node -> struct{}{}
	scRWLock      *sync.RWMutex
}

func NewWebsocketManager(
	wsConnOptions ...func(*wsConnection),
) *WebsocketManager {
	return &WebsocketManager{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logger:        log.NewNopLogger(),
		wsConnOptions: wsConnOptions,
		recvChanMap:   make(map[string]chan ResponseMsg),
		rcRWLock:      &sync.RWMutex{},
		sendChan:      make(map[string]chan types.RPCRequest),
		aliveNodes:    make(map[string]struct{}),
		scRWLock:      &sync.RWMutex{},
	}
}

func (wm *WebsocketManager) SetWsConnOptions(wsConnOptions ...func(*wsConnection)) {
	wm.wsConnOptions = wsConnOptions
}

func (wm *WebsocketManager) SetLogger(l log.Logger) {
	wm.logger = l
}

func (wm *WebsocketManager) AliveNodes() []string {
	ret := make([]string, 0)
	for node := range wm.aliveNodes {
		ret = append(ret, node)
	}
	return ret
}

func (wm *WebsocketManager) RegisterResChannel(requestId string, recvChan chan ResponseMsg, stopChan chan struct{}) error {
	wm.rcRWLock.Lock()
	defer wm.rcRWLock.Unlock()
	wm.recvChanMap[requestId] = recvChan

	go func() {
		<-stopChan // block util stop
		wm.unregisterRecvChan(requestId)
	}()

	return nil
}

func (wm *WebsocketManager) SendMsg(msg RequestMsg) error {
	wm.scRWLock.RLock()
	defer wm.scRWLock.RUnlock()
	sendChan, ok := wm.sendChan[msg.TargetNode]
	if !ok {
		return fmt.Errorf("the node(%s) is lost", msg.TargetNode)
	}
	go func() {
		sendChan <- msg.RpcRequest
	}()
	return nil
}

func (wm *WebsocketManager) unregisterRecvChan(requestId string) {
	wm.rcRWLock.Lock()
	defer wm.rcRWLock.Unlock()
	delete(wm.recvChanMap, requestId)
}

func (wm *WebsocketManager) clientConnected(pubkey string, channel chan types.RPCRequest) {
	wm.scRWLock.Lock()
	defer wm.scRWLock.Unlock()

	if _, ok := wm.sendChan[pubkey]; !ok {
		wm.sendChan[pubkey] = channel
	}

	if wm.aliveNodes == nil {
		wm.aliveNodes = make(map[string]struct{})
	}
	wm.aliveNodes[pubkey] = struct{}{}
	wm.logger.Info("new node connected", "public key", pubkey)
}

func (wm *WebsocketManager) clientDisconnected(pubkey string) {
	wm.scRWLock.Lock()
	defer wm.scRWLock.Unlock()

	delete(wm.aliveNodes, pubkey)
	delete(wm.sendChan, pubkey)
	wm.logger.Info("node disconnected", "public key", pubkey)
}

func (wm *WebsocketManager) JudgeWssConnectPermission(activeTssMembers, inActiveTssMembers []string, nodePublicKey string) bool {
	for i := 0; i < len(activeTssMembers); i++ {
		if nodePublicKey == activeTssMembers[i] {
			return true
		}
	}
	for i := 0; i < len(inActiveTssMembers); i++ {
		if nodePublicKey == inActiveTssMembers[i] {
			return true
		}
	}
	return false
}

func (wm *WebsocketManager) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wm.Upgrade(w, r, nil)
	if err != nil {
		wm.logger.Error("Failed to upgrade connection", "err", err)
		return
	}
	defer func() {
		if err := wsConn.Close(); err != nil {
			wm.logger.Error("Failed to close connection", "err", err)
		}
	}()

	pubKey := r.Header.Get("pubKey")
	if len(pubKey) < publicKeyLength {
		wm.logger.Error("Failed to establish connection", "err", fmt.Errorf("invalid pubKey in header, expected length %d, actual length %d", publicKeyLength, len(pubKey)))
		return
	}
	sig := r.Header.Get("sig")
	if len(sig) < messageSignatureLength {
		wm.logger.Error("Failed to establish connection", "err", fmt.Errorf("failed to establish connection, expected length %d, actual length %d", messageSignatureLength, len(sig)))
		return
	}
	timeStr := r.Header.Get("time")
	if len(timeStr) == 0 {
		wm.logger.Error("Failed to establish connection", "err", fmt.Errorf("failed to establish connection, expected length %d, actual length %d", 0, len(timeStr)))
		return
	}
	timeInt64, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil || timeInt64 < 0 {
		wm.logger.Error("illegal timestamp", "err", err)
		return
	}
	if time.Now().Unix()-timeInt64 > legalTimeStampPeriod {
		wm.logger.Error("illegal timestamp", "err", errors.New("reject because illegal timestamp"))
		return
	}

	pubKeyBytes, pubErr := hex.DecodeString(pubKey)
	sigBytes, sigErr := hex.DecodeString(sig)
	if pubErr != nil || sigErr != nil {
		wm.logger.Error("hex decode error for pubkey or sig", "err", err)
		return
	}
	if len(sigBytes) < 64 {
		wm.logger.Error(fmt.Sprintf("invalid sigBytes, expected length is no less than 64, actual length is %d", len(sigBytes)))
		return
	}
	digestBz := crypto.Keccak256Hash([]byte(timeStr)).Bytes()
	if !crypto.VerifySignature(pubKeyBytes, digestBz, sigBytes[:64]) {
		wm.logger.Error("illegal signature", "publicKey", pubKey, "time", timeStr, "signature", sig)
		return
	}

	con := newWSConnection(wsConn, pubKey, wm.wsConnOptions...)
	con.SetLogger(wm.logger.With("remote", wsConn.RemoteAddr()))
	wm.logger.Info("New websocket connection", "remote", con.remoteAddr)

	err = con.Start() // BLOCKING
	if err != nil {
		wm.logger.Error("Failed to start connection", "err", err)
		return
	}
	if err := con.Stop(); err != nil {
		wm.logger.Error("error while stopping connection", "error", err)
	}
}

type wsConnection struct {
	service.BaseService
	remoteAddr        string
	baseConn          *websocket.Conn
	nodePublicKey     string
	responseChan      chan types.RPCResponse
	requestChan       chan types.RPCRequest
	readRoutineQuit   chan struct{}
	writeChanCapacity int
	writeWait         time.Duration
	readWait          time.Duration
	pingPeriod        time.Duration
	readLimit         int64
	onDisconnect      func(remoteAddr, pubKey string)
	ctx               context.Context
	cancel            context.CancelFunc
}

func newWSConnection(
	baseConn *websocket.Conn,
	publicKey string,
	options ...func(*wsConnection),
) *wsConnection {
	wsc := &wsConnection{
		remoteAddr:        baseConn.RemoteAddr().String(),
		baseConn:          baseConn,
		nodePublicKey:     publicKey,
		writeWait:         defaultWSWriteWait,
		writeChanCapacity: defaultWSWriteChanCapacity,
		readWait:          defaultWSReadWait,
		pingPeriod:        defaultWSPingPeriod,
		readRoutineQuit:   make(chan struct{}),
	}
	wsc.responseChan = make(chan types.RPCResponse, wsc.writeChanCapacity)
	wsc.requestChan = make(chan types.RPCRequest, wsc.writeChanCapacity)
	for _, option := range options {
		option(wsc)
	}
	wsc.baseConn.SetReadLimit(wsc.readLimit)
	wsc.BaseService = *service.NewBaseService(nil, "wsConnection", wsc)
	return wsc
}

func OnDisconnect(onDisconnect func(remoteAddr, pubKey string)) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.onDisconnect = onDisconnect
	}
}

func WriteWait(writeWait time.Duration) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.writeWait = writeWait
	}
}

func WriteChanCapacity(cap int) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.writeChanCapacity = cap
	}
}

func ReadWait(readWait time.Duration) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.readWait = readWait
	}
}

func PingPeriod(pingPeriod time.Duration) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.pingPeriod = pingPeriod
	}
}

func ReadLimit(readLimit int64) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.readLimit = readLimit
	}
}

func (wsc *wsConnection) OnStart() error {
	go wsc.readRoutine()
	wsc.writeRoutine()
	return nil
}

func (wsc *wsConnection) OnStop() {
	if wsc.onDisconnect != nil {
		wsc.onDisconnect(wsc.remoteAddr, wsc.nodePublicKey)
	}

	if wsc.ctx != nil {
		wsc.cancel()
	}
}

func (wsc *wsConnection) Output() chan types.RPCResponse {
	return wsc.responseChan
}

func (wsc *wsConnection) GetRemoteAddr() string {
	return wsc.remoteAddr
}

func (wsc *wsConnection) WriteRPCResponse(ctx context.Context, resp types.RPCResponse) error {
	select {
	case <-wsc.Quit():
		return errors.New("connection was stopped")
	case <-ctx.Done():
		return ctx.Err()
	case wsc.responseChan <- resp:
		return nil
	}
}

func (wsc *wsConnection) TryWriteRPCResponse(resp types.RPCResponse) bool {
	select {
	case <-wsc.Quit():
		return false
	case wsc.responseChan <- resp:
		return true
	default:
		return false
	}
}

func (wsc *wsConnection) Context() context.Context {
	if wsc.ctx != nil {
		return wsc.ctx
	}
	wsc.ctx, wsc.cancel = context.WithCancel(context.Background())
	return wsc.ctx
}

func (wsc *wsConnection) readRoutine() {
	writeCtx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("WSJSONRPC: %v", r)
			}
			wsc.Logger.Error("Panic in WSJSONRPC handler", "err", err, "stack", string(debug.Stack()))
			if err := wsc.WriteRPCResponse(writeCtx, types.RPCInternalError(types.JSONRPCIntID(-1), err)); err != nil {
				wsc.Logger.Error("Error writing RPC response", "err", err)
			}
			go wsc.readRoutine()
		}
	}()

	wsc.baseConn.SetPongHandler(func(m string) error {
		return wsc.baseConn.SetReadDeadline(time.Now().Add(wsc.readWait))
	})

	for {
		select {
		case <-wsc.Quit():
			return
		default:
			if err := wsc.baseConn.SetReadDeadline(time.Now().Add(wsc.readWait)); err != nil {
				wsc.Logger.Error("failed to set read deadline", "err", err)
			}

			_, r, err := wsc.baseConn.NextReader()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					wsc.Logger.Info("Client closed the connection")
				} else {
					wsc.Logger.Error("Failed to read request", "err", err)
				}
				if err := wsc.Stop(); err != nil {
					wsc.Logger.Error("Error closing websocket connection", "err", err)
				}
				close(wsc.readRoutineQuit)
				return
			}

			dec := json.NewDecoder(r)
			var response types.RPCResponse
			err = dec.Decode(&response)
			if err != nil {
				wsc.Logger.Error("error unmarshaling response", "err", err)
				continue
			}

			if response.ID == nil {
				wsc.Logger.Info("[WS]received response with no ID, drop it")
				continue
			}
			if err := wsc.WriteRPCResponse(writeCtx, response); err != nil {
				wsc.Logger.Error("Error writing RPC response", "err", err)
			}

		}
	}
}

func (wsc *wsConnection) writeRoutine() {
	pingTicker := time.NewTicker(wsc.pingPeriod)
	defer pingTicker.Stop()

	pongs := make(chan string, 1)
	wsc.baseConn.SetPingHandler(func(m string) error {
		select {
		case pongs <- m:
		default:
		}
		return nil
	})

	for {
		select {
		case <-wsc.Quit():
			return
		case <-wsc.readRoutineQuit: // error in readRoutine
			return
		case m := <-pongs:
			err := wsc.writeMessageWithDeadline(websocket.PongMessage, []byte(m))
			if err != nil {
				wsc.Logger.Info("Failed to write pong (client may disconnect)", "err", err)
			}
		case <-pingTicker.C:
			err := wsc.writeMessageWithDeadline(websocket.PingMessage, []byte{})
			if err != nil {
				wsc.Logger.Error("Failed to write ping", "err", err)
				return
			}
		case msg := <-wsc.requestChan:
			wsc.Logger.Info("send msg from requestChan to target client", "method", msg.Method)
			jsonBytes, err := json.Marshal(msg)
			if err != nil {
				wsc.Logger.Error("Failed to marshal RPCRequest to JSON", "err", err)
				continue
			}
			if err = wsc.writeMessageWithDeadline(websocket.TextMessage, jsonBytes); err != nil {
				wsc.Logger.Error("Failed to write request", "err", err, "msg", msg)
				return
			}
		}
	}
}

func (wsc *wsConnection) writeMessageWithDeadline(msgType int, msg []byte) error {
	if err := wsc.baseConn.SetWriteDeadline(time.Now().Add(wsc.writeWait)); err != nil {
		return err
	}
	return wsc.baseConn.WriteMessage(msgType, msg)
}
