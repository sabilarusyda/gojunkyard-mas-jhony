package cache

import (
	"encoding/json"
	"errors"
	"time"

	"devcode.xeemore.com/systech/gojunkyard/conn"

	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
)

type RedisOption func(*Redis) error

type Redis struct {
	redis *redis.Pool
}

func NewRedis(opts ...RedisOption) (*Redis, error) {
	c := new(Redis)
	// step 1. set option to cache object
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	// step 2. set storage to local not set yet
	if c.redis == nil {
		return nil, errors.New("redis connection is not initialized")
	}
	// step 3. return cache object
	return c, nil
}

func (r *Redis) SetPool(pool *redis.Pool) {
	r.redis = pool
}

func SetRedisConfig(cfg conn.RedisConfig) RedisOption {
	return func(redis *Redis) error {
		pool, err := conn.InitRedis(cfg)
		if err != nil {
			return err
		}
		redis.SetPool(pool)
		return nil
	}
}

func SetRedisPool(pool *redis.Pool) RedisOption {
	return func(redis *Redis) error {
		redis.SetPool(pool)
		return nil
	}
}

func (r *Redis) Get(key uint64) (*object, error) {
	conn := r.redis.Get()
	byt, err := redis.Bytes(conn.Do("GET", key))
	conn.Close()
	if err == redis.ErrNil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var obj object
	err = jsoniter.ConfigFastest.Unmarshal(byt, &obj)
	if err != nil {
		return nil, nil
	}

	return &obj, nil
}

func (r *Redis) Set(key uint64, obj *object, expire time.Duration) error {
	if expire < 0 {
		return nil
	}

	// step 1. prepare data
	byt, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// step 2. set data to redis
	conn := r.redis.Get()
	if expire == 0 {
		_, err = conn.Do("SET", key, byt)
	} else {
		_, err = conn.Do("SET", key, byt, "EX", int64(expire.Seconds()))
	}
	conn.Close()
	return err
}

func (r *Redis) Delete(key uint64) error {
	conn := r.redis.Get()
	_, err := conn.Do("DEL", key)
	conn.Close()
	return err
}
