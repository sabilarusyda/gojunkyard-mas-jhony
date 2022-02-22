package writer

import (
	"strings"
)

type LEVEL uint8

const (
	DEBUG LEVEL = iota + 1
	INFO
	WARNING
	ERROR
	FATAL
)

// Decode ...
func (l *LEVEL) Decode(value string) error {
	switch strings.ToUpper(value) {
	case "DEBUG":
		*l = DEBUG
	case "WARNING":
		*l = WARNING
	case "ERROR":
		*l = ERROR
	case "FATAL":
		*l = FATAL
	default:
		*l = INFO
	}
	return nil
}
