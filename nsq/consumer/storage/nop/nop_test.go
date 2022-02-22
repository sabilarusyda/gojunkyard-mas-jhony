package redigo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNop(t *testing.T) {
	const key = "ahuehuehe:ahuhuhu"
	nop := New()
	ok, err := nop.SetNX(key, time.Second)
	assert.Nil(t, err)
	assert.True(t, ok)
	err = nop.Delete(key)
	assert.Nil(t, err)
}
