package model

type Exchange struct {
	Name string
}

type SupportAssetExchange struct {
	BaseAsset  string     `json:"base_asset"`
	QouteAsset string     `json:"qoute_asset"`
	Exchange   []Exchange `json:"exchange"`
}

type SupportAssetExchangeResponse struct {
	ReturnCode           uint64                 `json:"return_code"`
	Message              string                 `json:"message"`
	SupportAssetExchange []SupportAssetExchange `json:"support_asset_exchange"`
}

type MarketPrice struct {
	AssetName   string `json:"asset_name"`
	AssetPrice  string `json:"asset_price"`
	AssetVolume string `json:"asset_volume"`
	AssetRate   string `json:"asset_rate"`
}

type MarketPriceResponse struct {
	ReturnCode      uint64        `json:"return_code"`
	Message         string        `json:"message"`
	MarketPriceList []MarketPrice `json:"market_price_list"`
}
