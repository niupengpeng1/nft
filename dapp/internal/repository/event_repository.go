package repository

import (
	"dapp/internal/models"
	"errors"

	"gorm.io/gorm"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

// SaveEvent saves a contract event to database
func (r *EventRepository) SaveEvent(event *models.ContractEvent) error {
	result := DB.Create(event)
	return result.Error
}

// SaveEvents saves multiple contract events to database
func (r *EventRepository) SaveEvents(events []*models.ContractEvent) error {
	result := DB.CreateInBatches(events, 100)
	return result.Error
}

// GetEventByTxHash gets an event by transaction hash
func (r *EventRepository) GetEventByTxHash(txHash string) (*models.ContractEvent, error) {
	var event models.ContractEvent
	result := DB.Where("tx_hash = ?", txHash).First(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// GetEventsByBlockNumber gets events by block number
func (r *EventRepository) GetEventsByBlockNumber(blockNumber uint64) ([]*models.ContractEvent, error) {
	var events []*models.ContractEvent
	result := DB.Where("block_number = ?", blockNumber).Find(&events)
	return events, result.Error
}

// GetEventsByContractAddress gets events by contract address
func (r *EventRepository) GetEventsByContractAddress(address string) ([]*models.ContractEvent, error) {
	var events []*models.ContractEvent
	result := DB.Where("contract_address = ?", address).Find(&events)
	return events, result.Error
}

// UpdateBlockHeight updates the last processed block height
func (r *EventRepository) UpdateBlockHeight(height uint64) error {
	var bh models.BlockHeight
	result := DB.FirstOrCreate(&bh)
	if result.Error != nil {
		return result.Error
	}

	bh.Height = height
	return DB.Save(&bh).Error
}

// GetBlockHeight gets the last processed block height
func (r *EventRepository) GetBlockHeight() (uint64, error) {
	var bh models.BlockHeight
	result := DB.First(&bh)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return bh.Height, nil
}
