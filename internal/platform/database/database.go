package database

import (
	"crypto/tls"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type DbConfig struct {
	Host            string
	User            string
	Password        string
	Database        string
	ApplicationName string
}

func OpenConnection(cfg *DbConfig) *bun.DB {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithApplicationName("go_gamereview"),
		pgdriver.WithAddr(cfg.Host),
		pgdriver.WithTimeout(3*time.Second),
		pgdriver.WithDialTimeout(3*time.Second),
		pgdriver.WithReadTimeout(3*time.Second),
		pgdriver.WithWriteTimeout(3*time.Second),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)

	sqlDb := sql.OpenDB(pgconn)
	return bun.NewDB(sqlDb, pgdialect.New())
}
