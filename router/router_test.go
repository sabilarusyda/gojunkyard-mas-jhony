package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, NewWithEngine(HTTPRouter), New())
}

func TestNewWithEngine(t *testing.T) {
	assert.Equal(t, &Router{
		middlewares: make([]middleware, 0),
		engine:      getRouterEngine(HTTPRouter),
	}, NewWithEngine(HTTPRouter))
}

func TestRouter_Use(t *testing.T) {
	var f = func(http.Handler) http.Handler {
		return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}
	router := New()
	assert.Len(t, router.middlewares, 0)
	router = router.Use(f)
	assert.Len(t, router.middlewares, 1)
	router = router.Use(f)
	assert.Len(t, router.middlewares, 2)
}

func TestRouter_Group(t *testing.T) {
	var f = func(http.Handler) http.Handler {
		return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}

	router := New()
	assert.Len(t, router.middlewares, 0)
	router = router.Group("/api/v1/videos", f, f, f)
	assert.Len(t, router.middlewares, 3)
	assert.Equal(t, router.path, "/api/v1/videos")
	router = router.Group("/playlists", f, f, f)
	assert.Len(t, router.middlewares, 6)
	assert.Equal(t, router.path, "/api/v1/videos/playlists")
}

func TestRouter_Handle(t *testing.T) {
	router := New()
	router = router.Group("/api/v1/videos", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("1"))
			next.ServeHTTP(w, r)
		})
	}, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("2"))
			next.ServeHTTP(w, r)
		})
	})

	router = router.Group("/playlists", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("3"))
			next.ServeHTTP(w, r)
		})
	}, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("4"))
			next.ServeHTTP(w, r)
		})
	})

	router.Handle(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Playlist"))
	}))
	router.Handle(http.MethodGet, "/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Handle"))
	}))

	var (
		r = httptest.NewRequest(http.MethodGet, "/api/v1/videos/playlists/abc", nil)
		w = httptest.NewRecorder()
	)
	router.ServeHTTP(w, r)
	assert.Equal(t, "1234Handle", w.Body.String())

	r = httptest.NewRequest(http.MethodGet, "/api/v1/videos/playlists", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, "1234Playlist", w.Body.String())
}

func TestRouter_HandleFunc(t *testing.T) {
	router := New()
	router.HandleFunc(http.MethodPatch, "/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HandleFunc"))
	})

	var (
		r = httptest.NewRequest(http.MethodPatch, "/abc", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "HandleFunc", w.Body.String())
}

func TestRouter_GET(t *testing.T) {
	router := New()
	router.GET("/get/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET"))
	})

	var (
		r = httptest.NewRequest(http.MethodGet, "/get/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "GET", w.Body.String())
}

func TestRouter_HEAD(t *testing.T) {
	router := New()
	router.HEAD("/head/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HEAD"))
	})

	var (
		r = httptest.NewRequest(http.MethodHead, "/head/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "HEAD", w.Body.String())
}

func TestRouter_OPTIONS(t *testing.T) {
	router := New()
	router.OPTIONS("/options/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OPTIONS"))
	})

	var (
		r = httptest.NewRequest(http.MethodOptions, "/options/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "OPTIONS", w.Body.String())
}

func TestRouter_POST(t *testing.T) {
	router := New()
	router.POST("/post/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("POST"))
	})

	var (
		r = httptest.NewRequest(http.MethodPost, "/post/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "POST", w.Body.String())
}

func TestRouter_PUT(t *testing.T) {
	router := New()
	router.PUT("/put/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PUT"))
	})

	var (
		r = httptest.NewRequest(http.MethodPut, "/put/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "PUT", w.Body.String())
}

func TestRouter_PATCH(t *testing.T) {
	router := New()
	router.PATCH("/patch/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PATCH"))
	})

	var (
		r = httptest.NewRequest(http.MethodPatch, "/patch/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "PATCH", w.Body.String())
}

func TestRouter_DELETE(t *testing.T) {
	router := New()
	router.DELETE("/delete/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("DELETE"))
	})

	var (
		r = httptest.NewRequest(http.MethodDelete, "/delete/id", nil)
		w = httptest.NewRecorder()
	)

	router.ServeHTTP(w, r)

	assert.Equal(t, "DELETE", w.Body.String())
}

func TestRouter_Static_HTTPRouter(t *testing.T) {
	router := NewWithEngine(HTTPRouter)
	assert.Panics(t, func() { router.Static("/filepath", "/var/www") })
	assert.Panics(t, func() { router.Static("/test/filepath", "/var/www") })

	var (
		r = httptest.NewRequest(http.MethodGet, "/assets/README.md", nil)
		w = httptest.NewRecorder()
	)
	router.Static("/assets/*filepath", ".")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, len(w.Body.String()) > 0)
}

func TestRouter_Static_Muxie(t *testing.T) {
	router := NewWithEngine(Muxie)
	assert.Panics(t, func() { router.Static("/filepath", "/var/www") })
	assert.Panics(t, func() { router.Static("/test/filepath", "/var/www") })

	var (
		r = httptest.NewRequest(http.MethodGet, "/assets/README.md", nil)
		w = httptest.NewRecorder()
	)
	router.Static("/assets/*filepath", ".")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, len(w.Body.String()) > 0)
}
