package mysql

import (
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/jmoiron/sqlx"
)

// DB DB struct
type DB struct {
	DSN             string
	Dialect         string
	MaxConnOpen     int
	MaxConnLifetime time.Duration
	MaxConnIdle     int
	Client          *sqlx.DB
	Breaker         *breaker.Breaker
}

// Config holds per database config
// DSN				: connection string. It usually contains the username, password, and address
// MaxConnOpen		: maximum number of open connections to the database. Then there is no limit on the number of open connections. Default: 0
// MaxConnDuration	: maximum amount of time a connection may be reused. Expired connections may be closed lazily before reuse. If duration <= 0, connections are reused forever. Default value: 0
// MaxConnIdle		: maximum number of connections in the idle connection pool. If MaxConnIdle > 0 && < MaxIdleConns then the MaxConnIdle will be reduced to MaxConnOpen limit. Default value: 0
type Config struct {
	Dialect         string        `envconfig:"DIALECT"`
	Host            string        `envconfig:"HOST"`
	Port            int           `envconfig:"PORT"`
	Name            string        `envconfig:"NAME"`
	Username        string        `envconfig:"USER_NAME"`
	Password        string        `envconfig:"PASSWORD"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE"`
}
