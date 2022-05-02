//go:build sonic && (linux || windows || darwin)

package json

import "github.com/bytedance/sonic"

var (
	json = sonic.ConfigStd

	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
	Valid         = json.Valid
)

type RawMessage = sonic.NoCopyRawMessage
