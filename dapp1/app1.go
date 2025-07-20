package dapp1

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 任务1，查询区块，发送交易
// 查询区块
func QueryBlock(blockNumber int64) {
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块哈希:", block.Hash().Hex())
	fmt.Println("区块时间戳:", block.Time())
	fmt.Println("区块交易数量:", len(block.Transactions()))
}

// ETH转账
func EthTra() {
	//连接客户端
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	// 加载私钥
	privateKey, err := crypto.HexToECDSA("private_key")
	if err != nil {
		log.Fatal(err)
	}
	// 先根据私钥获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 再根据公钥地址获取到自增数nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 转移eth数量
	value := big.NewInt(1000000000000000) // in wei (0.001 eth)
	// ETH 转账的燃气应设上限为“21000”单位。
	gasLimit := uint64(21000) // in units
	// 燃气价格必须以 wei 为单位设定。 在撰写本文时，将在一个区块中比较快的打包交易的燃气价格为 30 gwei。
	// gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	// 然而，燃气价格总是根据市场需求和用户愿意支付的价格而波动的，因此对燃气价格进行硬编码有时并不理想。
	// go-ethereum 客户端提供 SuggestGasPrice 函数，用于根据'x'个先前块来获得平均燃气价格。

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 发送的目的地址
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	// 现在我们最终可以通过导入 go-ethereum core/types 包并调用 NewTransaction 来生成我们的未签名以太坊事务，
	// 这个函数需要接收 nonce，地址，值，燃气上限值，燃气价格和可选发的数据。 发送 ETH 的数据字段为“nil”。 在与智能合约进行交互时，我们将使用数据字段，仅仅转账以太币是不需要数据字段的。
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 下一步是使用发件人的私钥对事务进行签名。 为此，我们调用 SignTx 方法，
	// 该方法接受一个未签名的事务和我们之前构造的私钥。 SignTx 方法需要 EIP155 签名者，这个也需要我们先从客户端拿到链 ID。
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 现在通过在 client 实例调用 SendTransaction 来将已签名的事务广播到整个网络。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0x77006fcb3938f648e2cc65bafd27dec30b9bfbe9df41f78498b9c8b7322a249e

}
