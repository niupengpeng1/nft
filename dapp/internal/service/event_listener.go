package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"dapp/internal/config"
	"dapp/internal/models"
	"dapp/internal/repository"
	"dapp/pkg/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventListenerService struct {
	client        *ethclient.Client
	config        *config.Web3Config
	eventRepo     *repository.EventRepository
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.Mutex
	isRunning     bool
	lastBlock     uint64
	checkInterval time.Duration
}

const myNFTABI = `[
    {
      "inputs": [],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "ERC721IncorrectOwner",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "ERC721InsufficientApproval",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "approver",
          "type": "address"
        }
      ],
      "name": "ERC721InvalidApprover",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        }
      ],
      "name": "ERC721InvalidOperator",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "ERC721InvalidOwner",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        }
      ],
      "name": "ERC721InvalidReceiver",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        }
      ],
      "name": "ERC721InvalidSender",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "ERC721NonexistentToken",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnableInvalidOwner",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "OwnableUnauthorizedAccount",
      "type": "error"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "approved",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "ApprovalForAll",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_fromTokenId",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_toTokenId",
          "type": "uint256"
        }
      ],
      "name": "BatchMetadataUpdate",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_tokenId",
          "type": "uint256"
        }
      ],
      "name": "MetadataUpdate",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Mint",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "getApproved",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        }
      ],
      "name": "isApprovedForAll",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "url",
          "type": "string"
        }
      ],
      "name": "mint",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "ownerOf",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "price",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "safeTransferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "data",
          "type": "bytes"
        }
      ],
      "name": "safeTransferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "setApprovalForAll",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes4",
          "name": "interfaceId",
          "type": "bytes4"
        }
      ],
      "name": "supportsInterface",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "tokenURI",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_price",
          "type": "uint256"
        }
      ],
      "name": "updatePrice",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "withdraw",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]

`

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
	if len(vLog.Topics) == 0 {
		log.Fatalf("无效日志")
	}
	eventTopic := vLog.Topics[0]
	var eventName string
	var eventSig abi.Event

	topicsJSON, _ := json.Marshal(vLog.Topics)
	parsedABI, err := abi.JSON(strings.NewReader(myNFTABI))

	if err != nil {
		log.Fatalf("failed to parse ABI: %v", err)
	}

	// 遍历 ABI 中定义的所有事件，查找匹配的事件签名
	for name, event := range parsedABI.Events {
		// 计算事件的签名哈希
		eventSigHash := crypto.Keccak256Hash([]byte(event.Sig))
		if eventSigHash == eventTopic {
			eventName = name
			eventSig = event
			break
		}
	}
	log.Print(eventSig)

	// Topics[0] 是事件签名，所以 indexed 参数从 Topics[1] 开始
	// 注意：topicIndex 只针对 indexed 参数计数，不考虑非 indexed 参数
	indexedParamIndex := 0
	var topicInputs string
	for i, input := range eventSig.Inputs {
		if !input.Indexed {
			continue
		}
		// indexed 参数在 Topics 中的位置 = 1 + indexed 参数的索引
		topicIndex := 1 + indexedParamIndex
		indexedParamIndex++

		if topicIndex >= len(vLog.Topics) {
			continue
		}

		topic := vLog.Topics[topicIndex]
		fmt.Printf("    [%d] %s (%s): ", i+1, input.Name, input.Type)

		// 根据类型解析 indexed 参数
		switch input.Type.T {
		case abi.AddressTy:
			// address 类型：去除前 12 字节的 0 填充，后 20 字节是地址
			addr := common.BytesToAddress(topic.Bytes())
			topicInputs += fmt.Sprintf("(%s: %s\n) ", input.Name, addr.Hex())
			fmt.Printf("%s\n", addr.Hex())
		case abi.IntTy, abi.UintTy:
			// 整数类型：直接转换为 big.Int
			value := new(big.Int).SetBytes(topic.Bytes())
			topicInputs += fmt.Sprintf("(%s: %s\n) ", input.Name, value.String())
			fmt.Printf("%s\n", value.String())
		case abi.BoolTy:
			// bool 类型：检查最后一个字节
			topicInputs += fmt.Sprintf("(%s: %t\n) ", input.Name, topic[31] != 0)
			fmt.Printf("%t\n", topic[31] != 0)
		case abi.BytesTy, abi.FixedBytesTy:
			topicInputs += fmt.Sprintf("(%s: %s\n) ", input.Name, topic.Hex())
			// bytes 类型：直接显示十六进制
			fmt.Printf("%s\n", topic.Hex())
		default:
			topicInputs += fmt.Sprintf("(%s: %s\n) ", input.Name, topic.Hex())
			// 其他类型：显示原始十六进制
			fmt.Printf("%s (raw)\n", topic.Hex())
		}
	}
	event := &models.ContractEvent{
		EventName:       eventName, // Can be enhanced by parsing topic[0]
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
		event.FromAddress = topicInputs
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
