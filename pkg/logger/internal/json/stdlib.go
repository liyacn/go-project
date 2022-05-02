//go:build !go_json && !sonic

package json

import "encoding/json"

var (
	NewEncoder = json.NewEncoder
	Valid      = json.Valid
)

type RawMessage = json.RawMessage
