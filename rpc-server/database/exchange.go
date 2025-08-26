package database

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Exchange struct {
	GUID      uuid.UUID `gorm:"primaryKey" json:"guid"`
	Name      string    `gorm:"name" json:"name"`
	Config    string    `gorm:"config" json:"config"`
	Timestamp uint64    `json:"timestamp"`
}

func (Exchange) TableName() string {
	return "exchange"
}

type ExchangeView interface {
	QueryExchangeGuid(string) (*Exchange, error)
}

type ExchangeDB interface {
	ExchangeView
	StoreExchanges([]Exchange) error
}

type exchangeDB struct {
	gorm *gorm.DB
}

func NewExchangeDB(db *gorm.DB) ExchangeDB {
	return &exchangeDB{gorm: db}
}

func (exd *exchangeDB) QueryExchangeGuid(guid string) (*Exchange, error) {
	var exchangeItem Exchange
	result := exd.gorm.Table("exchange").Where("guid=?", guid).Take(&exchangeItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return &exchangeItem, nil
}

func (exd *exchangeDB) StoreExchanges(exchangeList []Exchange) error {
	result := exd.gorm.Table("exchange").CreateInBatches(&exchangeList, len(exchangeList))
	return result.Error
}
