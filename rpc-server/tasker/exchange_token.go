package tasker

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/rpc-server/common/retry"
	"github.com/the-web3/rpc-server/database"
)

type MarketAssetAndExchange struct {
	db             *database.DB
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
}

func NewMarketAssetAndExchange(db *database.DB) (*MarketAssetAndExchange, error) {
	resCtx, resCancel := context.WithCancel(context.Background())

	return &MarketAssetAndExchange{
		db:             db,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
	}, nil
}

func (mae *MarketAssetAndExchange) ConfigAssetAndExchange() error {
	var exchangeList []database.Exchange

	binanceUuid := uuid.New()
	bybitUuid := uuid.New()
	dapplinkUuid := uuid.New()

	btcUsdtUuid := uuid.New()
	ethUsdtUuid := uuid.New()
	solUsdtUuid := uuid.New()

	exchangeList = append(exchangeList, database.Exchange{
		GUID:      binanceUuid,
		Name:      "binance",
		Config:    "api-key",
		Timestamp: uint64(time.Now().Unix()),
	})
	exchangeList = append(exchangeList, database.Exchange{
		GUID:      bybitUuid,
		Name:      "bybit",
		Config:    "api-key",
		Timestamp: uint64(time.Now().Unix()),
	})
	exchangeList = append(exchangeList, database.Exchange{
		GUID:      dapplinkUuid,
		Name:      "dapplink",
		Config:    "api-key",
		Timestamp: uint64(time.Now().Unix()),
	})

	var supportTokenList []database.SupportToken
	supportTokenList = append(supportTokenList, database.SupportToken{
		GUID:       btcUsdtUuid,
		SymbolName: "BTC/USDT",
		BaseAsset:  "BTC",
		QouteAsset: "USDT",
		Timestamp:  uint64(time.Now().Unix()),
	})
	supportTokenList = append(supportTokenList, database.SupportToken{
		GUID:       ethUsdtUuid,
		SymbolName: "ETH/USDT",
		BaseAsset:  "ETH",
		QouteAsset: "USDT",
		Timestamp:  uint64(time.Now().Unix()),
	})
	supportTokenList = append(supportTokenList, database.SupportToken{
		GUID:       solUsdtUuid,
		SymbolName: "SOL/USDT",
		BaseAsset:  "SOL",
		QouteAsset: "USDT",
		Timestamp:  uint64(time.Now().Unix()),
	})

	var SupportTokenExchangeList []database.SupportTokenExchange
	SupportTokenExchangeList = append(SupportTokenExchangeList, database.SupportTokenExchange{
		GUID:       uuid.New(),
		TokenId:    btcUsdtUuid.String(),
		ExchangeId: binanceUuid.String(),
	})

	SupportTokenExchangeList = append(SupportTokenExchangeList, database.SupportTokenExchange{
		GUID:       uuid.New(),
		TokenId:    btcUsdtUuid.String(),
		ExchangeId: bybitUuid.String(),
	})

	SupportTokenExchangeList = append(SupportTokenExchangeList, database.SupportTokenExchange{
		GUID:       uuid.New(),
		TokenId:    ethUsdtUuid.String(),
		ExchangeId: bybitUuid.String(),
	})

	SupportTokenExchangeList = append(SupportTokenExchangeList, database.SupportTokenExchange{
		GUID:       uuid.New(),
		TokenId:    ethUsdtUuid.String(),
		ExchangeId: dapplinkUuid.String(),
	})

	SupportTokenExchangeList = append(SupportTokenExchangeList, database.SupportTokenExchange{
		GUID:       uuid.New(),
		TokenId:    solUsdtUuid.String(),
		ExchangeId: binanceUuid.String(),
	})

	SupportTokenExchangeList = append(SupportTokenExchangeList, database.SupportTokenExchange{
		GUID:       uuid.New(),
		TokenId:    solUsdtUuid.String(),
		ExchangeId: dapplinkUuid.String(),
	})

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](mae.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := mae.db.Transaction(func(tx *database.DB) error {

			if len(SupportTokenExchangeList) > 0 {
				err := mae.db.SupportTokenExchange.StoreSupportTokenExchange(SupportTokenExchangeList)
				if err != nil {
					log.Error("store support token exchange fail", "err", err)
					return err
				}
			}

			if len(supportTokenList) > 0 {
				err := mae.db.SupportToken.StoreSupportToken(supportTokenList)
				if err != nil {
					log.Error("store support token fail", "err", err)
					return err
				}
			}

			if len(exchangeList) > 0 {
				err := mae.db.Exchange.StoreExchanges(exchangeList)
				if err != nil {
					log.Error("store exchanges error", "err", err)
					return err
				}
			}
			return nil
		}); err != nil {
			log.Info("unable to persist batch", err)
			return nil, fmt.Errorf("unable to persist batch: %w", err)
		}
		return nil, nil
	}); err != nil {
		return err
	}

	return nil
}
