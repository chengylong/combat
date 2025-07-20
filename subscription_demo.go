// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// func main() {
// 	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
// 	if err != nil {
// 		log.Fatal("连接失败:", err)
// 	}

// 	// 1. 获取当前状态
// 	currentBlock, err := client.BlockNumber(context.Background())
// 	if err != nil {
// 		log.Fatal("获取当前区块失败:", err)
// 	}

// 	fmt.Printf("=== 订阅开始 ===\n")
// 	fmt.Printf("当前区块号: %d\n", currentBlock)
// 	fmt.Printf("开始监听区块 #%d 之后的区块...\n", currentBlock+1)
// 	fmt.Printf("等待新区块产生...\n\n")

// 	// 2. 开始订阅
// 	headers := make(chan *types.Header)
// 	sub, err := client.SubscribeNewHead(context.Background(), headers)
// 	if err != nil {
// 		log.Fatal("订阅失败:", err)
// 	}
// 	defer sub.Unsubscribe()

// 	// 3. 监听新区块
// 	blockCount := 0
// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			log.Fatal("订阅错误:", err)

// 		case header := <-headers:
// 			blockCount++
// 			newBlockNumber := header.Number.Uint64()

// 			fmt.Printf("=== 收到新区块 #%d ===\n", newBlockNumber)
// 			fmt.Printf("区块哈希: %s\n", header.Hash().Hex())
// 			fmt.Printf("时间戳: %s\n", time.Unix(int64(header.Time), 0).Format("2006-01-02 15:04:05"))
// 			fmt.Printf("距离开始订阅: %d 个区块\n", newBlockNumber-currentBlock)
// 			fmt.Printf("已收到: %d 个新区块\n", blockCount)
// 			fmt.Printf("区块间隔: ~12秒\n\n")

// 			// 只显示前3个区块，然后退出
// 			if blockCount >= 3 {
// 				fmt.Printf("演示完成！\n")
// 				return
// 			}
// 		}
// 	}
// }
