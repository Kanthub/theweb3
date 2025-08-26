package database

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/log"
)

type SupportToken struct {
	GUID       uuid.UUID `gorm:"primaryKey" json:"guid"`
	SymbolName string    `gorm:"symbol_name" json:"symbol_name"`
	BaseAsset  string    `gorm:"base_asset" json:"base_asset"`
	QouteAsset string    `gorm:"qoute_asset" json:"qoute_asset"`
	Timestamp  uint64    `json:"timestamp"`
}

func (SupportToken) TableName() string {
	return "support_token"
}

type SupportTokenView interface {
	QueryTokenNameByTokenId(tokenId string) (*SupportToken, error)
	ListSupportToken() []SupportToken
}

type SupportTokenDB interface {
	SupportTokenView
	StoreSupportToken([]SupportToken) error
}

type supportTokenDB struct {
	gorm *gorm.DB
}

func NewSupportTokenDB(db *gorm.DB) SupportTokenDB {
	return &supportTokenDB{gorm: db}
}

func (st *supportTokenDB) ListSupportToken() []SupportToken {
	var supportTokenList []SupportToken
	qErr := st.gorm.Table("support_token").Find(&supportTokenList).Error
	if qErr != nil {
		log.Error("list supportTokenList fail", "err", qErr)
	}
	return supportTokenList
}

func (st *supportTokenDB) StoreSupportToken(tokenList []SupportToken) error {
	result := st.gorm.Table("support_token").CreateInBatches(&tokenList, len(tokenList))
	return result.Error
}

func (st *supportTokenDB) QueryTokenNameByTokenId(tokenId string) (*SupportToken, error) {
	var supportToken SupportToken
	result := st.gorm.Table("support_token").Where("guid = ?", tokenId).Take(&supportToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return &supportToken, nil
}
