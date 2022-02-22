package tcp

import (
	"net"
	"time"
)

// TCP ...
type TCP struct {
	addr string
	name string
}

// New ...
func New(addr string) *TCP {
	return &TCP{addr: addr}
}

// SetName ...
func (t *TCP) SetName(name string) {
	t.name = name
}

// Name ...
func (t *TCP) Name() string {
	if len(t.name) == 0 {
		return "TCP"
	}
	return t.name
}

// Check ...
func (t *TCP) Check() error {
	conn, err := net.DialTimeout("tcp", t.addr, 2*time.Second)
	if err != nil {
		return err
	}
	return conn.Close()
}
