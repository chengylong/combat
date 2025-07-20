package main

import (
	"encoding/hex"
	"strings"

	"github.com/chengylong/combat/dapp1"
)

func main() {

	dapp1.RunContract()
	// 使用示例
	// keyHex := "64656d6f5f736176655f76616c75655f7573655f6162695f3131313131000000"
	// keyStr := decodeHexString(keyHex)
	// fmt.Println("键:", keyStr) // 输出: demo_save_key_use_abi
}
func decodeHexString(hexStr string) string {
	// 移除末尾的零填充
	hexStr = strings.TrimRight(hexStr, "0")

	// 每两个字符转换为一个字节
	var bytes []byte
	for i := 0; i < len(hexStr); i += 2 {
		if i+1 < len(hexStr) {
			b, _ := hex.DecodeString(hexStr[i : i+2])
			bytes = append(bytes, b...)
		}
	}

	return string(bytes)
}
