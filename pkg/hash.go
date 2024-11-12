package pkg

import (
	"crypto/sha256"
	"math/big"
)

func SashToRange(input string, size int64) int64 {
	// 使用 SHA-256 计算哈希
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	// 将哈希值转换为大整数
	hashInt := new(big.Int).SetBytes(hashBytes)

	// 取模得到 0-7 的范围，加 1 变成 1-8
	return hashInt.Mod(hashInt, big.NewInt(size)).Int64() + 1
}
