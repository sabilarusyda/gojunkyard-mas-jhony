package redigo

import (
	"github.com/gomodule/redigo/redis"
)

// Redigo ...
type Redigo struct {
	name string
	pool *redis.Pool
}

// New ...
func New(pool *redis.Pool) *Redigo {
	return &Redigo{pool: pool}
}

// SetName ...
func (r *Redigo) SetName(name string) {
	r.name = name
}

// Name ...
func (r *Redigo) Name() string {
	if len(r.name) == 0 {
		return "REDIS"
	}
	return r.name
}

// Check ...
func (r *Redigo) Check() error {
	conn := r.pool.Get()
	_, err := conn.Do("PING")
	conn.Close()
	return err
}
