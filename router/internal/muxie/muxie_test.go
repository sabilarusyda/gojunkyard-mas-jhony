package muxie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"devcode.xeemore.com/systech/gojunkyard/router/internal/param"

	"github.com/kataras/muxie"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		got  = New()
		want = &Router{muxie.NewMux()}
	)

	// cannot compare "got" & "want" directly due to a function variable
	assert.Equal(t, got.PathCorrection, want.PathCorrection)
	assert.Equal(t, got.Routes, want.Routes)
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

func TestRouter_Handle_MultipleMethod(t *testing.T) {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param.GetParam(r, "id")
		w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	})

	router := New()
	router.Handle(http.MethodGet, "/users/:id", f)
	router.Handle(http.MethodPost, "/users/:id", f)

	// condition 1. GET
	var (
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/users/23", nil)
	)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "23", w.Body.String())

	// condition 2. POST
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/users/24", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "24", w.Body.String())
}

func TestRouter_Lookup(t *testing.T) {
	const route = "/users"
	const method = http.MethodGet

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param.GetParam(r, "id")
		w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	})

	muxie := New()
	assert.False(t, muxie.Lookup(method, route))

	muxie.Handle(method, route, f)
	assert.True(t, muxie.Lookup(method, route))
}
