package redigo

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	client *redis.Pool
}

func New(client *redis.Pool) *Redis {
	return &Redis{
		client: client,
	}
}

func (r *Redis) SetNX(key string, ttl time.Duration) (bool, error) {
	conn := r.client.Get()
	defer conn.Close()

	conn.Send("SETNX", key, "")
	conn.Send("EXPIRE", key, 180)
	conn.Flush()

	return redis.Bool(conn.Receive())
}

func (r *Redis) Delete(key string) error {
	conn := r.client.Get()
	_, err := conn.Do("DEL", key)
	conn.Close()
	return err
}
