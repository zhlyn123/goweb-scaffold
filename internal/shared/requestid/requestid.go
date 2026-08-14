package requestid

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	HeaderName  = "X-Request-ID"
	ContextKey  = "request_id"
	ConstextKey = ContextKey
)

func New() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(bytes)
}
