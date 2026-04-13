package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Hello, World!")
	var contractAddress = common.HexToAddress("0x<your-contract-address>")
	//连接ethers客户端
	clien := creatClient()

	//获取合约Json格式字符串，得到ABI
	abi, err := getAbi()
	if err != nil {
		log.Fatalf("获取ABI失败 %v", err)
	}
	log.Print(abi)

	query := ethereum.FilterQuery{Addresses: []common.Address{contractAddress}}
	logsCh := make(chan types.Log)

	clien.SubscribeFilterLogs(context.Background(), query, logsCh)
}

func getAbi() (abi.ABI, error) {
	var contractStr = ``
	abi, err := abi.JSON(strings.NewReader(contractStr))
	if err != nil {
		log.Fatalf("解析ABI失败 %v", err)
	}
	return abi, nil
}

func creatClient() *ethclient.Client {
	rpcUrl := ""

	tx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := ethclient.DialContext(tx, rpcUrl)
	if err != nil {
		log.Fatalf("创建客户端失败 %v", err)
	}

	return client
}
