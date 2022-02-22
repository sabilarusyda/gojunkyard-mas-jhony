package param

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetParam(t *testing.T) {
	const key = "===KEY==="
	const value = "===VALUE==="
	var r = httptest.NewRequest(http.MethodGet, "/", nil)
	r = SetParam(r, key, value)
	assert.Equal(t, value, GetParam(r, key))
}
