package database

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/the-web3/rpc-server/services/rest/model"
	"gorm.io/gorm"
)

type SupportTokenExchange struct {
	GUID       uuid.UUID `gorm:"primaryKey" json:"guid"`
	TokenId    string    `gorm:"token_id" json:"token_id"`
	ExchangeId string    `gorm:"exchange_id" json:"exchange_id"`
}

func (SupportTokenExchange) TableName() string {
	return "support_token_exchange"
}

type SupportTokenExchangeView interface {
	QueryExchangesByTokenId(string) ([]model.Exchange, error)
}

type SupportTokenExchangeDB interface {
	StoreSupportTokenExchange(tokenExchangeList []SupportTokenExchange) error
	SupportTokenExchangeView
}

type supportTokenExchangeDB struct {
	gorm *gorm.DB
}

func (st *supportTokenExchangeDB) StoreSupportTokenExchange(tokenExchangeList []SupportTokenExchange) error {
	result := st.gorm.Table("support_token_exchange").CreateInBatches(&tokenExchangeList, len(tokenExchangeList))
	return result.Error
}

func (st *supportTokenExchangeDB) QueryExchangesByTokenId(tokenId string) ([]model.Exchange, error) {
	var steList []SupportTokenExchange
	result := st.gorm.Table("support_token_exchange").Where(SupportTokenExchange{TokenId: tokenId}).Find(&steList)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error("query support token exchange record not found", "err", result.Error)
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	log.Info("query support exchange success", "steList", steList)

	var exchangeList []model.Exchange
	for _, steItem := range steList {
		var exchangeItem Exchange
		exchangeResult := st.gorm.Table("exchange").Where("guid=?", steItem.ExchangeId).Take(&exchangeItem)
		if result.Error != nil {
			if errors.Is(exchangeResult.Error, gorm.ErrRecordNotFound) {
				log.Error("query exchange record not found", "err", exchangeResult.Error)
				return nil, nil
			} else {
				return nil, exchangeResult.Error
			}
		}
		exchangeList = append(exchangeList, model.Exchange{Name: exchangeItem.Name})
	}
	return exchangeList, nil
}

func NewSupportTokenExchangeDB(db *gorm.DB) SupportTokenExchangeDB {
	return &supportTokenExchangeDB{gorm: db}
}
