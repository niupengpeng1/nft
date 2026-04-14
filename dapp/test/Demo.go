package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Hello, World!")
	var contractAddress = common.HexToAddress("0x8a791620dd6260079bf849dc5567adc3f2fdc318")
	//连接ethers客户端
	clien := creatClient()

	//获取合约Json格式字符串，得到ABI
	paseAbi, err := getAbi()
	if err != nil {
		log.Fatalf("获取ABI失败 %v", err)
	}
	log.Print(paseAbi)

	query := ethereum.FilterQuery{Addresses: []common.Address{contractAddress}}
	logsCh := make(chan types.Log)

	subscription, err := clien.SubscribeFilterLogs(context.Background(), query, logsCh)

	if err != nil {
		log.Fatalf("合约订阅失败%v", err)
	}

	log.Printf("Subscribed to logs of contract %s  ,Rpcurl %s \n", contractAddress.Hex(), "指定地址")
	log.Print("***********开始监听***********************")

	for {
		select {
		case vLog := <-logsCh:
			parseLogEvent(&vLog, paseAbi)
		case err := <-subscription.Err():
			log.Printf("subscription error: %v", err)
			return
		}

	}
}

func parseLogEvent(vLog *types.Log, paseAbi abi.ABI) {
	//检查是否为有效日志
	if len(vLog.Topics) == 0 {
		log.Fatalf("无效日志")
	}

	//监测topics日志入参
	eventTopic := vLog.Topics[0]

	var eventName string
	var eventSig abi.Event
	for name, event := range paseAbi.Events {
		cryptoEventHash := crypto.Keccak256Hash([]byte(event.Sig))
		if eventTopic == cryptoEventHash {
			eventName = name
			eventSig = event
			break
		}
	}
	if eventName == "" {
		log.Printf("无法识别的日志:Block : %d. Tx: %s,Topic[0]: %s",
			vLog.BlockNumber, vLog.TxHash.Hex(), eventTopic.Hex())
	}

	//解析topic入参
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Printf("事件名称:%s \n", eventName)
	fmt.Printf("区块号: %d \n", vLog.BlockNumber)
	fmt.Printf("交易hash： %s \n", vLog.TxHash.Hex())
	fmt.Printf("log日志事件索引: %d \n", vLog.Index)
	fmt.Printf("合约地址: %s \n", vLog.Address.Hex())
	fmt.Printf("Topics count :%d \n", len(vLog.Topics))

	//topics 入参是从下指标1开始
	indexTopic := 0
	for i, input := range eventSig.Inputs {
		if !input.Indexed {
			continue
		}

		startINdex := indexTopic + 1
		indexTopic++
		fmt.Printf("[%d] %s (%s): \n", i+1, input.Name, input.Type)
		topic := vLog.Topics[startINdex]
		fmt.Print("入参值是：")
		switch input.Type.T {
		case abi.AddressTy:
			// address 类型：去除前 12 字节的 0 填充，后 20 字节是地址
			addr := common.BytesToAddress(topic.Bytes())
			fmt.Printf("%s\n", addr.Hex())
		case abi.IntTy, abi.UintTy:
			// 整数类型：直接转换为 big.Int
			value := new(big.Int).SetBytes(topic.Bytes())
			fmt.Printf("%s\n", value.String())
		case abi.BoolTy:
			// bool 类型：检查最后一个字节
			fmt.Printf("%t\n", topic[31] != 0)
		case abi.BytesTy, abi.FixedBytesTy:
			// bytes 类型：直接显示十六进制
			fmt.Printf("%s\n", topic.Hex())
		default:
			// 其他类型：显示原始十六进制
			fmt.Printf("%s (raw)\n", topic.Hex())
		}
	}

	// Data 字段包含所有非 indexed 参数的编码数据
	if len(vLog.Data) > 0 {
		fmt.Printf("\n  Non-Indexed Parameters (from Data):\n")

		// 创建一个结构体来接收解码后的参数
		// 注意：这里使用通用方法，实际应用中可能需要根据具体事件定义结构体
		nonIndexedInputs := make([]abi.Argument, 0)
		for _, input := range eventSig.Inputs {
			if !input.Indexed {
				nonIndexedInputs = append(nonIndexedInputs, input)
			}
		}

		if len(nonIndexedInputs) > 0 {
			// 使用 ABI 解码 Data 字段
			// 方法 1: 使用 UnpackIntoInterface（需要预定义结构体）
			// 方法 2: 使用 Unpack（返回 []interface{}）
			values, err := paseAbi.Unpack(eventName, vLog.Data)
			if err != nil {
				fmt.Printf("    Error decoding data: %v\n", err)
			} else {
				// 只输出非 indexed 参数
				nonIndexedIdx := 0
				for i, input := range eventSig.Inputs {
					if !input.Indexed {
						if nonIndexedIdx < len(values) {
							value := values[nonIndexedIdx]
							fmt.Printf("    [%d] %s (%s): ", i+1, input.Name, input.Type)

							// 根据类型格式化输出
							switch v := value.(type) {
							case *big.Int:
								fmt.Printf("%s\n", v.String())
							case common.Address:
								fmt.Printf("%s\n", v.Hex())
							case []byte:
								fmt.Printf("0x%x\n", v)
							default:
								fmt.Printf("%v\n", v)
							}
							nonIndexedIdx++
						}
					}
				}
			}
		}
	} else {
		fmt.Printf("\n  Non-Indexed Parameters: None\n")
	}
}

func getAbi() (abi.ABI, error) {
	var contractStr = `[
	
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
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "price",
          "type": "uint256"
        }
      ],
      "name": "updatePriceEvent",
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
  
	]`
	abi, err := abi.JSON(strings.NewReader(contractStr))
	if err != nil {
		log.Fatalf("解析ABI失败 %v", err)
	}
	return abi, nil
}

func creatClient() *ethclient.Client {
	rpcUrl := "ws://127.0.0.1:8545"

	tx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := ethclient.DialContext(tx, rpcUrl)
	if err != nil {
		log.Fatalf("创建客户端失败 %v", err)
	}

	return client
}
