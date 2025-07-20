# combat
合约的创建和部署过程：
1. 编写智能合约（Solidity）
示例：计数器合约 Counter.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 public count;

    event CountedTo(uint256 newCount);

    function increment() public {
        count += 1;
        emit CountedTo(count);
    }

    function getCount() public view returns (uint256) {
        return count;
    }
}
2. 编译智能合约，生成 ABI 和字节码
Remix 编译
打开 remix.ethereum.org
新建 Counter.sol 文件，粘贴合约代码
选择 Solidity 编译器，点击“编译”
在“Artifacts”中下载 ABI 和字节码
3. 使用 abigen 生成 Go 绑定代码，才可以在go项目中使用合约中的方法
abigen 是 go-ethereum 提供的工具，可以根据 ABI 和 bin 文件生成 Go 代码。
安装 abigen（如果未安装）：
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
安装后配置好环境变量才可以使用
写好sol文件
生成bin文件和abi文件
solcjs --bin --abi Counter.sol
生成 Go 绑定代码
abigen --bin=Counter_sol_Counter.bin --abi=Counter_sol_Counter.abi --pkg=counter --out=counter.go
4. 结果
你会得到一个 counter.go 文件，里面包含了 Counter 合约的 Go 绑定代码。
之后你就可以在 Go 代码中像操作普通对象一样与合约交互了。
5.部署合约
在代码中实现
