package cache

import (
	"sync"
	"time"
)

type inMemoryObject struct {
	Obj        *object    `json:"obj"`
	ExpireTime *time.Time `json:"expireTime"`
}

type InMemoryOption func(*InMemory)
type InMemory struct {
	storage map[uint64]*inMemoryObject
	mux     sync.RWMutex
}

func NewInMemory(opts ...InMemoryOption) *InMemory {
	const defaultCapacity = 1000
	im := new(InMemory)
	// step 1. set option to cache object
	for _, opt := range opts {
		opt(im)
	}
	// step 2. set storage to local not set yet
	if im.storage == nil {
		im.SetCapacity(defaultCapacity)
	}
	// step 3. return cache object
	return im
}

func (im *InMemory) SetCapacity(cap int) {
	im.mux.Lock()
	if len(im.storage) < cap {
		storage := make(map[uint64]*inMemoryObject, cap)
		for k, v := range im.storage {
			storage[k] = v
		}
		im.storage = storage
	}
	im.mux.Unlock()
}

func SetInMemoryCapacity(cap int) InMemoryOption {
	return func(im *InMemory) {
		im.SetCapacity(cap)
	}
}

func (im *InMemory) Get(key uint64) (obj *object, err error) {
	im.mux.RLock()
	v, ok := im.storage[key]
	im.mux.RUnlock()
	if !ok || v == nil || (v.ExpireTime != nil && v.ExpireTime.Before(now())) {
		return nil, nil
	}
	return v.Obj, nil
}

func (im *InMemory) Set(key uint64, obj *object, expire time.Duration) error {
	if expire < 0 {
		return nil
	}
	var t *time.Time
	if expire > 0 {
		et := now().Add(expire)
		t = &et
	}
	im.mux.Lock()
	im.storage[key] = &inMemoryObject{obj, t}
	im.mux.Unlock()
	return nil
}

func (im *InMemory) Delete(key uint64) error {
	im.mux.Lock()
	delete(im.storage, key)
	im.mux.Unlock()
	return nil
}
