package grpc

import (
	"context"

	"github.com/the-web3/rpc-server/proto/market"
)

func (mps *MarketRpcService) GetSupportSymbolList(ctx context.Context, in *market.SupportSymbolListReq) (*market.SupportSymbolListRep, error) {
	supportTokenList := mps.db.SupportToken.ListSupportToken()
	var tokenSymbolList []*market.TokenSymbol
	for _, sptItem := range supportTokenList {
		exchanges, _ := mps.db.SupportTokenExchange.QueryExchangesByTokenId(sptItem.GUID.String())
		var mExchangeList []*market.Exchange
		for _, exchange := range exchanges {
			mexItem := &market.Exchange{
				Name: exchange.Name,
			}
			mExchangeList = append(mExchangeList, mexItem)
		}
		tokenSymbolList = append(tokenSymbolList, &market.TokenSymbol{
			BaseAsset:  sptItem.BaseAsset,
			QouteAsset: sptItem.QouteAsset,
			Exchanges:  mExchangeList,
		})
	}
	return &market.SupportSymbolListRep{
		Code:         1,
		Msg:          "fetch suppport token list success",
		TokenSymbols: tokenSymbolList,
	}, nil
}

func (mps *MarketRpcService) GetMarketPriceList(ctx context.Context, in *market.MarketPriceListReq) (*market.MarketPriceListRep, error) {
	tokenPirces, err := mps.db.TokenPrice.QueryTokenPriceById()
	if err != nil {
		return nil, err
	}
	var mklist []*market.MarketPrice
	for _, tokenPriceItem := range tokenPirces {
		asset, _ := mps.db.SupportToken.QueryTokenNameByTokenId(tokenPriceItem.TokenId)
		mpPirceItem := &market.MarketPrice{
			BaseAsset:    asset.BaseAsset,
			QouteAsset:   asset.QouteAsset,
			SellPrice:    tokenPriceItem.AvgPrice,
			BuyPrice:     tokenPriceItem.AvgPrice,
			MarketPrice:  tokenPriceItem.AvgPrice,
			Rate:         "0.81",
			Hours24Trade: tokenPriceItem.Volume24h,
		}
		mklist = append(mklist, mpPirceItem)
	}
	return &market.MarketPriceListRep{
		Code:         200,
		Msg:          "fetch market price success",
		MarketPrices: mklist,
	}, nil
}
