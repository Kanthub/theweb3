package database

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenPrice struct {
	GUID      uuid.UUID `json:"guid" gorm:"primaryKey;DEFAULT replace(uuid_generate_v4()::text,'-','')"`
	TokenId   string    `gorm:"token_id" json:"token_id"`
	SellPrice string    `gorm:"sell_price" json:"sell_price"`
	Rate      string    `gorm:"rate" json:"rate"`
	AvgPrice  string    `gorm:"avg_price" json:"avg_price"`
	Volume24h string    `gorm:"volume24h" json:"volume24h"`
	Timestamp uint64    `json:"timestamp"`
}

func (TokenPrice) TableName() string {
	return "token_price"
}

type TokenPriceView interface {
	QueryTokenPriceById() ([]TokenPrice, error)
}

type TokenPriceDB interface {
	TokenPriceView
	StoreAndUpdateTokenPrices([]TokenPrice) error
}

type tokenPriceDB struct {
	gorm *gorm.DB
}

func NewTokenPriceDB(db *gorm.DB) TokenPriceDB {
	return &tokenPriceDB{gorm: db}
}

func (tpb *tokenPriceDB) QueryTokenPriceById() ([]TokenPrice, error) {
	var tokenPriceList []TokenPrice
	result := tpb.gorm.Table("token_price").Find(&tokenPriceList)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return tokenPriceList, nil
}

func (tpb *tokenPriceDB) StoreAndUpdateTokenPrices(tokenPriceList []TokenPrice) error {
	for _, tokenItem := range tokenPriceList {
		var tokenPrice TokenPrice
		result := tpb.gorm.Model(&TokenPrice{}).Where("token_id = ?", tokenItem.TokenId).First(&tokenPrice)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := tpb.gorm.Create(&tokenItem).Error; err != nil {
					log.Error("store token price fail", "tokenId", tokenItem.TokenId, "err", err)
					return err
				}
			} else {
				log.Error("query token price fail", "tokenId", tokenItem.TokenId, "err", result.Error)
				return result.Error
			}
		} else {
			updates := map[string]interface{}{
				"sell_price": tokenItem.SellPrice,
				"avg_price":  tokenItem.AvgPrice,
				"rate":       tokenItem.Rate,
			}
			if err := tpb.gorm.Model(&TokenPrice{}).Where("token_id = ?", tokenItem.TokenId).Updates(updates).Error; err != nil {
				log.Error("update token price fail", "tokenId", tokenItem.TokenId, "err", err)
				return err
			}
		}
	}
	return nil
}
