package cache

import (
	"encoding/base32"
	"project/pkg/random"
)

func GenerateToken() string {
	return base32.StdEncoding.EncodeToString(random.Bytes(20))
}
