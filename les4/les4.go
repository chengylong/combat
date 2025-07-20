package les4

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/chengylong/combat/store"
	"github.com/chengylong/combat/token"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// 创建钱包
func CW() {
	//生成私钥 即钱包
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	//转换成字节并转换为16进制字符串，截去首2位的0x 得到用于签署交易的私钥，将被视为密码，永远不应该被共享给别人，因为谁拥有它可以访问你的所有资产。
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("00000:", hexutil.Encode(privateKeyBytes)[2:])

	//通过私钥获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("publicKeyBytes", hexutil.Encode(publicKeyBytes)[4:])
	// 再根据公钥获取公共地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address", address)
	fmt.Println(address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("full:", hexutil.Encode(hash.Sum(nil)[:]))
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去12位，保留后20位

}

// eth转账
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

// 代币转账
func DBTra() {
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	// 获取私钥
	privateKey, err := crypto.HexToECDSA("71ddf4944784d1b5d946484bc7d3c5e39a9c86f3a4b0eba706e0e2e1f8223329")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 私钥->公钥->发送地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("fromAddress:", fromAddress.Hex())

	// 代币传输缺乏传输ETH，因此将交易“值”设置为“0”。
	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 首先您需要发送代币的地址存储在变量中。
	toAddress := common.HexToAddress("0x02061c94109EFF3A02BA09D24f7932692331ab96")
	// 让我们将代币合约地址赋予标志。
	tokenAddress := common.HexToAddress("0x4592D8f8D7B001e72Cb26A73e4Fa1806a51aC79d")
	//函数名将是传递函数的名称，即ERC-20规范中的transfer和参数类型。
	// 第一个参数类型是address（令牌的接收者），第二个类型是uint256（要发送的代币数量）。不需要没有空格和参数名称。我们还需要用字节切片格式。
	transferFnSignature := []byte("transfer(address,uint256)")
	// 我们现在队列go-ethereum导入crypto/sha3包以生成函数签名的Keccak256存储。然后我们只使用前4个字节来获取方法ID。
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb
	// 接下来，我们需要将给我们发送代币的地址左填充到 32 字节。
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d
	// 接下来我们确定要发送多少代币，在这个例子中是 1,000 个，并且我们需要在big.Int中格式化为 wei。
	amount := new(big.Int)
	amount.SetString("100000000000000000", 10) // 1000 tokens
	// 代币量还需要填充到32个字节。
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
	// 接下来我们只需将方法 ID，填充后的地址和填充后的字节量，即可成为我们数据字段的字节片。
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	// 燃气上限制将取决于交易数据的大小和智能合约执行必须的计算步骤。幸运的是，客户端提供了EstimateGas方法，
	// 它可以为我们提示所需的燃气量。该函数从ethereum包中获取CallMsg结构，我们在其中指定数据和地址。将返回我们提示的完成交易所需的估计燃气上限。
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit) // 23256
	// 接下来我们需要做的是构建交易类型，这类似于在 ETH 交互部分中的，除了_to_字段将是代币智能合约地址。这常让人困惑。我们还必须在调用中包含 0 ETH 的值字段并看到刚刚生成的数据字节。
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
	// 下一步是使用发件人的私钥对事务进行签名。SignTx方法需要EIP155igner，需要我们先从客户端获取链ID。

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 最后发送交易。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
}

// 查询账户余额
func QueryBalance() {
	//连接客户端
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	account := common.HexToAddress("0x9beEfA25246858238D8FC16D7D81875768a04AaB")
	//最新区块
	balance, err := client.BalanceAt(context.Background(), account, nil)
	fmt.Println(balance)

	// 指定区块
	// blockNumber := big.NewInt(8791585)
	// balance, err := client.BalanceAt(context.Background(), account, blockNumber)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue) //0.29899992256045196

	fmt.Println(balance) //298999922560451960

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println("pendingBalance:", pendingBalance)

}

