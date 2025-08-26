package vrf_node

import (
	"context"

	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	clien2 "github.com/the-web3-contracts/vrf-node/client"
	common2 "github.com/the-web3-contracts/vrf-node/common"
	"github.com/the-web3-contracts/vrf-node/config"
	"github.com/the-web3-contracts/vrf-node/database"
	"github.com/the-web3-contracts/vrf-node/event"
	"github.com/the-web3-contracts/vrf-node/synchronizer"
	"github.com/the-web3-contracts/vrf-node/synchronizer/node"
	"github.com/the-web3-contracts/vrf-node/worker"
)

const BlockSize = 3000

type VrfNode struct {
	db *database.DB

	synchronizer *synchronizer.Synchronizer
	eventsParser *event.EventsParser
	worker       *worker.Worker

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewVrfNode(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*VrfNode, error) {
	ethClient, err := node.DialEthClient(ctx, cfg.Chain.ChainRpcUrl)
	if err != nil {
		log.Error("new eth syncer client fail", "err", err)
		return nil, err
	}

	callEthClient, err := clien2.DialEthClientWithTimeout(context.Background(), cfg.Chain.ChainRpcUrl, false)
	// 支持设置超时、避免调用阻塞
	if err != nil {
		log.Error("new eth caller client fail", "err", err)
		return nil, err
	}

	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("new database fail", "err", err)
		return nil, err
	}

	syncer, err := synchronizer.NewSynchronizer(cfg, db, ethClient, shutdown)
	if err != nil {
		log.Error("new synchronizer fail", "err", err)
		return nil, err
	}

	epConfig := &event.EventsParserConfig{
		DappLinkVrfAddress:        cfg.Chain.DappLinkVrfContractAddress,
		DappLinkVrfFactoryAddress: cfg.Chain.DappLinkVrfFactoryContractAddress,
		EventLoopInterval:         cfg.Chain.EventInterval,
		StartHeight:               big.NewInt(int64(cfg.Chain.StartingHeight)),
		BlockSize:                 BlockSize,
	}

	eventsParser, err := event.NewEventsParser(db, epConfig, shutdown)

	ecdsaPrivteKey, _ := common2.ParsePrivateKeyStr(cfg.Chain.PrivateKey)

	workConf := &worker.WorkerConfig{
		ChainClient:               callEthClient,
		ChainId:                   big.NewInt(int64(cfg.Chain.ChainId)),
		DappLinkVrfManagerAddress: common.HexToAddress(cfg.Chain.DappLinkVrfContractAddress),
		CallerAddress:             common.HexToAddress(cfg.Chain.CallerAddress),
		PrivateKey:                ecdsaPrivteKey,
		NumConfirmations:          cfg.Chain.NumConfirmations,
		SafeAbortNonceToLowCount:  cfg.Chain.SafeAbortNonceTooLowCount,
		LoopInternal:              cfg.Chain.CallInterval,
	}

	workerF, err := worker.NewWorker(db, workConf, shutdown)
	if err != nil {
		log.Error("new worker fail", "err", err)
		return nil, err
	}

	return &VrfNode{
		db:           db,
		synchronizer: syncer,
		eventsParser: eventsParser,
		worker:       workerF,
		shutdown:     shutdown,
	}, nil
}

func (vn *VrfNode) Start(ctx context.Context) error { // 并没有用到传入的ctx
	err := vn.synchronizer.Start()
	if err != nil {

		return err
	}
	err = vn.eventsParser.Start()
	if err != nil {
		return err
	}
	err = vn.worker.Start()
	if err != nil {
		return err
	}
	return nil
}

func (vn *VrfNode) Stop(ctx context.Context) error { // 并没有用到传入的ctx
	if vn.stopped.Load() {
		log.Error("Vrf Node already stopped")
		return nil
	}
	vn.stopped.Store(true)

	err := vn.synchronizer.Close()
	if err != nil {
		return err
	}
	err = vn.eventsParser.Close()
	if err != nil {
		return err
	}
	err = vn.worker.Close()
	if err != nil {
		return err
	}
	return nil
}

func (vn *VrfNode) Stopped() bool {
	return vn.stopped.Load()
}
