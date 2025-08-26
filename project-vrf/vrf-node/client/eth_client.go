package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultDialTimeout = 5 * time.Second
)

// 安全地连接 RPC 节点，支持禁用 HTTP/2 和连接超时
func DialEthClientWithTimeout(ctx context.Context, url string, disableHTTP2 bool) (*ethclient.Client, error) {
	// 创建一个 5 秒超时的子 context。防止连接操作阻塞太久。
	ctxt, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()
	if strings.HasPrefix(url, "http") {
		httpClient := new(http.Client)
		httpClient.Timeout = defaultDialTimeout

		if disableHTTP2 {
			log.Debug("Disabled HTTP/2 support in  eth client")
			httpClient.Transport = &http.Transport{
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			}
		}
		rpcClient, err := rpc.DialHTTPWithClient(url, httpClient)
		if err != nil {
			return nil, err
		}
		return ethclient.NewClient(rpcClient), nil
	}
	// 如果不是 HTTP 协议，比如 ws:// 或 ipc://，就使用默认的 Geth DialContext 方法连接
	return ethclient.DialContext(ctxt, url)
}

// 用于生成带链 ID 的交易签名配置（TransactOpts）
func NewTransactOpts(ctx context.Context, chainId uint64, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	var opts *bind.TransactOpts
	var err error

	if privateKey == nil {
		return nil, errors.New("no private key provided")
	}

	opts, err = bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainId))
	if err != nil {
		return nil, fmt.Errorf("new keyed transactor fail, err: %v", err)
	}

	opts.Context = ctx
	opts.NoSend = true // 表示不会自动发送交易

	return opts, err
}

// 典型的“等待交易上链”的封装函数，适合异步交易确认逻辑
func GetTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	ticker := time.NewTicker(10 * time.Second) // 设置每 10 秒轮询一次
	for {
		<-ticker.C
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		if errors.Is(err, ethereum.NotFound) {
			continue
		}
		return receipt, nil
	}
}
