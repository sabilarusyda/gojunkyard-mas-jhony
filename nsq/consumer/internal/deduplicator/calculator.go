package deduplicator

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

func calculateKey(topic, channel string, msg []byte) string {
	return fmt.Sprintf("NSQ_DEDUPLICATOR:%s:%s:%x", topic, channel, sha3.Sum256(msg))
}
