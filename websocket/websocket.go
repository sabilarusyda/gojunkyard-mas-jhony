package websocket

import (
	"time"

	"devcode.xeemore.com/systech/gojunkyard/websocket/connection"
)

// Config ...
type Config struct {
	HandshakeTimeout time.Duration `envconfig:"HANDSHAKE_TIMEOUT"` // HandshakeTimeout specifies the duration for the handshake to complete.
	ReadBufferSize   int           `envconfig:"READ_BUFFER_SIZE"`  // ReadBufferSize specify I/O buffer sizes.
	WriteBufferSize  int           `envconfig:"WRITE_BUFFER_SIZE"` // WriteBufferSize specify I/O buffer sizes
	Origins          []string      `envconfig:"ORIGINS"`           // Origins is parameter that used for prevent cross-site request forgery.
	OriginValid      bool          // OriginValid is parameter for CheckOrigin function

}

// New ...
func New(cfg *Config) *connection.Option {
	return &connection.Option{
		HandshakeTimeout: cfg.HandshakeTimeout,
		ReadBufferSize:   cfg.ReadBufferSize,
		WriteBufferSize:  cfg.WriteBufferSize,
		Origins:          cfg.Origins,
		OriginValid:      cfg.OriginValid,
	}
}
