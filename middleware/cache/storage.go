package cache

import "time"

type Storage interface {
	Get(key uint64) (obj *object, err error)
	Set(key uint64, obj *object, expire time.Duration) (err error)
	Delete(key uint64) (err error)
}

type StorageOption func(Storage)
