package driver

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3-contracts/vrf-node/bindings/vrf"
	"github.com/the-web3-contracts/vrf-node/txmgr"
)

var (
	errMaxPriorityFeePerGasNotFound = errors.New(
		"Method eth_maxPriorityFeePerGas not found",
	)
	FallbackGasTipCap = big.NewInt(1500000000)
)

type CallerConfig struct {
	ChainClient               *ethclient.Client
	ChainId                   *big.Int
	DappLinkVrfManagerAddress common.Address
	CallerAddress             common.Address
	PrivateKey                *ecdsa.PrivateKey
	NumConfirmations          uint64
	SafeAbortNonceToLowCount  uint64
}

type Caller struct {
	Ctx                     context.Context
	Cfg                     *CallerConfig
	DappLinkVrfContracts    *vrf.DappLinkVRFManager
	RawDappLinkVrfContracts *bind.BoundContract
	DappLinkVrfContractsAbi *abi.ABI
	TxMrg                   txmgr.TxManager
}

func NewCaller(Ctx context.Context, Cfg *CallerConfig) (*Caller, error) {
	dappLinkVrfContracts, err := vrf.NewDappLinkVRFManager(Cfg.DappLinkVrfManagerAddress, Cfg.ChainClient)
	if err != nil {
		log.Error("New DappLink Vrf Manager Fail", "err", err)
		return nil, err
	}

	// 获取合约ABI
	parsed, err := abi.JSON(strings.NewReader(vrf.DappLinkVRFManagerMetaData.ABI))
	if err != nil {
		log.Error("abi parsed fail", "err", err)
		return nil, err
	}

	dappLinkVrfContractsAbi, err := vrf.DappLinkVRFManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get abi fail", "err", err)
		return nil, err
	}

	// 构造绑定底层原始合约
	rawDappLinkVrfContracts := bind.NewBoundContract(Cfg.DappLinkVrfManagerAddress, parsed, Cfg.ChainClient, Cfg.ChainClient, Cfg.ChainClient)

	txManagerConfig := txmgr.Config{
		ResubmissionTimeout:       time.Second * 5,
		ReceiptQueryInterval:      time.Second,
		NumConfirmations:          Cfg.NumConfirmations,
		SafeAbortNonceTooLowCount: Cfg.SafeAbortNonceToLowCount,
	}

	txManager := txmgr.NewSimpleTxManager(txManagerConfig, Cfg.ChainClient)

	return &Caller{
		Ctx:                     Ctx,
		Cfg:                     Cfg,
		DappLinkVrfContracts:    dappLinkVrfContracts,
		RawDappLinkVrfContracts: rawDappLinkVrfContracts,
		DappLinkVrfContractsAbi: dappLinkVrfContractsAbi,
		TxMrg:                   txManager,
	}, nil
}

func (caller *Caller) isMaxPriorityFeePerGasNotFoundError(err error) bool {
	return strings.Contains(err.Error(), errMaxPriorityFeePerGasNotFound.Error())
}

func (caller *Caller) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return caller.Cfg.ChainClient.SendTransaction(ctx, tx)
}

func (caller *Caller) UpdateGasPrice(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	var opts *bind.TransactOpts
	var err error

	// 交易实例
	opts, err = bind.NewKeyedTransactorWithChainID(caller.Cfg.PrivateKey, caller.Cfg.ChainId)
	if err != nil {
		log.Error("new keyed transactor with chain id fail", "err", err)
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(tx.Nonce())
	opts.NoSend = true

	// 构造raw transaction
	finalTx, err := caller.RawDappLinkVrfContracts.RawTransact(opts, tx.Data())
	switch {
	case err == nil:
		return finalTx, nil

	case caller.isMaxPriorityFeePerGasNotFoundError(err):
		log.Info("Don't support priority fee")
		opts.GasTipCap = FallbackGasTipCap
		return caller.RawDappLinkVrfContracts.RawTransact(opts, tx.Data())

	default:
		return nil, err
	}
}

func (caller *Caller) fulfillRandomWords(ctx context.Context, requestId *big.Int, randomList []*big.Int) (*types.Transaction, error) {
	// 获取nonce
	nonce, err := caller.Cfg.ChainClient.NonceAt(ctx, caller.Cfg.CallerAddress, nil)
	if err != nil {
		log.Error("get eth nonce fail", "err", err)
		return nil, err
	}

	// 获取交易发送器 opts
	opts, err := bind.NewKeyedTransactorWithChainID(caller.Cfg.PrivateKey, caller.Cfg.ChainId)
	if err != nil {
		log.Error("new keyed transactor with chain id fail", "err", err)
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.NoSend = true

	var msgHash [32]byte
	var blsParam vrf.IBLSApkRegistryVrfNoSignerAndSignature

	// 在go binding 合约 构造交易
	// 使用强类型合约（DappLinkVrfContracts）调用合约函数 FulfillRandomWords
	tx, err := caller.DappLinkVrfContracts.FulfillRandomWords(opts, requestId, randomList, msgHash, big.NewInt(100), blsParam)

	switch {
	case err == nil:
		return tx, nil

	case caller.isMaxPriorityFeePerGasNotFoundError(err):
		log.Info("Don't support priority fee")
		opts.GasTipCap = FallbackGasTipCap
		return caller.DappLinkVrfContracts.FulfillRandomWords(opts, requestId, randomList, msgHash, big.NewInt(100), blsParam)
	// 合约构造交易时报错，有可能是当前链不支持 EIP-1559 的 maxPriorityFeePerGas 模式。
	// 所以它 fallback 到 legacy 模式下，设置 GasTipCap 为固定值，再重试一次。
	default:
		log.Error("fulfill random words fail", "err", err)
		return nil, err
	}
}

func (caller *Caller) FulfillRandomWords(requestId *big.Int, randomList []*big.Int) (*types.Receipt, error) {

	//  强类型合约调用 + raw 合约调用 的组合使用技巧
	tx, err := caller.fulfillRandomWords(caller.Ctx, requestId, randomList)
	if err != nil {
		log.Error("build request random words tx fail", "err", err)
		return nil, err
	}
	// 强类型合约调用提供 标准交易体tx（你不用自己拼 calldata）

	// raw 调用可以 手动调整一些参数（比如 GasTipCap、GasFeeCap）
	// raw 可以直接重新构造一个 types.Transaction 实例（把 gas 设置得更合适）
	updateGasPrice := func(ctx context.Context) (*types.Transaction, error) {
		return caller.UpdateGasPrice(ctx, tx)
	}

	// 用txManager 发送 rawTransaction
	receipt, err := caller.TxMrg.Send(caller.Ctx, updateGasPrice, caller.SendTransaction)
	if err != nil {
		log.Error("send tx fail", "err", err)
		return nil, err
	}
	return receipt, nil
}
