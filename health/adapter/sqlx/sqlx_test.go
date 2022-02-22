package sqlx

import (
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		want = &Sqlx{db: nil}
		got  = New(nil)
	)
	assert.Equal(t, want, got)
}

func TestRedigo_GetSetName(t *testing.T) {
	db, _, _ := sqlmock.New()

	sqlx := New(sqlx.NewDb(db, "mysql"))
	assert.Equal(t, "MYSQL", sqlx.Name())

	const name = "mysql.svc.local"
	sqlx.SetName(name)
	assert.Equal(t, name, sqlx.Name())
}

func TestRedigo_Check(t *testing.T) {
	db, _, _ := sqlmock.New()

	sqlx := New(sqlx.NewDb(db, "mysql"))
	assert.Nil(t, sqlx.Check())
}
