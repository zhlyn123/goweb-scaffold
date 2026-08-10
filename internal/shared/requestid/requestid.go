package requestid

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	HeaderName  = "X-Request-ID"
	ConstextKey = "request_id"
)

func New() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(bytes)
}
