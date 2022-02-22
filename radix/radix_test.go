package radix

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseConfig(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		lenWant int
	}{
		{
			name: "filled all",
			args: args{
				cfg: &Config{
					Host:                   "127.0.0.1",
					Port:                   6379,
					MaxConnection:          100,
					Database:               1,
					PingInterval:           time.Second,
					PipelineWindowDeadline: time.Microsecond * 100,
					PipelineWindowLimit:    1000,
				},
			},
			lenWant: 3,
		},
		{
			name: "not filled",
			args: args{
				cfg: &Config{},
			},
			lenWant: 2,
		},
		{
			name: "invalid config",
			args: args{
				cfg: &Config{
					PipelineWindowDeadline: -time.Second,
					PipelineWindowLimit:    -10,
				},
			},
			lenWant: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseConfig(tt.args.cfg)
			assert.Len(t, got, tt.lenWant)
		})
	}
}
