package random

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strconv"
)

// Bytes 随机字节数组
func Bytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// UUID V4版本不带-号的UUID
func UUID() string { return hex.EncodeToString(Bytes(16)) }

// NumCode 安全的六位数字码
func NumCode() string {
	num, _ := rand.Int(rand.Reader, big.NewInt(1e6))
	return strconv.FormatInt(num.Int64()+1e6, 10)[1:]
} //优于fmt.Sprintf("%06d", num.Int64())
