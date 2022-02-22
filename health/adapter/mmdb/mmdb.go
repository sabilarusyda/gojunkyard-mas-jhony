package mmdb

import (
	"errors"

	"github.com/oschwald/maxminddb-golang"
)

// MMDB ...
type MMDB struct {
	name   string
	reader *maxminddb.Reader
}

// New ...
func New(reader *maxminddb.Reader) *MMDB {
	return &MMDB{reader: reader}
}

// SetName ...
func (m *MMDB) SetName(name string) {
	m.name = name
}

// Name ...
func (m *MMDB) Name() string {
	if len(m.name) == 0 {
		return "MMDB"
	}
	return m.name
}

// Check ...
func (m *MMDB) Check() error {
	if m.reader == nil {
		return errors.New("NOT READY")
	}
	return nil
}
