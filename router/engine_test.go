package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRouterEngine(t *testing.T) {
	assert.NotNil(t, getRouterEngine(HTTPRouter))
	assert.NotNil(t, getRouterEngine(Muxie))
	assert.NotPanics(t, func() { getRouterEngine(HTTPRouter) })
	assert.NotPanics(t, func() { getRouterEngine(Muxie) })
	assert.PanicsWithValue(t, "Router engine is not found", func() { getRouterEngine(20) })
}
