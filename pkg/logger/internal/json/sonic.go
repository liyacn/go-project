//go:build sonic

package json

import "github.com/bytedance/sonic"

var NewEncoder = sonic.ConfigDefault.NewEncoder
