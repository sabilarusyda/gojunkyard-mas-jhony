package redis

import (
	"time"
)

// IRedis Interface of redis
type IRedis interface {
	Connected() bool
	Exist(result interface{}, key string) error
	Expire(key string, expiration time.Duration) error
	Del(key string) error
	Get(result interface{}, key string) (bool, error)
	Set(key string, value interface{}) error
	SetEX(key string, value interface{}, expiration time.Duration) error
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	HDel(key string, fields ...string) error
	HGet(result interface{}, key, field string) error
	HGetAll(result interface{}, key string) error
	HMGet(result interface{}, key string, fields ...string) error
	HMSet(key string, fields map[string]interface{}) error
	HSet(key, field string, value interface{}) error
	HSetEX(key, field string, value interface{}, expiration time.Duration) error
	SAdd(key string, value interface{}) error
	SRem(key string, value interface{}) error
	ZAdd(key, score, value string) error
	ZRange(key, start, end string) (result []string, err error)
	ZRevRange(key, start, end string) (result []string, err error)
	ZRem(key, value string) error
	ZRemRangeByscore(key, min, max string) error
	TTL(result interface{}, key string) error
	Incr(key string) error
	Pipeline([]Cmd) error
	MGet(result interface{}, keys []string) error
	Keys(key string) (result []string, err error)
}

// Cmd is command data used for pipeline.
type Cmd struct {
	// Data retrieved from redis.
	// Should be a pointer to a variable.
	// Like when using json.Unmarshal().
	Return interface{}
	// Redis command. SET, GET, DEL, etc.
	Command string
	// Redis key.
	Key string
	// Additional arguments.
	Args []interface{}
}
