package lm

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_getStatusCode(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "status 0",
			args: args{
				status: 0,
			},
			want: 200,
		},
		{
			name: "status 200",
			args: args{
				status: 200,
			},
			want: 200,
		},
		{
			name: "status 300",
			args: args{
				status: 300,
			},
			want: 300,
		},
		{
			name: "status 400",
			args: args{
				status: 400,
			},
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getStatusCode(tt.args.status))
		})
	}
}

func Test_getRemoteAddr(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "exist x-forwarded-for",
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/", nil)
					r.Header.Set("X-Forwarded-for", "127.0.0.1")
					return r
				}(),
			},
			want: "127.0.0.1",
		},
		{
			name: "not exist x-forwarded-for",
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/", nil)
					r.RemoteAddr = "127.0.0.1:89888"
					return r
				}(),
			},
			want: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getRemoteAddr(tt.args.r), "getRemoteAddr() is not equal")
		})
	}
}

func Test_getEvent(t *testing.T) {
	logger := zerolog.New(ioutil.Discard)

	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want *zerolog.Event
	}{
		{
			name: "status ok",
			args: args{
				status: 200,
			},
			want: logger.Info(),
		},
		{
			name: "status redirect",
			args: args{
				status: 300,
			},
			want: logger.Info(),
		},
		{
			name: "status bad request",
			args: args{
				status: 400,
			},
			want: logger.Warn(),
		},
		{
			name: "status internal server error",
			args: args{
				status: 500,
			},
			want: logger.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := getEvent(&logger, tt.args.status)
			assert.Equal(t, tt.want, event, "getEvent() is not equal")
		})
	}
}
