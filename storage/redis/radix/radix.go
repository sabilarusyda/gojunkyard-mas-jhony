package radix

import (
	"encoding/json"
	"time"

	"devcode.xeemore.com/systech/gojunkyard/storage/redis"

	"github.com/eapache/go-resiliency/breaker"
	_radix "github.com/mediocregopher/radix/v3"
)

// To check if Radix implements IRedis.
var _ redis.IRedis = &Radix{}

// Radix radix struct
type Radix struct {
	Address       string
	MaxConnection int
	Pool          *_radix.Pool
	PoolOpts      []_radix.PoolOpt
	Breaker       *breaker.Breaker
}

// Exist ..
func (radix *Radix) Exist(result interface{}, key string) error {
	mn := _radix.MaybeNil{Rcv: result}
	err := radix.Pool.Do(_radix.FlatCmd(&mn, "EXISTS", key))
	return err
}

// Expire ...
func (radix *Radix) Expire(key string, expiration time.Duration) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "EXPIRE", key, int64(expiration.Seconds())))
	return err
}

// Del ...
func (radix *Radix) Del(key string) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "DEL", key))
	return err
}

// Get ...
func (radix *Radix) Get(result interface{}, key string) (bool, error) {
	mn := _radix.MaybeNil{Rcv: result}
	err := radix.Pool.Do(_radix.FlatCmd(&mn, "GET", key))
	if mn.Nil {
		return true, err
	}

	return false, err
}

// Set ...
func (radix *Radix) Set(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = radix.Pool.Do(_radix.FlatCmd(nil, "SET", key, b))
	return err
}

// SetEX ...
func (radix *Radix) SetEX(key string, value interface{}, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = radix.Pool.Do(_radix.FlatCmd(nil, "SET", key, b, "EX", int64(expiration.Seconds())))
	return err
}

// SetNX ...
func (radix *Radix) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	var cond bool
	b, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	err = radix.Pool.Do(_radix.FlatCmd(&cond, "SET", key, b, "EX", int64(expiration.Seconds()), "NX"))
	return cond, err
}

// HDel ...
func (radix *Radix) HDel(key string, fields ...string) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "HDEL", key, fields))
	return err
}

// HGet ...
func (radix *Radix) HGet(result interface{}, key, field string) error {
	mn := _radix.MaybeNil{Rcv: result}
	err := radix.Pool.Do(_radix.FlatCmd(&mn, "HGET", key, field))
	return err
}

// HGetAll ...
func (radix *Radix) HGetAll(result interface{}, key string) error {
	mn := _radix.MaybeNil{Rcv: result}
	err := radix.Pool.Do(_radix.FlatCmd(&mn, "HGETALL", key))
	return err
}

// HMGet ...
func (radix *Radix) HMGet(result interface{}, key string, fields ...string) error {
	mn := _radix.MaybeNil{Rcv: result}
	err := radix.Pool.Do(_radix.FlatCmd(&mn, "HMGET", key, fields))
	return err
}

// MGet ...
func (radix *Radix) MGet(result interface{}, keys []string) error {
	return radix.Pool.Do(_radix.Cmd(&result, "MGET", keys...))
}

// HMSet ...
func (radix *Radix) HMSet(key string, fields map[string]interface{}) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "HMSET", key, fields))
	return err
}

// HSet ...
func (radix *Radix) HSet(key, field string, value interface{}) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "HSET", key, field, value))
	return err
}

// HSetEX ...
func (radix *Radix) HSetEX(key, field string, value interface{}, expiration time.Duration) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "HSET", key, field, value, "EX", int64(expiration.Seconds())))
	return err
}

// SAdd ...
func (radix *Radix) SAdd(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = radix.Pool.Do(_radix.FlatCmd(nil, "SADD", key, b))
	return err
}

// SRem ...
func (radix *Radix) SRem(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = radix.Pool.Do(_radix.FlatCmd(nil, "SREM", key, b))
	return err
}

// ZAdd ...
func (radix *Radix) ZAdd(key, score, value string) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "ZADD", key, score, value))
	return err
}

// ZRange ...
func (radix *Radix) ZRange(key, start, end string) (result []string, err error) {
	err = radix.Pool.Do(_radix.FlatCmd(&result, "ZRANGE", key, start, end))
	return result, err
}

// ZRevRange ...
func (radix *Radix) ZRevRange(key, start, end string) (result []string, err error) {
	err = radix.Pool.Do(_radix.FlatCmd(&result, "ZREVRANGE", key, start, end))
	return result, err
}

// ZRem ...
func (radix *Radix) ZRem(key, value string) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "ZREM", key, value))
	return err
}

// ZRemRangeByscore ...
func (radix *Radix) ZRemRangeByscore(key, min, max string) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "ZREMRANGEBYSCORE", key, min, max))
	return err
}

// TTL ...
func (radix *Radix) TTL(result interface{}, key string) error {
	mn := _radix.MaybeNil{Rcv: result}
	err := radix.Pool.Do(_radix.FlatCmd(&mn, "TTL", key))

	return err
}

// Incr ...
func (radix *Radix) Incr(key string) error {
	err := radix.Pool.Do(_radix.FlatCmd(nil, "INCR", key))
	return err
}

// Keys ...
func (radix *Radix) Keys(key string) (result []string, err error) {
	err = radix.Pool.Do(_radix.FlatCmd(&result, "KEYS", key))
	return result, err
}

// Pipeline ...
// Example:
//
//  tmp := "sample data"
//  expiration := 5 * time.Second
//
//  var returnedData string
//  err := Pipeline([]redis.Cmd{
//  	{
//  		Command: "SET",
//  		Key:     "key",
//  		Args:    []interface{}{tmp, "EX", int64(expiration.Seconds()), "NX"},
//  	},
//  	{
//  		Return:  &returnedData,
//  		Command: "GET",
//  		Key:     "key",
//  	},
//  })
//
func (radix *Radix) Pipeline(commands []redis.Cmd) error {
	var cmds []_radix.CmdAction
	for _, cmd := range commands {
		cmds = append(cmds, _radix.FlatCmd(cmd.Return, cmd.Command, cmd.Key, cmd.Args...))
	}
	return radix.Pool.Do(_radix.Pipeline(cmds...))
}
