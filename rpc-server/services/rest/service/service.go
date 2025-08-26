package service

import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/rpc-server/database"
	"github.com/the-web3/rpc-server/services/rest/model"
)

type RestService interface {
	GetSupportSymbolList() (*model.SupportAssetExchangeResponse, error)
	GetMarketPrice() (*model.MarketPriceResponse, error)
}

type HandleSvc struct {
	v                        *Validator
	supportTokenView         database.SupportTokenView
	tokenPriceView           database.TokenPriceView
	supportTokenExchangeView database.SupportTokenExchangeView
}

func NewHandleSvc(v *Validator, supportTokenView database.SupportTokenView, tokenPriceView database.TokenPriceView, supportTokenExchangeView database.SupportTokenExchangeView) RestService {
	return &HandleSvc{
		v:                        v,
		supportTokenView:         supportTokenView,
		tokenPriceView:           tokenPriceView,
		supportTokenExchangeView: supportTokenExchangeView,
	}
}

func (h HandleSvc) GetSupportSymbolList() (*model.SupportAssetExchangeResponse, error) {
	supportTokenList := h.supportTokenView.ListSupportToken()
	var saeList []model.SupportAssetExchange
	for _, supportEx := range supportTokenList {
		exchanges, err := h.supportTokenExchangeView.QueryExchangesByTokenId(supportEx.GUID.String())
		if err != nil {
			log.Error("query support token list fail", "err", err)
			return nil, err
		}
		log.Info("query exchangs success", "exchanges", exchanges)
		saeItem := model.SupportAssetExchange{
			BaseAsset:  supportEx.BaseAsset,
			QouteAsset: supportEx.QouteAsset,
			Exchange:   exchanges,
		}
		saeList = append(saeList, saeItem)
	}
	return &model.SupportAssetExchangeResponse{
		ReturnCode:           100,
		Message:              "get support asset success",
		SupportAssetExchange: saeList,
	}, nil
}

func (h HandleSvc) GetMarketPrice() (*model.MarketPriceResponse, error) {
	priceList, err := h.tokenPriceView.QueryTokenPriceById()
	if err != nil {
		log.Error("query token price error", "err", err)
		return nil, err
	}

	var marketPriceList []model.MarketPrice

	for _, priceItem := range priceList {
		sToken, err := h.supportTokenView.QueryTokenNameByTokenId(priceItem.TokenId)
		if err != nil {
			log.Info("query token name by tokenId fail", "err", err)
			return nil, err
		}
		mp := model.MarketPrice{
			AssetName:   sToken.SymbolName,
			AssetPrice:  priceItem.AvgPrice,
			AssetVolume: priceItem.Volume24h,
			AssetRate:   priceItem.Rate,
		}

		marketPriceList = append(marketPriceList, mp)
	}

	return &model.MarketPriceResponse{
		ReturnCode:      100,
		Message:         "get market price success",
		MarketPriceList: marketPriceList,
	}, nil
}
