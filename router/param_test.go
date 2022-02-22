package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"devcode.xeemore.com/systech/gojunkyard/router/internal/param"
)

func TestGetParam(t *testing.T) {
	const key = "===KEY==="
	const value = "===VALUE==="
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = param.SetParam(r, key, value)
	assert.Equal(t, value, GetParam(r, key))
}
