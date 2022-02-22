package radix

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	_radix "github.com/mediocregopher/radix/v3"
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

func parseConfig(cfg *Config) []_radix.PoolOpt {
	var (
		dialOpts = make([]_radix.DialOpt, 0, 1)
		poolOpts = make([]_radix.PoolOpt, 0, 3)
	)

	if cfg.Database > 0 {
		dialOpts = append(dialOpts, _radix.DialSelectDB(cfg.Database))
	}

	if cfg.PingInterval > 0 {
		poolOpts = append(poolOpts, _radix.PoolPingInterval(cfg.PingInterval))
	}

	if cfg.PipelineWindowDeadline >= 0 && cfg.PipelineWindowLimit >= 0 {
		poolOpts = append(poolOpts, _radix.PoolPipelineWindow(cfg.PipelineWindowDeadline, cfg.PipelineWindowLimit))
	}

	return append(poolOpts, _radix.PoolConnFunc(func(network, addr string) (_radix.Conn, error) {
		return _radix.Dial(network, addr, dialOpts...)
	}))
}

// New Instantiate radix pool based on configuration
func New(cfg Config) (*Radix, error) {
	var (
		addr     = cfg.Host + ":" + strconv.Itoa(cfg.Port)
		poolOpts = parseConfig(&cfg)
		cb       = breaker.New(10, 1, 10*time.Second)
	)

	radixPool, err := _radix.NewPool("tcp", addr, cfg.MaxConnection, poolOpts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	radix := Radix{
		Address:       addr,
		MaxConnection: cfg.MaxConnection,
		Pool:          radixPool,
		PoolOpts:      poolOpts,
		Breaker:       cb,
	}

	return &radix, err
}

// Ping Function to ping the current connection of Redis
func (radix *Radix) Ping() bool {
	if radix.Pool == nil {
		return false
	}

	err := radix.Pool.Do(_radix.Cmd(nil, "PING"))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// Close Function to close database connection
func (radix *Radix) Close() {
	if radix.Pool != nil {
		radix.Pool.Close()
	}
}

// CheckAvailableConns Function to check available connections
func (radix *Radix) CheckAvailableConns() bool {
	if radix.Pool != nil {
		availableConn := radix.Pool.NumAvailConns()
		if availableConn > 0 {
			return true
		}
	}

	return false
}

// Connected Function to check connection to Redis
func (radix *Radix) Connected() bool {
	res := radix.Breaker.Run(func() error {
		if !radix.Ping() {
			return errors.New("failed ping to redis")
		}
		return nil
	})
	return res != breaker.ErrBreakerOpen
}
