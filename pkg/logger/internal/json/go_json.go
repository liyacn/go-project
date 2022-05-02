//go:build go_json

package json

import json "github.com/goccy/go-json"

var (
	NewDecoder = json.NewDecoder
	Valid      = json.Valid
)

type RawMessage = json.RawMessage
