package driver

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/dapplink-vrf/bindings"
	"github.com/the-web3/dapplink-vrf/txmgr"
)

var (
	errMaxPriorityFeePerGasNotFound = errors.New(
		"Method eth_maxPriorityFeePerGas not found",
	)
	FallbackGasTipCap = big.NewInt(1500000000)
)

type DriverEingineConfig struct {
	ChainClient               *ethclient.Client
	ChainId                   *big.Int
	DappLinkVrfAddress        common.Address
	CallerAddress             common.Address
	PrivateKey                *ecdsa.PrivateKey // CallerAddress 和 PrivateKey 是一一对应的
	NumConfirmations          uint64
	SafeAbortNonceTooLowCount uint64
}

type DriverEingine struct {
	Ctx                    context.Context
	Cfg                    *DriverEingineConfig
	DappLinkVrfContract    *bindings.DappLinkVRF
	RawDappLinkVrfContract *bind.BoundContract
	DappLinkVrfContractAbi *abi.ABI
	TxMgr                  txmgr.TxManager
	cancel                 func()
	wg                     sync.WaitGroup
}

func NewDriverEingine(ctx context.Context, cfg *DriverEingineConfig) (*DriverEingine, error) {
	_, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	dappLinkVrfContract, err := bindings.NewDappLinkVRF(cfg.DappLinkVrfAddress, cfg.ChainClient)
	if err != nil {
		log.Error("new dapplink vrf fail", "err", err)
		return nil, err
	}

	parsed, err := abi.JSON(strings.NewReader(bindings.DappLinkVRFMetaData.ABI))
	if err != nil {
		log.Error("parsed abi fail", "err", err)
		return nil, err
	}

	dappLinkVrfContractAbi, err := bindings.DappLinkVRFMetaData.GetAbi()
	if err != nil {
		log.Error("get dapplink vrf meta data fail", "err", err)
		return nil, err
	}

	rawDappLinkVrfContract := bind.NewBoundContract(cfg.DappLinkVrfAddress, parsed, cfg.ChainClient, cfg.ChainClient, cfg.ChainClient)

	txManagerConfig := txmgr.Config{
		ResubmissionTimeout:       time.Second * 5,
		ReceiptQueryInterval:      time.Second,
		NumConfirmations:          cfg.NumConfirmations,
		SafeAbortNonceTooLowCount: cfg.SafeAbortNonceTooLowCount,
	}

	txManager := txmgr.NewSimpleTxManager(txManagerConfig, cfg.ChainClient)

	return &DriverEingine{
		Ctx:                    ctx,
		Cfg:                    cfg,
		DappLinkVrfContract:    dappLinkVrfContract,
		RawDappLinkVrfContract: rawDappLinkVrfContract,
		DappLinkVrfContractAbi: dappLinkVrfContractAbi,
		TxMgr:                  txManager,
		cancel:                 cancel,
	}, nil
}

func (de *DriverEingine) UpdateGasPrice(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	var opts *bind.TransactOpts
	var err error
	opts, err = bind.NewKeyedTransactorWithChainID(de.Cfg.PrivateKey, de.Cfg.ChainId)
	if err != nil {
		log.Error("new keyed transactor with chain id fail", "err", err)
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(tx.Nonce())
	opts.NoSend = true

	finalTx, err := de.RawDappLinkVrfContract.RawTransact(opts, tx.Data())
	switch {
	case err == nil:
		return finalTx, nil

	case de.isMaxPriorityFeePerGasNotFoundError(err):
		log.Info("Don't support priority fee")
		opts.GasTipCap = FallbackGasTipCap
		return de.RawDappLinkVrfContract.RawTransact(opts, tx.Data())

	default:
		return nil, err
	}
}

func (de *DriverEingine) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return de.Cfg.ChainClient.SendTransaction(ctx, tx)
}

func (de *DriverEingine) isMaxPriorityFeePerGasNotFoundError(err error) bool {
	return strings.Contains(err.Error(), errMaxPriorityFeePerGasNotFound.Error())
}

func (de *DriverEingine) fulfillRandomWords(ctx context.Context, requestId *big.Int, randomList []*big.Int) (*types.Transaction, error) {
	nonce, err := de.Cfg.ChainClient.NonceAt(ctx, de.Cfg.CallerAddress, nil)
	if err != nil {
		log.Error("get nonce error", "err", err)
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(de.Cfg.PrivateKey, de.Cfg.ChainId)
	if err != nil {
		log.Error("new keyed transactor with chain id fail", "err", err)
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.NoSend = true

	tx, err := de.DappLinkVrfContract.FulfillRandomWords(opts, requestId, randomList)
	switch {
	case err == nil:
		return tx, nil

	case de.isMaxPriorityFeePerGasNotFoundError(err):
		log.Info("Don't support priority fee")
		opts.GasTipCap = FallbackGasTipCap
		return de.DappLinkVrfContract.FulfillRandomWords(opts, requestId, randomList)

	default:
		return nil, err
	}
}

func (de *DriverEingine) FulfillRandomWords(requestId *big.Int, randomList []*big.Int) (*types.Receipt, error) {
	tx, err := de.fulfillRandomWords(de.Ctx, requestId, randomList)
	if err != nil {
		log.Error("build request random words tx fail", "err", err)
		return nil, err
	}
	updateGasPrice := func(ctx context.Context) (*types.Transaction, error) {
		return de.UpdateGasPrice(ctx, tx)
	}

	receipt, err := de.TxMgr.Send(de.Ctx, updateGasPrice, de.SendTransaction)
	if err != nil {
		log.Error("send tx fail", "err", err)
		return nil, err
	}
	return receipt, nil
}