// 查询代币余额
func QueryERCBalance() {
	//连接客户端
	client, err := ethclient.Dial("https://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	// 合约地址
	tokenAddress := common.HexToAddress("0x9b0971Aff0Bd1371A2b38f341FBbc010F6401283")
	// token实例来源于 abi，abi文件就是在remix上当我部署和编译完成我的合约后生成的MyToken_metadata.json 中的"abi"属性
	//在项目目录下创建token.abi,把刚才复制的 abi代码复制进去得到 abi文件
	// 得到对应得abi文件后使用 abigen --abi=token.abi --pkg=token --out=token/erc20.go 得到最终在go代码中可以调用的合约方法。
	instance, err := token.NewToken(tokenAddress, client)

	if err != nil {
		log.Fatal(err)
	}
	account := common.HexToAddress("0x9beEfA25246858238D8FC16D7D81875768a04AaB")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
}

// 订阅区块链中新增的区块
func SubscribeBlock() {
	//连接客户端
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	// 接下来，我们将创建一个新的通道，用于接收最新的区块头。
	headers := make(chan *types.Header)
	// 现在我们调用客户端的SubscribeNewHead方法，它接收我们刚刚创建的区块头通道，该方法将返回一个订阅对象。
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	// 订阅将订阅新的块头事件到我们的通道，因此可以使用一个选择语句来监听新消息。订阅对象还包括一个错误通道，该通道将在订阅失败时发送消息。
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(block.Hash().Hex())      // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64()) // 3477413
			fmt.Println(block.Time())            // 1529525947
			fmt.Println(block.Nonce())           // 130524141876765836fmt.Println(len(block.Transactions())) // 7
		}

	}

}

// 部署合约
func DeployContract() {
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("71ddf4944784d1b5d946484bc7d3c5e39a9c86f3a4b0eba706e0e2e1f8223329")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	// 执行合约前查询余额
	monitorBalance(client, fromAddress)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	// address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}
	// instance.Items()
	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())
	_ = instance
	key := stringToBytes32("hello")
	value := stringToBytes32("world")

	instance.SetItem(auth, key, value)
	//执行合约后查询余额
	monitorBalance(client, fromAddress)

	// instance.WatchItemSet()
}
func stringToBytes32(str string) [32]byte {
	var result [32]byte
	copy(result[:], []byte(str))
	return result
}

// 加载合约
func LoadContract() {
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	storeContract, err := store.NewStore(common.HexToAddress("0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8"), client)
	if err != nil {
		log.Fatal(err)
	}

	_ = storeContract
	version, err := storeContract.Version(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询版本失败:", err)
	}
	fmt.Printf("合约版本: %s\n", version)
	// 	key := stringToBytes32("hello")
	// value := stringToBytes32("world")
	// storeContract.SetItem(key)
}

// 执行合约 71ddf4944784d1b5d946484bc7d3c5e39a9c86f3a4b0eba706e0e2e1f8223329
func RunContract() {
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	// 创建合约实例
	storeContract, err := store.NewStore(common.HexToAddress("0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8"), client)
	if err != nil {
		log.Fatal(err)
	}
	// 获取私钥准备执行合约
	privateKey, err := crypto.HexToECDSA("71ddf4944784d1b5d946484bc7d3c5e39a9c86f3a4b0eba706e0e2e1f8223329")
	if err != nil {
		log.Fatal(err)
	}
	// 获取账户地址
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 查看执行前的余额
	balanceBefore, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("获取余额失败:", err)
	}

	ethBalanceBefore := new(big.Float).Quo(new(big.Float).SetInt(balanceBefore), big.NewFloat(1e18))
	fmt.Printf("执行前余额: %f ETH\n", ethBalanceBefore)

	// 调用合约方法
	// 准备数据
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("demo_save_key"))
	copy(value[:], []byte("demo_save_value11111"))

	// 初始化交易opt实例
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}
	// 调用合约方法
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex())
	// 等待交易确认
	fmt.Println("等待交易确认...")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal("等待交易确认失败:", err)
	}
	if receipt.Status == 1 {
		fmt.Println("✅ 交易成功")

		// 计算燃气费用
		gasUsed := receipt.GasUsed
		effectiveGasPrice := receipt.EffectiveGasPrice
		totalCost := new(big.Int).Mul(big.NewInt(int64(gasUsed)), effectiveGasPrice)

		fmt.Printf("燃气消耗: %d 单位\n", gasUsed)
		fmt.Printf("燃气价格: %s wei\n", effectiveGasPrice.String())
		fmt.Printf("总费用: %s wei\n", totalCost.String())

		// 转换为 ETH
		ethCost := new(big.Float).Quo(new(big.Float).SetInt(totalCost), big.NewFloat(1e18))
		fmt.Printf("总费用: %f ETH\n", ethCost)

		// 查看执行后的余额
		balanceAfter, err := client.BalanceAt(context.Background(), fromAddress, nil)
		if err != nil {
			log.Fatal("获取余额失败:", err)
		}

		ethBalanceAfter := new(big.Float).Quo(new(big.Float).SetInt(balanceAfter), big.NewFloat(1e18))
		fmt.Printf("执行后余额: %f ETH\n", ethBalanceAfter)

		// 计算实际消耗
		balanceDiff := new(big.Int).Sub(balanceBefore, balanceAfter)
		ethBalanceDiff := new(big.Float).Quo(new(big.Float).SetInt(balanceDiff), big.NewFloat(1e18))
		fmt.Printf("实际消耗: %f ETH\n", ethBalanceDiff)

	} else {
		fmt.Println("❌ 交易失败")
	}

}

