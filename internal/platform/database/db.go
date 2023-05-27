package database

import (
	"context"
	"crypto/tls"
	"database/sql"
	"embed"
	"log"
	"runtime"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
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
		pgdriver.WithDatabase(cfg.Database),
		pgdriver.WithTimeout(3*time.Second),
		pgdriver.WithDialTimeout(3*time.Second),
		pgdriver.WithReadTimeout(3*time.Second),
		pgdriver.WithWriteTimeout(3*time.Second),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithInsecure(true),
	)

	sqlDb := sql.OpenDB(pgconn)
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqlDb.SetMaxOpenConns(maxOpenConns)
	sqlDb.SetMaxIdleConns(maxOpenConns)

	db := bun.NewDB(sqlDb, pgdialect.New())

	_, err := db.Exec("SELECT 1")
	if err != nil {
		panic("error connecting to database - " + err.Error())
	}

	return db
}

//go:embed migrations/*.sql
var sqlMigrations embed.FS

func Migrate(ctx context.Context, db *bun.DB) error {
	migrations := migrate.NewMigrations()
	err := migrations.Discover(sqlMigrations)
	if err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	migrator.Init(ctx)

	migrator.Lock(ctx)
	defer migrator.Unlock(ctx)

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.ID == 0 {
		log.Println("no new migrations to run")
	}

	return nil
}
