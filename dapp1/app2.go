package dapp1

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/chengylong/combat/counter"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 任务2合约代码生成 任务目标

// 部署计数器合约
func DeployContract() {
	// 1. 连接到 Sepolia 测试网
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	// 2. 加载私钥（Metamask 导出的私钥，去掉0x前缀）
	privateKey, err := crypto.HexToECDSA("private key")
	if err != nil {
		log.Fatal(err)
	}

	// 3. 获取部署账户地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("部署账户地址:", fromAddress.Hex())

	// 4. 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 5. 获取 gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 6. 获取 chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 7. 构造授权对象
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // 部署合约不需要转账
	auth.GasLimit = uint64(3000000) // 建议300万
	auth.GasPrice = gasPrice

	// 8. 部署合约（Counter 没有构造参数，如果有参数，写在后面）
	address, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("合约部署中，等待上链...")
	fmt.Println("合约地址:", address.Hex())
	fmt.Println("部署交易哈希:", tx.Hash().Hex())

	// 9. 等待合约部署上链
	fmt.Println("等待合约部署上链...")
	bind.WaitMined(context.Background(), client, tx)
	fmt.Println("合约已上链！")

	// 10. 测试调用合约方法
	count, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("初始计数器值:", count)
	// // 让计数器加1
	// tx, err = instance.Increment(auth) // 需要签名和gas
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("计数器自增:", count)
	// fmt.Println("tx:", tx)
	// // 等待自增交易上链
	// bind.WaitMined(context.Background(), client, tx)

	// // 查询当前计数器值
	// count, err = instance.GetCount(&bind.CallOpts{}) // 只读
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("增加后计数器值:", count)

}

// 执行合约打印结果
// 合约hash：0xddF5Ac09594B38a7D6933b6Ea5CFD947Ad2381DD
func RunContract() {
	// 1. 连接到 Sepolia 测试网

	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	// 2. 加载私钥
	privateKey, err := crypto.HexToECDSA("private key")
	if err != nil {
		log.Fatal(err)
	}
	// 3. 获取 chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 4. 构造授权对象（不手动设置 Nonce，自动管理）
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	// 5. 用已部署合约地址创建实例
	contractAddress := common.HexToAddress("0xddF5Ac09594B38a7D6933b6Ea5CFD947Ad2381DD") // 替换为你的合约地址
	instance, err := counter.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 6. 调用自增方法
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("increment 交易哈希:", tx.Hash().Hex())

	// 7. 等待交易上链
	bind.WaitMined(context.Background(), client, tx)

	// 8. 查询最新计数器值
	count, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("自增后的计数器值:", count)
}
