package tasker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/the-web3/rpc-server/common/tasks"
	"github.com/the-web3/rpc-server/config"
	"github.com/the-web3/rpc-server/database"
	"github.com/the-web3/rpc-server/market_price"
)

type MarketPriceTasker struct {
	loopInternal   time.Duration
	db             *database.DB
	exClient       *market_price.Client
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
	stopped        atomic.Bool
}

func NewMarketPriceTasker(ctx context.Context, db *database.DB, conf *config.Config, shutdown context.CancelCauseFunc) (*MarketPriceTasker, error) {

	exClient, err := market_price.NewExchangeClient(conf.BaseUrl)
	if err != nil {
		log.Error("new exchange client fail", "err", err)
		return nil, err
	}

	resCtx, resCancel := context.WithCancel(ctx)
	return &MarketPriceTasker{
		db:           db,
		loopInternal: conf.LoopInternal,
		exClient:     exClient,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in market price: %w", err))
		}},
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
	}, nil
}

func (mpt *MarketPriceTasker) Start(ctx context.Context) error {
	tickerTask := time.NewTicker(mpt.loopInternal)
	mpt.tasks.Go(func() error {
		for range tickerTask.C {
			tokenList := mpt.db.SupportToken.ListSupportToken()
			var tokenPriceList []database.TokenPrice
			for _, tokenItem := range tokenList {
				log.Info("handle token price", "tokeName", tokenItem.BaseAsset)
				marketPrice, err := mpt.exClient.GetMarketPriceByBaseAsset(tokenItem.BaseAsset)
				if err != nil || marketPrice == nil {
					log.Error("get market price fail", "err", err)
					return err
				}
				log.Info("get mnarket price", "marketPrice", marketPrice.Price)
				tokePriceItem := database.TokenPrice{
					GUID:      uuid.New(),
					TokenId:   tokenItem.GUID.String(),
					SellPrice: fmt.Sprintf("%2f", marketPrice.Price),
					AvgPrice:  fmt.Sprintf("%2f", marketPrice.Price),
					Rate:      "0.81",
					Volume24h: fmt.Sprintf("%2f", marketPrice.Volume24h),
					Timestamp: uint64(time.Now().Unix()),
				}
				tokenPriceList = append(tokenPriceList, tokePriceItem)
			}
			if len(tokenPriceList) > 0 {
				log.Info("tokenPriceList length", "length", len(tokenPriceList))
				err := mpt.db.TokenPrice.StoreAndUpdateTokenPrices(tokenPriceList)
				if err != nil {
					log.Error("update token price fail", "err", err)
					return err
				}
			}
		}
		return nil
	})
	return nil
}

func (mpt *MarketPriceTasker) Stop(ctx context.Context) error {
	mpt.resourceCancel()
	return mpt.tasks.Wait()
}

func (mpt *MarketPriceTasker) Stopped() bool {
	return mpt.stopped.Load()
}
