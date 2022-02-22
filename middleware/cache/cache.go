package cache

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/cespare/xxhash"
)

type CacheOption func(*cache)

func SetCacheStorage(storage Storage) CacheOption {
	return func(c *cache) {
		c.SetStorage(storage)
	}
}

func SetCacheTTL(ttl time.Duration) CacheOption {
	return func(c *cache) {
		c.SetTTL(ttl)
	}
}

type cache struct {
	storage Storage
	ttl     time.Duration
}

func newCache(opts ...CacheOption) *cache {
	c := new(cache)
	// step 1. set option to cache object
	for _, opt := range opts {
		opt(c)
	}
	// step 2. set storage to local not set yet
	if c.storage == nil {
		c.SetStorage(NewInMemory())
	}
	// step 3. return cache object
	return c
}

func (c *cache) SetStorage(storage Storage) {
	c.storage = storage
}

func (c *cache) SetTTL(ttl time.Duration) {
	c.ttl = ttl
}

func renderResponse(w http.ResponseWriter, obj *object) {
	header := w.Header()
	for k, v := range obj.Header {
		for _, v := range v {
			header.Add(k, v)
		}
	}
	w.WriteHeader(obj.StatusCode)
	w.Write(obj.Body)
}

func normalize(url *url.URL) {
	qs := url.Query()
	for _, q := range qs {
		sort.Slice(q, func(i, j int) bool {
			return q[i] < q[j]
		})
	}
	url.RawQuery = qs.Encode()
}

func generateKey(r *http.Request) uint64 {
	return xxhash.Sum64String(r.URL.String())
}

type object struct {
	Header     http.Header     `json:"header"`
	Body       json.RawMessage `json:"body"`
	StatusCode int             `json:"code"`
}

var now = time.Now
