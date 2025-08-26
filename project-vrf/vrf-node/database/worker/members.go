package worker

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	_ "github.com/the-web3/dapplink-vrf/database/utils/serializers"
)

type Members struct {
	GUID      uuid.UUID `gorm:"primaryKey" json:"guid"`
	Member    string    `json:"member"`
	IsActive  uint8     `json:"is_active"`
	Timestamp uint64
}

type MembersView interface {
	GetActiveMember() ([]string, error)
}

type MembersDB interface {
	MembersView
	StoreMembers([]Members) error
}

type membersDB struct {
	gorm *gorm.DB
}

func NewMembersDB(db *gorm.DB) MembersDB {
	return &membersDB{gorm: db}
}

func (mdb *membersDB) StoreMembers(members []Members) error {
	panic("implement me")
}

func (mdb *membersDB) GetActiveMember() ([]string, error) {
	var memberList []string
	return memberList, nil
}
