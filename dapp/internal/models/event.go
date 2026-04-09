package models

import (
	"time"
)

// ContractEvent represents a blockchain contract event
type ContractEvent struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	EventName       string    `gorm:"size:100;not null;index" json:"event_name"`
	ContractAddress string    `gorm:"size:42;not null;index" json:"contract_address"`
	TxHash          string    `gorm:"size:66;not null;" json:"tx_hash"`
	BlockNumber     uint64    `gorm:"not null;index" json:"block_number"`
	BlockHash       string    `gorm:"size:66" json:"block_hash"`
	LogIndex        uint      `gorm:"not null" json:"log_index"`
	FromAddress     string    `gorm:"size:600" json:"from_address"`
	ToAddress       string    `gorm:"size:42" json:"to_address"`
	Data            string    `gorm:"type:text" json:"data"`
	Topics          string    `gorm:"type:text" json:"topics"`
	Timestamp       time.Time `gorm:"autoCreateTime" json:"timestamp"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BlockHeight tracks the last processed block
type BlockHeight struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Height    uint64    `gorm:"not null" json:"height"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ContractEvent) TableName() string {
	return "contract_events"
}

func (BlockHeight) TableName() string {
	return "block_heights"
}
