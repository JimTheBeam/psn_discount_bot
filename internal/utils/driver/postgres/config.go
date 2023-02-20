package postgres

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmptyAddr     = errors.New("addr is empty")
	ErrEmptyDatabase = errors.New("database is empty")
	ErrEmptyUser     = errors.New("user is empty")
	ErrEmptyPassword = errors.New("password is empty")
	ErrEmptyPoolSize = errors.New("pool_size is empty")
	ErrEmptyTimeout  = errors.New("timeout is empty")
)

type Config struct {
	Addr     string        `yaml:"addr" json:"addr"`
	Database string        `yaml:"database" json:"database"`
	User     string        `yaml:"username" json:"user"`
	Password string        `yaml:"password" json:"password"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout"`
	Debug    bool          `yaml:"debug" json:"debug"`
	PoolSize int           `yaml:"pool_size" json:"pool_size"`
}

func (cfg *Config) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Addr, cfg.Database)
}

func (cfg *Config) Validate() error {
	if cfg.Addr == "" {
		return ErrEmptyAddr
	}

	if cfg.Database == "" {
		return ErrEmptyDatabase
	}

	if cfg.User == "" {
		return ErrEmptyUser
	}

	if cfg.Password == "" {
		return ErrEmptyPassword
	}

	if cfg.PoolSize <= 0 {
		return ErrEmptyPoolSize
	}

	if cfg.Timeout <= 0 {
		return ErrEmptyTimeout
	}

	return nil
}
