package util

import (
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
)

// GenerateUUID returns new unique id without dash
func GenerateUUID() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}