// 执行合约  仅使用 ethclient 包调用合约
func RunContractByEth() {
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	// 获取私钥准备执行合约
	privateKey, err := crypto.HexToECDSA("71ddf4944784d1b5d946484bc7d3c5e39a9c86f3a4b0eba706e0e2e1f8223329")
	if err != nil {
		log.Fatal(err)
	}
	// 获取账户地址
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 准备交易所需数据
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 合约abi
	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`))
	if err != nil {
		log.Fatal(err)
	}

	methodName := "setItem"
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("demo_save_key_use_abi"))
	copy(value[:], []byte("demo_save_value_use_abi_11111"))
	input, err := contractABI.Pack(methodName, key, value)
	// 创建和交易签名
	contractAddr := "0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8"
	chainID := big.NewInt(int64(11155111))
	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddr), big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 签名后发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	// call查询
	callInput, err := contractABI.Pack("items", key)
	if err != nil {
		log.Fatal(err)
	}
	to := common.HexToAddress(contractAddr)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	var unpacked [32]byte
	contractABI.UnpackIntoInterface(&unpacked, "items", result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("is value saving in contract equals to origin value:", unpacked == value)

}
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}

// 查询合约中的数据
func QueryContract() {
	// 已知key
	// client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // 创建合约实例
	// contractAddress := common.HexToAddress("0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8")
	// instance, err := store.NewStore(contractAddress, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("=== 查询合约中存储的数据 ===")

	// // 查询你存储的数据
	// keys := []string{
	// 	"demo_save_key",
	// 	"hello",
	// 	"user_name",
	// 	"balance",
	// 	"status",
	// 	"key1",
	// 	"key2",
	// 	"key3",
	// }

	// for _, keyStr := range keys {
	// 	key := stringToBytes32(keyStr)

	// 	// 调用合约的 items 映射查询
	// 	value, err := instance.Items(&bind.CallOpts{}, key)
	// 	if err != nil {
	// 		fmt.Printf("查询 %s 失败: %v\n", keyStr, err)
	// 		continue
	// 	}

	// 	// 检查是否有数据
	// 	if isEmptyBytes32(value) {
	// 		fmt.Printf("%s: 无数据\n", keyStr)
	// 	} else {
	// 		valueStr := bytes32ToString(value)
	// 		fmt.Printf("%s = %s\n", keyStr, valueStr)
	// 	}
	// }
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8")

	fmt.Println("=== 使用 ethclient 查询数据 ===")

	keys := []string{"demo_save_key", "hello", "user_name"}

	for _, keyStr := range keys {
		key := stringToBytes32(keyStr)

		// 构造查询数据
		data := constructItemsQueryData(key)

		// 发送查询
		msg := ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		}

		result, err := client.CallContract(context.Background(), msg, nil)
		if err != nil {
			fmt.Printf("查询 %s 失败: %v\n", keyStr, err)
			continue
		}

		// 解析结果
		if len(result) >= 32 {
			var value [32]byte
			copy(value[:], result[:32])

			if isEmptyBytes32(value) {
				fmt.Printf("%s: 无数据\n", keyStr)
			} else {
				valueStr := bytes32ToString(value)
				fmt.Printf("%s = %s\n", keyStr, valueStr)
			}
		}
	}

}

// 构造 items 映射查询数据
func constructItemsQueryData(key [32]byte) []byte {
	// 函数选择器：items(bytes32)
	functionSignature := []byte("items(bytes32)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(functionSignature)
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)
	data = append(data, key[:]...)

	return data
}

// 检查 bytes32 是否为空
func isEmptyBytes32(bytes [32]byte) bool {
	for _, b := range bytes {
		if b != 0 {
			return false
		}
	}
	return true
}

// bytes32 转字符串
func bytes32ToString(bytes [32]byte) string {
	end := 0
	for i, b := range bytes {
		if b == 0 {
			end = i
			break
		}
		end = i + 1
	}
	return string(bytes[:end])
}

// 监控账户余额
func monitorBalance(client *ethclient.Client, address common.Address) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal("获取余额失败:", err)
	}

	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("当前余额: %f ETH\n", ethBalance)

	// 更安全的转换
	balanceFloat, accuracy := ethBalance.Float64()
	if accuracy != big.Exact {
		fmt.Println("警告：余额转换可能不精确")
	}

	costPerOp := 0.0001
	remainingOps := int(balanceFloat / costPerOp)
	fmt.Printf("预估剩余操作次数: %d\n", remainingOps)
}

// 合约事件
func ContractEvent() {
	var StoreABI = `[{"inputs":
	[{"internalType":"string","name":"_version","type":"string"}]
	,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":
	[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},
	{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],
	"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],
	"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`
	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress("0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(8797325),
		ToBlock:   big.NewInt(8797504),
		Addresses: []common.Address{
			contractAddress,
		},
		// Topics: [][]common.Hash{
		//  {},
		//  {},
		// },
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}
	for _, vLog := range logs {
		fmt.Println(vLog.BlockHash.Hex())
		fmt.Println(vLog.BlockNumber)
		fmt.Println(vLog.TxHash.Hex())
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(common.Bytes2Hex(event.Key[:]))
		fmt.Println(common.Bytes2Hex(event.Value[:]))
		var topics []string
		for i := range vLog.Topics {
			topics = append(topics, vLog.Topics[i].Hex())
		}
		fmt.Println("topics[0]=", topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
		if len(topics) > 1 {
			fmt.Println("indexed topics:", topics[1:])
		}
	}
}

// 合约事件+
func ContractEventPlus() {
	var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

	client, err := ethclient.Dial("wss://ethereum-sepolia.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x9C024f2A5379b605D26Dd548dbA646B1bef5D9A8")

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(8797325),
		ToBlock:   big.NewInt(8797504),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		fmt.Println("=== 事件信息 ===")
		fmt.Println("区块哈希:", vLog.BlockHash.Hex())
		fmt.Println("区块号:", vLog.BlockNumber)
		fmt.Println("交易哈希:", vLog.TxHash.Hex())

		// 解析事件数据
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		// 获取实际存储的 key（从事件参数中）
		storageKey := string(bytes.TrimRight(event.Value[:], "\x00"))
		fmt.Println("存储的 Key:", storageKey)

		// 将 key 转换为 [32]byte 格式
		keyArray := [32]byte{}
		copy(keyArray[:], []byte(storageKey))

		// 构建调用数据
		data, err := contractAbi.Pack("items", keyArray)
		if err != nil {
			fmt.Printf("构建调用数据失败: %v\n", err)
			continue
		}

		// 调用合约
		msg := ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		}

		result, err := client.CallContract(context.Background(), msg, nil)
		if err != nil {
			fmt.Printf("调用合约失败: %v\n", err)
			continue
		}

		// 解析返回结果
		var actualValue [32]byte
		err = contractAbi.UnpackIntoInterface(&actualValue, "items", result)
		if err != nil {
			fmt.Printf("解析返回结果失败: %v\n", err)
			continue
		}

		fmt.Println("真正的存储 Value (hex):", common.Bytes2Hex(actualValue[:]))
		fmt.Println("真正的存储 Value (string):", string(bytes.TrimRight(actualValue[:], "\x00")))

		// 如果 value 是数字
		valueNumber := new(big.Int).SetBytes(actualValue[:])
		fmt.Println("真正的存储 Value (number):", valueNumber.String())

		fmt.Println("---")
	}
}
