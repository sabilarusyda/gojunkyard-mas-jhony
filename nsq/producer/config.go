package nsqproducer

// Config holds all configuration of nsqd and asynchronous
type Config struct {
	NSQD    string `envconfig:"NSQD"`
	IsAsync bool   `envconfig:"IS_ASYNC"`
}

// NewConfig returns pointer of Config object
func NewConfig(addr string, isAsync bool) *Config {
	return &Config{
		NSQD:    addr,
		IsAsync: false, // there's a bug in async producer, so it's disabled
	}
}
