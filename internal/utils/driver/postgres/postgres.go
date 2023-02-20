package postgres

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"os"
)

type IDB interface {
	DB() *bun.DB
	IsConnected() bool
	Close() error
}

type db struct {
	config *Config
	db     *bun.DB
	errors chan error
}

func NewDB(cfg *Config) (IDB, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.GetConnectionString()),
		pgdriver.WithTimeout(cfg.Timeout),
	))

	sqldb.SetMaxOpenConns(cfg.PoolSize)

	client := db{
		config: cfg,
		db:     bun.NewDB(sqldb, pgdialect.New()),
	}

	if cfg.Debug {
		client.db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithWriter(os.Stdout),
		))
	}

	if err := client.db.Ping(); err != nil {
		_ = client.Close()

		return nil, err
	}

	return &client, nil
}

func (db *db) DB() *bun.DB {
	return db.db
}

func (db *db) IsConnected() bool {
	if db == nil {
		return false
	}

	if err := db.db.Ping(); err != nil {
		return false
	}

	return true
}

func (db *db) Close() error {
	if db.db == nil {
		return nil
	}

	return db.db.Close()
}
