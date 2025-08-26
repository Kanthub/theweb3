package dapplink_vrf

import (
	"context"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	common2 "github.com/the-web3/dapplink-vrf/common"
	"github.com/the-web3/dapplink-vrf/config"
	"github.com/the-web3/dapplink-vrf/database"
	"github.com/the-web3/dapplink-vrf/driver"
	"github.com/the-web3/dapplink-vrf/event"
	"github.com/the-web3/dapplink-vrf/synchronizer"
	"github.com/the-web3/dapplink-vrf/synchronizer/node"
	"github.com/the-web3/dapplink-vrf/worker"
)

type DappLinkVrf struct {
	db *database.DB

	synchronizer  *synchronizer.Synchronizer
	eventsHandler *event.EventsHandler
	worker        *worker.Worker

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewDappLinkVrf(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*DappLinkVrf, error) {

	ethClient, err := node.DialEthClient(ctx, cfg.Chain.ChainRpcUrl)
	if err != nil {
		log.Error("new eth client fail", "err", err)
		return nil, err
	}

	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("new database fail", "err", err)
		return nil, err
	}

	synchronizerS, err := synchronizer.NewSynchronizer(cfg, db, ethClient, shutdown)
	if err != nil {
		log.Error("new synchronizer fail", "err", err)
		return nil, err
	}

	eventConfigm := &event.EventsHandlerConfig{
		DappLinkVrfAddress:        cfg.Chain.DappLinkVrfContractAddress,
		DappLinkVrfFactoryAddress: cfg.Chain.DappLinkVrfFactoryContractAddress,
		LoopInterval:              cfg.Chain.EventInterval,
		StartHeight:               big.NewInt(int64(cfg.Chain.StartingHeight)),
		Epoch:                     500,
	}

	eventHandler, err := event.NewEventsHandler(db, eventConfigm, shutdown)
	if err != nil {
		return nil, err
	}

	ethcli, err := driver.EthClientWithTimeout(ctx, cfg.Chain.ChainRpcUrl)
	if err != nil {
		log.Error("new eth cli fail", "err", err)
		return nil, err
	}

	callerPrivateKey, _, err := common2.ParseWalletPrivKeyAndContractAddr(
		"ContractCaller", cfg.Chain.Mnemonic, cfg.Chain.CallerHDPath,
		cfg.Chain.PrivateKey, cfg.Chain.DappLinkVrfContractAddress, cfg.Chain.Passphrase,
	)

	decg := &driver.DriverEingineConfig{
		ChainClient:               ethcli,
		ChainId:                   big.NewInt(int64(cfg.Chain.ChainId)),
		DappLinkVrfAddress:        common.HexToAddress(cfg.Chain.DappLinkVrfContractAddress),
		CallerAddress:             common.HexToAddress(cfg.Chain.CallerAddress),
		PrivateKey:                callerPrivateKey,
		NumConfirmations:          cfg.Chain.Confirmations,
		SafeAbortNonceTooLowCount: cfg.Chain.SafeAbortNonceTooLowCount,
	}
	eingine, err := driver.NewDriverEingine(ctx, decg)
	if err != nil {
		log.Error("new driver eingine fail", "err", err)
		return nil, err
	}

	workerConfig := &worker.WorkerConfig{
		LoopInterval: cfg.Chain.CallInterval,
	}

	workerProcessor, err := worker.NewWorker(db, eingine, workerConfig, shutdown)
	if err != nil {
		log.Error("new event processor fail", "err", err)
		return nil, err
	}

	return &DappLinkVrf{
		db:            db,
		synchronizer:  synchronizerS,
		eventsHandler: eventHandler,
		worker:        workerProcessor,
		shutdown:      shutdown,
	}, nil
}

func (dvrf *DappLinkVrf) Start(ctx context.Context) error {
	err := dvrf.synchronizer.Start()
	if err != nil {
		return err
	}
	err = dvrf.eventsHandler.Start()
	if err != nil {
		return err
	}
	err = dvrf.worker.Start()
	if err != nil {
		return err
	}
	return nil
}

func (dvrf *DappLinkVrf) Stop(ctx context.Context) error {
	err := dvrf.synchronizer.Close()
	if err != nil {
		return err
	}
	err = dvrf.eventsHandler.Close()
	if err != nil {
		return err
	}
	err = dvrf.worker.Close()
	if err != nil {
		return err
	}
	return nil
}

func (dvrf *DappLinkVrf) Stopped() bool {
	return dvrf.stopped.Load()
}
