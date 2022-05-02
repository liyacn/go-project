package files

import (
	"crypto/sha1"
	"encoding/base32"
)

var enc = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)

func GenFileName(b []byte) string {
	sum := sha1.Sum(b)
	return enc.EncodeToString(sum[:])
}
