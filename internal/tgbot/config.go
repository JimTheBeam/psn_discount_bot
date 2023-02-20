package tgbot

import "errors"

type Config struct {
	Token          string `yaml:"token"`
	Debug          bool   `yaml:"debug"`
	UpdatesTimeout int    `yaml:"updates_timeout"`
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("token is empty")
	}

	if c.UpdatesTimeout <= 0 {
		return errors.New("updates_timeout must be greater than zero")
	}

	return nil
}
