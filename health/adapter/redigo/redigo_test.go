package redigo

import (
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		pool = new(redis.Pool)
		want = &Redigo{pool: pool}
		got  = New(pool)
	)
	assert.Equal(t, want, got)
}

func TestRedigo_GetSetName(t *testing.T) {
	redigo := New(nil)
	assert.Equal(t, "REDIS", redigo.Name())

	const name = "Redis Account 127.0.0.1"
	redigo.SetName(name)
	assert.Equal(t, name, redigo.Name())
}

func TestRedigo_Check(t *testing.T) {
	redigo := New(&redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn := redigomock.NewConn()
			conn.Command("PING").Expect("PONG")
			return conn, nil
		},
	})
	assert.Equal(t, nil, redigo.Check())

	redigo = New(&redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn := redigomock.NewConn()
			conn.Command("PING").ExpectError(redis.ErrPoolExhausted)
			return conn, nil
		},
	})
	assert.Equal(t, redis.ErrPoolExhausted, redigo.Check())
}
