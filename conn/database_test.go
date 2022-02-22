package conn

import (
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestInitDB(t *testing.T) {
	db, _, _ := sqlmock.New()
	type mock struct {
		sqlxConnect func(string, string) (*sqlx.DB, error)
	}
	type args struct {
		cfg DBConfig
	}
	tests := []struct {
		name    string
		args    args
		mock    mock
		want    *sqlx.DB
		wantErr bool
	}{
		{
			name: "Invalid DSN",
			args: args{
				cfg: DBConfig{},
			},
			wantErr: true,
		},
		{
			name: "MySQL connection failed",
			args: args{
				cfg: DBConfig{
					Dialect:         "mysql",
					Host:            "127.0.0.1",
					Port:            3306,
					Name:            "v2_videos",
					Username:        "root",
					Password:        "root",
					MaxConnOpen:     1,
					MaxConnLifetime: time.Second,
					MaxConnIdle:     1,
				},
			},
			mock: mock{
				sqlxConnect: func(dialect, dsn string) (*sqlx.DB, error) {
					const wantDialect = "mysql"
					const wantDSN = "root:root@tcp(127.0.0.1:3306)/v2_videos?parseTime=true"
					if wantDialect != dialect {
						log.Fatalf("error dialect. expected: [%s], got: [%s]", wantDialect, dialect)
					}
					if wantDSN != dsn {
						log.Fatalf("error dsn. expected: [%s], got: [%s]", wantDSN, dsn)
					}
					return nil, sql.ErrConnDone
				},
			},
			wantErr: true,
		},
		{
			name: "Postgres connection failed",
			args: args{
				cfg: DBConfig{
					Dialect:         "postgres",
					Host:            "127.0.0.1",
					Port:            3306,
					Name:            "v2_videos",
					Username:        "root",
					Password:        "root",
					MaxConnOpen:     1,
					MaxConnLifetime: time.Second,
					MaxConnIdle:     1,
				},
			},
			mock: mock{
				sqlxConnect: func(dialect, dsn string) (*sqlx.DB, error) {
					const wantDialect = "postgres"
					const wantDSN = "postgres://root:root@127.0.0.1:3306/v2_videos?sslmode=disable"
					if wantDialect != dialect {
						log.Fatalf("error dialect. expected: [%s], got: [%s]", wantDialect, dialect)
					}
					if wantDSN != dsn {
						log.Fatalf("error dsn. expected: [%s], got: [%s]", wantDSN, dsn)
					}
					return nil, sql.ErrConnDone
				},
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				cfg: DBConfig{
					Dialect:         "mysql",
					Host:            "127.0.0.1",
					Port:            3306,
					Name:            "v2_videos",
					Username:        "root",
					Password:        "root",
					MaxConnOpen:     1,
					MaxConnLifetime: time.Second,
					MaxConnIdle:     1,
				},
			},
			mock: mock{
				sqlxConnect: func(dialect, dsn string) (*sqlx.DB, error) {
					const wantDialect = "mysql"
					const wantDSN = "root:root@tcp(127.0.0.1:3306)/v2_videos?parseTime=true"
					if wantDialect != dialect {
						log.Fatalf("error dialect. expected: [%s], got: [%s]", wantDialect, dialect)
					}
					if wantDSN != dsn {
						log.Fatalf("error dsn. expected: [%s], got: [%s]", wantDSN, dsn)
					}
					return sqlx.NewDb(db, "sqlmock"), nil
				},
			},
			want:    sqlx.NewDb(db, "sqlmock"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlxConnect = tt.mock.sqlxConnect
			got, err := InitDB(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
