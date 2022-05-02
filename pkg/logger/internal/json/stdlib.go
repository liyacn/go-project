//go:build !go_json && !sonic

package json

import "encoding/json"

var NewEncoder = json.NewEncoder
