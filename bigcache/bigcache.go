package bigcache

import (
	"time"

	"github.com/allegro/bigcache"
)

// Config holds all bigcache configuration
type Config struct {
	Shards             int           `envconfig:"SHARDS"`
	LifeWindow         time.Duration `envconfig:"LIFE_WINDOW"`
	CleanWindow        time.Duration `envconfig:"CLEAN_WINDOW"`
	MaxEntriesInWindow int           `envconfig:"MAX_ENTRIES_IN_WINDOW"`
	MaxEntrySize       int           `envconfig:"MAX_ENTRY_SIZE"`
	Verbose            bool          `envconfig:"VERBOSE"`
	HardMaxCacheSize   int           `envconfig:"HARD_MAX_CACHE_SIZE"`
}

// New instantiate bigcache based on configuration
func New(cfg Config) (*bigcache.BigCache, error) {
	bigcacheConfig := parseConfig(&cfg)
	return bigcache.NewBigCache(bigcacheConfig)
}

func parseConfig(cfg *Config) bigcache.Config {
	return bigcache.Config{
		Shards:             cfg.Shards,
		LifeWindow:         cfg.LifeWindow,
		CleanWindow:        cfg.CleanWindow,
		MaxEntriesInWindow: 1000 * int(cfg.LifeWindow/time.Second),
		MaxEntrySize:       cfg.MaxEntrySize,
		Verbose:            cfg.Verbose,
		HardMaxCacheSize:   cfg.HardMaxCacheSize,
	}
}
