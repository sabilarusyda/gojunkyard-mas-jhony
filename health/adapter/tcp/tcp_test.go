package tcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	const addr = "127.0.0.1:9090"
	var (
		want = &TCP{addr: addr}
		got  = New(addr)
	)
	assert.Equal(t, want, got)
}

func TestTCP_GetSetName(t *testing.T) {
	const addr = "127.0.0.1:9090"
	redigo := New(addr)
	assert.Equal(t, "TCP", redigo.Name())

	const name = "TCP " + addr
	redigo.SetName(name)
	assert.Equal(t, name, redigo.Name())
}
