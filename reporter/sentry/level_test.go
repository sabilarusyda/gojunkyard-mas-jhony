package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isDebug(t *testing.T) {
	assert.Equal(t, true, isDebug(0))
	assert.Equal(t, true, isDebug(DEBUG))
	assert.Equal(t, false, isDebug(INFO))
	assert.Equal(t, false, isDebug(WARNING))
	assert.Equal(t, false, isDebug(ERROR))
	assert.Equal(t, false, isDebug(FATAL))
	assert.Equal(t, false, isDebug(128))
}

func Test_isInfo(t *testing.T) {
	assert.Equal(t, true, isInfo(0))
	assert.Equal(t, true, isInfo(DEBUG))
	assert.Equal(t, true, isInfo(INFO))
	assert.Equal(t, false, isInfo(WARNING))
	assert.Equal(t, false, isInfo(ERROR))
	assert.Equal(t, false, isInfo(FATAL))
	assert.Equal(t, false, isInfo(128))
}

func Test_isWarning(t *testing.T) {
	assert.Equal(t, true, isWarning(0))
	assert.Equal(t, true, isWarning(DEBUG))
	assert.Equal(t, true, isWarning(INFO))
	assert.Equal(t, true, isWarning(WARNING))
	assert.Equal(t, false, isWarning(ERROR))
	assert.Equal(t, false, isWarning(FATAL))
	assert.Equal(t, false, isWarning(128))
}

func Test_isError(t *testing.T) {
	assert.Equal(t, true, isError(0))
	assert.Equal(t, true, isError(DEBUG))
	assert.Equal(t, true, isError(INFO))
	assert.Equal(t, true, isError(WARNING))
	assert.Equal(t, true, isError(ERROR))
	assert.Equal(t, false, isError(FATAL))
	assert.Equal(t, false, isError(128))
}

func Test_isFatal(t *testing.T) {
	assert.Equal(t, true, isFatal(0))
	assert.Equal(t, true, isFatal(DEBUG))
	assert.Equal(t, true, isFatal(INFO))
	assert.Equal(t, true, isFatal(WARNING))
	assert.Equal(t, true, isFatal(ERROR))
	assert.Equal(t, true, isFatal(FATAL))
	assert.Equal(t, false, isFatal(128))
}
