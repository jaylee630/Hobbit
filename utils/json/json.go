// +build !jsoniter

package json

import "encoding/json"

var (
	// Marshal is exported by utils/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by utils/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by utils/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by utils/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by utils/json package.
	NewEncoder = json.NewEncoder
)
