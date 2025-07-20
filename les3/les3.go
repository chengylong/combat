package les3

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// 查询收据
func QueryReceipt() {
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	//使用区块hash查询
	receiptByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHash, false))
	if err != nil {
		log.Fatal(err)
	}

	//使用区块高度 blockNumber 查询
	blockNumber := big.NewInt(5671744)
	receiptsByNum, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(receiptByHash[0])
	// fmt.Println(receiptsByNum[0])
	fmt.Println(receiptByHash[0] == receiptsByNum[0]) // true

	for _, receipt := range receiptByHash {
		fmt.Println(receipt.Status)           // 1
		fmt.Println(receipt.Logs)             // []
		fmt.Println(receipt.TxHash.Hex())     // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
		fmt.Println(receipt.TransactionIndex) // 0
	}
	block, err := client.BlockByNumber(context.Background(), blockNumber)

	//通过区块中某个交易的交易 hash获取收据
	if err != nil {
		log.Fatal(err)
	}
	tx := block.Transactions()[0]
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(receipt.Status) // 1
	fmt.Println(receipt.Logs)   // ...

}
