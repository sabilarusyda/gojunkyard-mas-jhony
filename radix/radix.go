package radix

import (
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v3"
)

// Config holds all radix pool configuration for connecting to redis
type Config struct {
	Host                   string        `envconfig:"HOST"`
	Port                   int           `envconfig:"PORT"`
	MaxConnection          int           `envconfig:"MAX_CONNECTION"`
	Database               int           `envconfig:"DATABASE"`
	PingInterval           time.Duration `envconfig:"PING_INTERVAL"`
	PipelineWindowDeadline time.Duration `envconfig:"PIPELINE_WINDOW_DURATION"`
	PipelineWindowLimit    int           `envconfig:"PIPELINE_WINDOW_LIMIT"`
}

func parseConfig(cfg *Config) []radix.PoolOpt {
	var (
		dialOpts = make([]radix.DialOpt, 0, 1)
		poolOpts = make([]radix.PoolOpt, 0, 3)
	)

	if cfg.Database > 0 {
		dialOpts = append(dialOpts, radix.DialSelectDB(cfg.Database))
	}

	if cfg.PingInterval > 0 {
		poolOpts = append(poolOpts, radix.PoolPingInterval(cfg.PingInterval))
	}

	if cfg.PipelineWindowDeadline >= 0 && cfg.PipelineWindowLimit >= 0 {
		poolOpts = append(poolOpts, radix.PoolPipelineWindow(cfg.PipelineWindowDeadline, cfg.PipelineWindowLimit))
	}

	return append(poolOpts, radix.PoolConnFunc(func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr, dialOpts...)
	}))
}

// New instantiate radix pool based on configuration
func New(cfg Config) (*radix.Pool, error) {
	var (
		addr     = cfg.Host + ":" + strconv.Itoa(cfg.Port)
		poolOpts = parseConfig(&cfg)
	)
	return radix.NewPool("tcp", addr, cfg.MaxConnection, poolOpts...)
}
