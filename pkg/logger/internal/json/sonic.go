//go:build sonic

package json

import "github.com/bytedance/sonic"

var (
	NewEncoder = sonic.ConfigDefault.NewEncoder
	Valid      = sonic.Valid
)

type RawMessage = sonic.NoCopyRawMessage
