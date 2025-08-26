package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/dapplink-vrf/common/tasks"
	"github.com/the-web3/dapplink-vrf/database"
	"github.com/the-web3/dapplink-vrf/driver"
)

type WorkerConfig struct {
	LoopInterval time.Duration
}

type Worker struct {
	workerConfig *WorkerConfig
	db           *database.DB
	deg          *driver.DriverEingine

	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewWorker(db *database.DB, deg *driver.DriverEingine, workerConfig *WorkerConfig, shutdown context.CancelCauseFunc) (*Worker, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Worker{
		db:             db,
		deg:            deg,
		workerConfig:   workerConfig,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in bridge processor: %w", err))
		}},
	}, nil
}

func (wk *Worker) Start() error {
	log.Info("starting worker processor...")
	tickerEventWorker := time.NewTicker(wk.workerConfig.LoopInterval)
	wk.tasks.Go(func() error {
		for range tickerEventWorker.C {
			log.Info("start handler random for vrf")
			err := wk.ProcessCallerVrf()
			if err != nil {
				log.Error("process caller vrf fail", "err", err)
				return err
			}
		}
		return nil
	})
	return nil
}

func (wk *Worker) Close() error {
	wk.resourceCancel()
	return wk.tasks.Wait()
}

func (wk *Worker) ProcessCallerVrf() error {
	// 获取 RequestSent 合约事件
	requestUnsentList, err := wk.db.RequestSend.QueryUnHandleRequestSendList()
	if err != nil {
		log.Error("query unhandle request sent list fail", "err", err)
		return err
	}
	for _, requestUnsent := range requestUnsentList {
		// 组装随机数据列表交易发到 Vrf 合约
		txRecepit, err := wk.deg.FulfillRandomWords(requestUnsent.RequestId, nil)
		if err != nil {
			log.Error("Fulfill random words fail", "err", err)
			return err
		}
		if txRecepit.Status == 1 {
			// 更新 RequestSent 合约事件表的状态
			err := wk.db.RequestSend.MarkRequestSendFinish(requestUnsent)
			if err != nil {
				log.Error("mark request sent event list fail", "err", err)
				return err
			}
		}
	}
	return nil
}

/*

 */
