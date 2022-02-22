package sqlx

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// Sqlx ...
type Sqlx struct {
	name string
	db   *sqlx.DB
}

// New ...
func New(db *sqlx.DB) *Sqlx {
	return &Sqlx{db: db}
}

// SetName ...
func (s *Sqlx) SetName(name string) {
	s.name = name
}

// Name ...
func (s *Sqlx) Name() string {
	if len(s.name) == 0 {
		return strings.ToUpper(s.db.DriverName())
	}
	return s.name
}

// Check ...
func (s *Sqlx) Check() error {
	return s.db.Ping()
}
