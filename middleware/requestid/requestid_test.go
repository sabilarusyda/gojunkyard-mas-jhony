package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Case 1. RequestId middleware is not set
	assert.Equal(t, "", GetFromRequest(nil))
	assert.Equal(t, "", GetFromContext(nil))

	// Case 2. RequestId is not forwarded from nginx
	{
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			w = httptest.NewRecorder()
			f = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id1 := GetFromRequest(r)
				id2 := GetFromContext(r.Context())
				assert.Len(t, id1, 32)
				assert.Equal(t, id1, id2)
			})
		)
		New()(f).ServeHTTP(w, r)
	}

	// Case 3. RequestId forwarded from another service or nginx
	{
		const expected = "0685d19c7a3741f09369271af0750eed"
		const header = "X-Request-Id"
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add(header, expected)
		var (
			w = httptest.NewRecorder()
			f = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id1 := GetFromRequest(r)
				id2 := GetFromContext(r.Context())
				assert.Equal(t, expected, id1)
				assert.Equal(t, expected, id2)
			})
		)
		New()(f).ServeHTTP(w, r)
		assert.Equal(t, expected, w.Header().Get(header))
	}
}
