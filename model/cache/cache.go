package cache

import (
	"crypto/sha1"
	"encoding/base32"
	"project/pkg/random"
)

func GenerateToken(b []byte) string {
	h := sha1.New()
	h.Write(b)
	h.Write(random.Bytes(16))
	return base32.StdEncoding.EncodeToString(h.Sum(nil))
}
