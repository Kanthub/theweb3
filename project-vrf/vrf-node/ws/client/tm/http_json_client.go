package tm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	tmsync "github.com/tendermint/tendermint/libs/sync"
	"github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

const (
	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
	protoUNIX  = "unix"
)

type parsedURL struct {
	url.URL

	isUnixSocket bool
}

func newParsedURL(remoteAddr string) (*parsedURL, error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = protoTCP
	}

	pu := &parsedURL{
		URL:          *u,
		isUnixSocket: false,
	}

	if u.Scheme == protoUNIX {
		pu.isUnixSocket = true
	}

	return pu, nil
}

func (u *parsedURL) SetDefaultSchemeHTTP() {
	switch u.Scheme {
	case protoHTTP, protoHTTPS, protoWS, protoWSS:
	default:
		// default to http for unknown protocols (ex. tcp)
		u.Scheme = protoHTTP
	}
}

func (u parsedURL) GetHostWithPath() string {
	// Remove protocol, userinfo and # fragment, assume opaque is empty
	return u.Host + u.EscapedPath()
}

func (u parsedURL) GetTrimmedHostWithPath() string {
	if !u.isUnixSocket {
		return u.GetHostWithPath()
	}
	return strings.ReplaceAll(u.GetHostWithPath(), "/", ".")
}

func (u parsedURL) GetDialAddress() string {
	if !u.isUnixSocket {
		return u.Host
	}
	return u.GetHostWithPath()
}

func (u parsedURL) GetTrimmedURL() string {
	return u.Scheme + "://" + u.GetTrimmedHostWithPath()
}

type HTTPClient interface {
	Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (interface{}, error)
}

type Caller interface {
	Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (interface{}, error)
}

type Client struct {
	address  string
	username string
	password string

	client *http.Client

	mtx       tmsync.Mutex
	nextReqID int
}

var _ HTTPClient = (*Client)(nil)
var _ Caller = (*Client)(nil)
var _ Caller = (*RequestBatch)(nil)

func New(remote string) (*Client, error) {
	httpClient, err := DefaultHTTPClient(remote)
	if err != nil {
		return nil, err
	}
	return NewWithHTTPClient(remote, httpClient)
}

func NewWithHTTPClient(remote string, client *http.Client) (*Client, error) {
	if client == nil {
		panic("nil http.Client provided")
	}

	parsedURL, err := newParsedURL(remote)
	if err != nil {
		return nil, fmt.Errorf("invalid remote %s: %s", remote, err)
	}

	parsedURL.SetDefaultSchemeHTTP()

	address := parsedURL.GetTrimmedURL()
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	rpcClient := &Client{
		address:  address,
		username: username,
		password: password,
		client:   client,
	}

	return rpcClient, nil
}

func (c *Client) Call(
	ctx context.Context,
	method string,
	params map[string]interface{},
	result interface{},
) (interface{}, error) {
	id := c.nextRequestID()

	request, err := types.MapToRequest(id, method, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode params: %w", err)
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	requestBuf := bytes.NewBuffer(requestBytes)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.address, requestBuf)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	if c.username != "" || c.password != "" {
		httpRequest.SetBasicAuth(c.username, c.password)
	}

	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("post failed: %w", err)
	}

	defer httpResponse.Body.Close()

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return unmarshalResponseBytes(responseBytes, id, result)
}

// NewRequestBatch starts a batch of requests for this client.
func (c *Client) NewRequestBatch() *RequestBatch {
	return &RequestBatch{
		requests: make([]*jsonRPCBufferedRequest, 0),
		client:   c,
	}
}

func (c *Client) sendBatch(ctx context.Context, requests []*jsonRPCBufferedRequest) ([]interface{}, error) {
	reqs := make([]types.RPCRequest, 0, len(requests))
	results := make([]interface{}, 0, len(requests))
	for _, req := range requests {
		reqs = append(reqs, req.request)
		results = append(results, req.result)
	}

	requestBytes, err := json.Marshal(reqs)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.address, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	if c.username != "" || c.password != "" {
		httpRequest.SetBasicAuth(c.username, c.password)
	}

	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}

	defer httpResponse.Body.Close()

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	ids := make([]types.JSONRPCIntID, len(requests))
	for i, req := range requests {
		ids[i] = req.request.ID.(types.JSONRPCIntID)
	}

	return unmarshalResponseBytesArray(responseBytes, ids, results)
}

func (c *Client) nextRequestID() types.JSONRPCIntID {
	c.mtx.Lock()
	id := c.nextReqID
	c.nextReqID++
	c.mtx.Unlock()
	return types.JSONRPCIntID(id)
}

type jsonRPCBufferedRequest struct {
	request types.RPCRequest
	result  interface{}
}

type RequestBatch struct {
	client *Client

	mtx      tmsync.Mutex
	requests []*jsonRPCBufferedRequest
}

func (b *RequestBatch) Count() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	return len(b.requests)
}

func (b *RequestBatch) enqueue(req *jsonRPCBufferedRequest) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.requests = append(b.requests, req)
}

func (b *RequestBatch) Clear() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	return b.clear()
}

func (b *RequestBatch) clear() int {
	count := len(b.requests)
	b.requests = make([]*jsonRPCBufferedRequest, 0)
	return count
}

func (b *RequestBatch) Send(ctx context.Context) ([]interface{}, error) {
	b.mtx.Lock()
	defer func() {
		b.clear()
		b.mtx.Unlock()
	}()
	return b.client.sendBatch(ctx, b.requests)
}

func (b *RequestBatch) Call(
	_ context.Context,
	method string,
	params map[string]interface{},
	result interface{},
) (interface{}, error) {
	id := b.client.nextRequestID()
	request, err := types.MapToRequest(id, method, params)
	if err != nil {
		return nil, err
	}
	b.enqueue(&jsonRPCBufferedRequest{request: request, result: result})
	return result, nil
}

func makeHTTPDialer(remoteAddr string) (func(string, string) (net.Conn, error), error) {
	u, err := newParsedURL(remoteAddr)
	if err != nil {
		return nil, err
	}

	protocol := u.Scheme

	switch protocol {
	case protoHTTP, protoHTTPS:
		protocol = protoTCP
	}

	dialFn := func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, u.GetDialAddress())
	}

	return dialFn, nil
}

func DefaultHTTPClient(remoteAddr string) (*http.Client, error) {
	dialFn, err := makeHTTPDialer(remoteAddr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
			Dial:               dialFn,
		},
	}

	return client, nil
}
