package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"dapp/internal/config"
	"dapp/internal/models"
	"dapp/internal/repository"
	"dapp/pkg/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventListenerService struct {
	client         *ethclient.Client
	config         *config.Web3Config
	eventRepo      *repository.EventRepository
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
	isRunning      bool
	lastBlock      uint64
	checkInterval  time.Duration
}

func NewEventListenerService(cfg *config.Web3Config) *EventListenerService {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventListenerService{
		config:        cfg,
		eventRepo:     repository.NewEventRepository(),
		ctx:           ctx,
		cancel:        cancel,
		checkInterval: 12 * time.Second, // Approximately one Ethereum block time
	}
}

// Init initializes the Web3 client
func (s *EventListenerService) Init() error {
	client, err := ethclient.Dial(s.config.RPCURL)
	if err != nil {
		return fmt.Errorf("failed to connect to ethereum node: %w", err)
	}
	s.client = client
	logger.Info("Connected to Ethereum node: %s", s.config.RPCURL)
	return nil
}

// Close closes the Web3 client connection
func (s *EventListenerService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

// Start starts listening for contract events
func (s *EventListenerService) Start() error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("event listener is already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	// Load last processed block
	var err error
	s.lastBlock, err = s.eventRepo.GetBlockHeight()
	if err != nil {
		logger.Warn("Failed to get last block height, using config start block: %v", err)
		s.lastBlock = s.config.StartBlock
	}

	logger.Info("Starting event listener from block %d", s.lastBlock)

	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			logger.Info("Event listener stopped")
			return nil
		case <-ticker.C:
			if err := s.processNewBlocks(); err != nil {
				logger.Error("Error processing blocks: %v", err)
			}
		}
	}
}

// Stop stops the event listener
func (s *EventListenerService) Stop() {
	s.cancel()
	s.mu.Lock()
	s.isRunning = false
	s.mu.Unlock()
	logger.Info("Stopping event listener...")
}

// processNewBlocks checks and processes new blocks
func (s *EventListenerService) processNewBlocks() error {
	currentBlock, err := s.client.BlockNumber(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	if currentBlock <= s.lastBlock {
		return nil
	}

	logger.Debug("Processing blocks from %d to %d", s.lastBlock+1, currentBlock)

	// Process blocks in batches
	batchSize := uint64(100)
	for from := s.lastBlock + 1; from <= currentBlock; from += batchSize {
		to := from + batchSize - 1
		if to > currentBlock {
			to = currentBlock
		}

		if err := s.processBlockRange(from, to); err != nil {
			logger.Error("Failed to process block range %d-%d: %v", from, to, err)
			continue
		}

		s.lastBlock = to
		if err := s.eventRepo.UpdateBlockHeight(to); err != nil {
			logger.Error("Failed to update block height: %v", err)
		}
	}

	return nil
}

// processBlockRange processes a range of blocks
func (s *EventListenerService) processBlockRange(from, to uint64) error {
	contractAddress := common.HexToAddress(s.config.ContractAddr)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: []common.Address{contractAddress},
	}

	logs, err := s.client.FilterLogs(s.ctx, query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	if len(logs) == 0 {
		return nil
	}

	logger.Info("Found %d events in blocks %d-%d", len(logs), from, to)

	events := make([]*models.ContractEvent, 0, len(logs))
	for _, vLog := range logs {
		event := s.logToContractEvent(&vLog)
		events = append(events, event)
	}

	if err := s.eventRepo.SaveEvents(events); err != nil {
		return fmt.Errorf("failed to save events: %w", err)
	}

	logger.Info("Saved %d events to database", len(events))
	return nil
}

// logToContractEvent converts a types.Log to ContractEvent model
func (s *EventListenerService) logToContractEvent(vLog *types.Log) *models.ContractEvent {
	topicsJSON, _ := json.Marshal(vLog.Topics)
	
	event := &models.ContractEvent{
		EventName:       "Unknown", // Can be enhanced by parsing topic[0]
		ContractAddress: vLog.Address.Hex(),
		TxHash:          vLog.TxHash.Hex(),
		BlockNumber:     vLog.BlockNumber,
		BlockHash:       vLog.BlockHash.Hex(),
		LogIndex:        uint(vLog.Index),
		Data:            common.Bytes2Hex(vLog.Data),
		Topics:          string(topicsJSON),
	}

	// Extract from and to addresses if available in topics
	if len(vLog.Topics) > 1 {
		// Common pattern: first topic is event signature, second is often from address
		event.FromAddress = common.BytesToAddress(vLog.Topics[1].Bytes()).Hex()
	}

	return event
}

// GetStatus returns the current status of the event listener
func (s *EventListenerService) GetStatus() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	status := map[string]interface{}{
		"is_running":    s.isRunning,
		"last_block":    s.lastBlock,
		"contract_addr": s.config.ContractAddr,
		"rpc_url":       s.config.RPCURL,
	}

	if s.client != nil {
		if blockNumber, err := s.client.BlockNumber(context.Background()); err == nil {
			status["current_block"] = blockNumber
		}
	}

	return status
}
