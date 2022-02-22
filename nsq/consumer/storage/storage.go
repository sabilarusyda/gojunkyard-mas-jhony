package storage

import "time"

type Storage interface {
	SetNX(key string, ttl time.Duration) (bool, error)
	Delete(key string) error
}
