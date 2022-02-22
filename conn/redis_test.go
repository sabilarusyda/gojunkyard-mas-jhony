package conn

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func Test_InitRedis(t *testing.T) {
	type mock struct {
		dial func(network, address string, options ...redis.DialOption) (redis.Conn, error)
	}
	type args struct {
		cfg RedisConfig
	}
	tests := []struct {
		name     string
		args     args
		mock     mock
		wantPool *redis.Pool
		wantErr  error
	}{
		{
			name: "Default Config",
			mock: mock{
				dial: func(network, address string, options ...redis.DialOption) (redis.Conn, error) {
					if network != "tcp" {
						log.Fatalln("Network must be tcp")
					}
					if address != ":0" {
						log.Fatalln("Address must be empty string")
					}
					return nil, errors.New("cannot connect")
				},
			},
			wantPool: nil,
			wantErr:  errors.New("cannot connect"),
		},
		{
			name: "Success",
			args: args{
				cfg: RedisConfig{
					Host:      "127.0.0.1",
					Port:      6379,
					MaxActive: 1000,
					Wait:      true,
				},
			},
			mock: mock{
				dial: func(network, address string, options ...redis.DialOption) (redis.Conn, error) {
					const wantAddress = "127.0.0.1:6379"
					if network != "tcp" {
						log.Fatalln("Network must be tcp")
					}
					if address != wantAddress {
						log.Fatalf("Address: %s, want: %s\n", address, wantAddress)
					}
					conn := redigomock.NewConn()
					conn.GenericCommand("PING").Expect("PONG")
					return conn, nil
				},
			},
			wantPool: &redis.Pool{
				MaxActive:       1000,
				Wait:            true,
				MaxIdle:         1000,
				IdleTimeout:     2 * time.Second,
				MaxConnLifetime: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisDial = tt.mock.dial

			gotPool, gotErr := InitRedis(tt.args.cfg)
			if gotPool != nil {
				assert.Equal(t, tt.wantPool.MaxActive, gotPool.MaxActive)
				assert.Equal(t, tt.wantPool.Wait, gotPool.Wait)
				assert.Equal(t, tt.wantPool.IdleTimeout, gotPool.IdleTimeout)
				assert.Equal(t, tt.wantPool.MaxConnLifetime, gotPool.MaxConnLifetime)
				assert.Equal(t, tt.wantPool.MaxIdle, gotPool.MaxIdle)
			}
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
