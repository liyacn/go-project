//go:build !jsoniter && !go_json && !(sonic && (linux || windows || darwin))

package json

import "encoding/json"

var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
	Valid         = json.Valid
)

type RawMessage = json.RawMessage
