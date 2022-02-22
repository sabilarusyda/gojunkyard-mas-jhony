package cache

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/julienschmidt/httprouter"
)

type HTTPRouter struct {
	*cache
}

func NewHTTPRouter(opts ...CacheOption) *HTTPRouter {
	return &HTTPRouter{newCache(opts...)}
}

func (router *HTTPRouter) Handle(h httprouter.Handle) httprouter.Handle {
	return router.HandleWithTTL(h, router.ttl)
}

func (router *HTTPRouter) HandleWithTTL(h httprouter.Handle, ttl time.Duration) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// step 1. if method is not get, then skip the cache
		if r.Method != http.MethodGet {
			h(w, r, ps)
			return
		}

		// step 2. normalize the url
		normalize(r.URL)

		var key = generateKey(r)
		if len(r.Header.Get("x-cache-refresh")) > 0 {
			// step 3a. delete key
			err := router.storage.Delete(key)
			if err != nil {
				h(w, r, ps)
				return
			}
		} else {
			// step 3b. get from cache
			obj, err := router.storage.Get(key)
			if err != nil {
				h(w, r, ps)
				return
			}
			if obj != nil {
				renderResponse(w, obj)
				return
			}
		}

		// step 4. if no cache then call the handler and retrieve the header, status code, and body
		recorder := httptest.NewRecorder()
		h(recorder, r, ps)

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
	}
}
