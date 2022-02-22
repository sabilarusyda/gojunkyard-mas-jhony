package cache

import (
	"net/http"
	"net/http/httptest"
	"time"
)

type Stdlib struct {
	*cache
}

func NewStdlib(opts ...CacheOption) *Stdlib {
	return &Stdlib{newCache(opts...)}
}

func (router *Stdlib) Handle(next http.Handler) http.Handler {
	return router.HandleWithTTL(next, router.ttl)
}

func (router *Stdlib) HandleFunc(next http.HandlerFunc) http.HandlerFunc {
	return router.Handle(next).(http.HandlerFunc)
}

func (router *Stdlib) HandleWithTTL(next http.Handler, ttl time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// step 1. if method is not get, then skip the cache
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// step 2. normalize the url
		normalize(r.URL)

		var key = generateKey(r)
		if len(r.Header.Get("x-cache-refresh")) > 0 {
			// step 3a. delete key
			err := router.storage.Delete(key)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
		} else {
			// step 3b. get from cache
			obj, err := router.storage.Get(key)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if obj != nil {
				renderResponse(w, obj)
				return
			}
		}

		// step 4. if no cache then call the handler and retrieve the header, status code, and body
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		obj := &object{
			Body:       recorder.Body.Bytes(),
			Header:     recorder.HeaderMap,
			StatusCode: recorder.Code,
		}

		// step 5. if response 4xx and 5xx, then do not cache
		if recorder.Code > 400 {
			renderResponse(w, obj)
			return
		}

		// step 6. if response 2xx and 3xx, then cache the response
		router.storage.Set(key, obj, ttl)
		renderResponse(w, obj)
	})
}
