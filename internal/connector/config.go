package connector

import (
	"fmt"
	"psn_discount_bot/internal/utils/driver/postgres"
)

type (
	Config struct {
		Postgres postgres.Config `yaml:"postgres"`
	}
)

func (c *Config) Validate() error {
	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres.%w", err)
	}

	return nil
}
