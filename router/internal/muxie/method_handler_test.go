package muxie

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/muxie"
	"github.com/stretchr/testify/assert"
)

func TestMethodHandler(t *testing.T) {
	mux := New()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		fmt.Fprintf(w, "GET: List all users\n")
	})

	mux.Mux.Handle("/user/:id", Methods().
		HandleFunc(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "GET: User details by user ID: %s\n", muxie.GetParam(w, "id"))
		}).
		HandleFunc(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "POST: save user with ID: %s\n", muxie.GetParam(w, "id"))
		}).
		HandleFunc(http.MethodDelete, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "DELETE: remove user with ID: %s\n", muxie.GetParam(w, "id"))
		}))

	srv := httptest.NewServer(mux)
	defer srv.Close()

	var w *httptest.ResponseRecorder
	var r *http.Request

	// condition 1. GET /users
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/users", nil)
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GET: List all users\n", w.Body.String())

	// condition 2. POST /users
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/users", nil)
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "Method Not Allowed\n", w.Body.String())

	// condition 3. GET /user/42
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/user/42", nil)
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GET: User details by user ID: 42\n", w.Body.String())

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/user/42", nil)
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "POST: save user with ID: 42\n", w.Body.String())

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodDelete, "/user/42", nil)
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "DELETE: remove user with ID: 42\n", w.Body.String())

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPut, "/user/42", nil)
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "Method Not Allowed\n", w.Body.String())

}

func TestMethodHandler_Lookup(t *testing.T) {
	mux := New()

	mux.Handle(http.MethodGet, "/user/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET: User details by user ID: %s\n", muxie.GetParam(w, "id"))
	}))
	mux.Handle(http.MethodPost, "/user/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "POST: save user with ID: %s\n", muxie.GetParam(w, "id"))
	}))
	mux.Handle(http.MethodDelete, "/user/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "DELETE: remove user with ID: %s\n", muxie.GetParam(w, "id"))
	}))

	assert.True(t, mux.Lookup(http.MethodGet, "/user/:id"))
	assert.True(t, mux.Lookup(http.MethodPost, "/user/:id"))
	assert.True(t, mux.Lookup(http.MethodDelete, "/user/:id"))
	assert.False(t, mux.Lookup(http.MethodHead, "/user/:id"))
	assert.False(t, mux.Lookup(http.MethodPatch, "/user/:id"))
}
