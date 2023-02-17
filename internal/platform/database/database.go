package database

import (
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DbConfig struct {
	Host            string
	User            string
	Password        string
	Database        string
	ApplicationName string
}

func OpenConnection(cfg DbConfig) *bun.DB {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithUser(cfg.User),
		pgdriver.WithPassword(cfg.Password),
		pgdriver.WithApplicationName(cfg.ApplicationName),
		pgdriver.WithAddr(cfg.Host),
		pgdriver.WithDatabase("gamereview"),
		pgdriver.WithTimeout(3*time.Second),
		pgdriver.WithDialTimeout(3*time.Second),
		pgdriver.WithReadTimeout(3*time.Second),
		pgdriver.WithWriteTimeout(3*time.Second),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithInsecure(true),
	)

	sqlDb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqlDb, pgdialect.New())

	_, err := db.Exec("SELECT 1")

	if err != nil {
		panic("error connecting to database - " + err.Error())
	}

	return db
}
