package radix

import (
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		pool = new(radix.Pool)
		want = &Radix{pool: pool}
		got  = New(pool)
	)
	assert.Equal(t, want, got)
}

func TestRadix_GetSetName(t *testing.T) {
	redigo := New(nil)
	assert.Equal(t, "REDIS", redigo.Name())

	const name = "Redis Account 127.0.0.1"
	redigo.SetName(name)
	assert.Equal(t, name, redigo.Name())
}
