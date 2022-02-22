package mysql

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/jmoiron/sqlx"

	_ "github.com/denisenkom/go-mssqldb"
)

// sqlxConnect is used to mock sqlx.Open
var sqlxConnect = sqlx.Connect

// New init the database from config to database connection
func New(cfg Config) (*DB, error) {
	dsn, err := cfg.toDSN()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := sqlxConnect(cfg.Dialect, dsn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	client.SetMaxOpenConns(cfg.MaxConnOpen)
	client.SetMaxIdleConns(cfg.MaxConnIdle)
	client.SetConnMaxLifetime(cfg.MaxConnLifetime)

	cb := breaker.New(10, 1, 10*time.Second)

	db := DB{
		DSN:             dsn,
		Dialect:         cfg.Dialect,
		MaxConnOpen:     cfg.MaxConnOpen,
		MaxConnLifetime: cfg.MaxConnLifetime,
		MaxConnIdle:     cfg.MaxConnIdle,
		Client:          client,
		Breaker:         cb,
	}
	return &db, nil
}

func (cfg Config) toDSN() (string, error) {
	switch cfg.Dialect {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name), nil
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name), nil
	case "sqlserver":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=60", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name), nil
	default:
		return "", fmt.Errorf("Dialect is not supported. expected: (mysql|postgres|sqlserver), got: %s", cfg.Dialect)
	}
}

// Ping Function to ping the current connection of MySQL
func (db *DB) Ping() bool {
	if db.Client == nil {
		return false
	}

	err := db.Client.Ping()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Close Function to close database connection
func (db *DB) Close() {
	if db.Client != nil {
		db.Client.Close()
	}
}

// Connected Function to check connection to MySQL DB
func (db *DB) Connected() bool {
	res := db.Breaker.Run(func() error {
		if !db.Ping() {
			return errors.New("failed ping to database")
		}
		return nil
	})
	return res != breaker.ErrBreakerOpen
}
