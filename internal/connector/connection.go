package connector

import (
	"errors"
	"fmt"
	"psn_discount_bot/internal/utils/driver/postgres"
)

type IConnector interface {
	Pg() postgres.IDB

	CheckConnections() error
	Close() error
}

type connector struct {
	pg postgres.IDB
}

func New(cfg *Config) (IConnector, error) {
	conn := new(connector)
	var err error

	if conn.pg, err = postgres.NewDB(&cfg.Postgres); err != nil {
		return nil, fmt.Errorf("connect to pg: %w", err)
	}

	return conn, nil
}

func (conn *connector) CheckConnections() error {
	if ok := conn.pg.IsConnected(); !ok {
		return errors.New("postgres: connection is lost")
	}

	return nil
}

func (conn *connector) Pg() postgres.IDB {
	return conn.pg
}

func (conn *connector) Close() error {
	if err := conn.pg.Close(); err != nil {
		return fmt.Errorf("postgres: %w", err)
	}

	return nil
}
