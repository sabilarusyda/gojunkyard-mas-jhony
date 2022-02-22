package redigo

import (
	"time"
)

type Nop struct{}

func New() *Nop {
	return new(Nop)
}

func (n *Nop) SetNX(key string, ttl time.Duration) (bool, error) {
	return true, nil
}

func (n *Nop) Delete(key string) error {
	return nil
}
