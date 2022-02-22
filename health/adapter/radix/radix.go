package radix

import (
	"github.com/mediocregopher/radix/v3"
)

// Radix ...
type Radix struct {
	name string
	pool *radix.Pool
}

// New ...
func New(pool *radix.Pool) *Radix {
	return &Radix{pool: pool}
}

// SetName ...
func (r *Radix) SetName(name string) {
	r.name = name
}

// Name ...
func (r *Radix) Name() string {
	if len(r.name) == 0 {
		return "REDIS"
	}
	return r.name
}

// Check ...
func (r *Radix) Check() error {
	return r.pool.Do(radix.Cmd(nil, "PING"))
}
