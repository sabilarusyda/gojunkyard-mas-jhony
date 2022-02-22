package httprouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"devcode.xeemore.com/systech/gojunkyard/router/internal/param"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, &Router{httprouter.New()}, New())
}

func TestRouter_Handle_Found(t *testing.T) {
	router := New()
	router.Handle(http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param.GetParam(r, "id")
		w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	}))

	var (
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/users/23", nil)
	)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "23", w.Body.String())
}

func TestRouter_Handle_NotFound(t *testing.T) {
	router := New()
	router.Handle(http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param.GetParam(r, "id")
		w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	}))

	var (
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/users/23/abc", nil)
	)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "404 page not found\n", w.Body.String())
}

func TestRouter_Lookup(t *testing.T) {
	router := New()
	router.Handle(http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param.GetParam(r, "id")
		w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	}))
	router.Handle(http.MethodPost, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param.GetParam(r, "id")
		w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	}))

	assert.True(t, router.Lookup(http.MethodGet, "/users/:id"))
	assert.True(t, router.Lookup(http.MethodPost, "/users/:id"))
	assert.False(t, router.Lookup(http.MethodOptions, "/users/:id"))
	assert.False(t, router.Lookup(http.MethodDelete, "/users/:id"))
}
