package market_price

import (
	"fmt"

	"github.com/pkg/errors"

	gresty "github.com/go-resty/resty/v2"
)

var errExchangeHTTPError = errors.New("Exchange Market Price Http Error")

type Data struct {
	BaseAsset      string  `json:"base_asset"`
	Exchange       string  `json:"exchange"`
	Price          float64 `json:"price"`
	PriceChange24h float64 `json:"price_change_24h"`
	Volume24h      float64 `json:"volume_24h"`
	Timestamp      string  `json:"timestamp"`
}

type MetaData struct {
	Ok     bool  `json:"ok"`
	Code   int   `json:"code"`
	Result *Data `json:"result"`
}

type ExchangeClient interface {
	GetMarketPriceByBaseAsset(baseAsset string) (*Data, error)
}

type Client struct {
	client *gresty.Client
}

func NewExchangeClient(baseUrl string) (*Client, error) {
	client := gresty.New()
	client.SetBaseURL(baseUrl)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errExchangeHTTPError)
		}
		return nil
	})
	return &Client{
		client: client,
	}, nil
}

// GetMarketPriceByBaseAsset dapplink 的行情业务模块
func (c *Client) GetMarketPriceByBaseAsset(baseAsset string) (*Data, error) {
	var metaData MetaData
	response, err := c.client.R().
		SetResult(&metaData).
		SetQueryParam("symbol", baseAsset).
		Get("/api/v1/ccxt/price")
	if err != nil {
		return nil, fmt.Errorf("cannot market price fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get market price fail")
	}
	if !metaData.Ok && metaData.Result == nil {
		return nil, errors.New("fetch meta data price fail")
	}
	return metaData.Result, nil
}
