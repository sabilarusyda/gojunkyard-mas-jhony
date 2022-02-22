package radix

import (
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v3"
)

type Radix struct {
	client radix.Client
}

func New(client radix.Client) *Radix {
	return &Radix{
		client: client,
	}
}

func (r *Radix) SetNX(key string, ttl time.Duration) (bool, error) {
	var b bool
	err := r.client.Do(radix.Pipeline(
		radix.Cmd(&b, "SETNX", key, ""),
		radix.Cmd(nil, "EXPIRE", key, strconv.FormatFloat(ttl.Seconds(), 'f', 0, 64)),
	))
	return b, err
}

func (r *Radix) Delete(key string) error {
	return r.client.Do(radix.Cmd(nil, "DEL", key))
}
