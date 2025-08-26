package worker

import (
	"errors"
	common2 "github.com/the-web3-contracts/vrf-node/database/common"
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/the-web3/dapplink-vrf/database/utils/serializers"
)

type EventBlocks struct {
	GUID       uuid.UUID   `gorm:"primaryKey"`
	Hash       common.Hash `gorm:"serializer:bytes"`
	ParentHash common.Hash `gorm:"serializer:bytes"`
	Number     *big.Int    `gorm:"serializer:u256"`
	Timestamp  uint64
}

type BlocksView interface {
	LatestEventBlockHeader() (*common2.BlockHeader, error)
}

type EventBlocksDB interface {
	BlocksView
	StoreEventBlocks([]EventBlocks) error
}

type evnetBlocksDB struct {
	gorm *gorm.DB
}

func (e evnetBlocksDB) LatestEventBlockHeader() (*common2.BlockHeader, error) {
	eventQuery := e.gorm.Where("number = (?)", e.gorm.Table("event_blocks").Select("MAX(number)"))
	var eventBlock common2.BlockHeader
	result := eventQuery.Take(&eventBlock)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &eventBlock, nil
}

func (e evnetBlocksDB) StoreEventBlocks(eventBlocks []EventBlocks) error {
	result := e.gorm.CreateInBatches(&eventBlocks, len(eventBlocks))
	return result.Error
}

func NewEventBlocksDB(db *gorm.DB) EventBlocksDB {
	return &evnetBlocksDB{gorm: db}
}
