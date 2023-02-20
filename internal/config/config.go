package config

import (
	"errors"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"psn_discount_bot/internal/connector"
	"psn_discount_bot/internal/httpclient"
	"psn_discount_bot/internal/process"
	"psn_discount_bot/internal/tgbot"
)

type (
	Config struct {
		App         App              `yaml:"app"`
		Connections connector.Config `yaml:"connections"`
	}

	App struct {
		Bot        tgbot.Config      `yaml:"bot"`
		Log        log               `yaml:"log"`
		HTTPClient httpclient.Config `yaml:"http_client"`
		Processors Processors        `yaml:"processors"`
	}

	Processors struct {
		PriceChecker  process.PriceCheckerConfig  `yaml:"price_checker"`
		PriceInformer process.PriceInformerConfig `yaml:"price_informer"`
	}
)

type log struct {
	Level string `yaml:"level"`
}

func (c *Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return fmt.Errorf("app.%w", err)
	}

	if err := c.Connections.Validate(); err != nil {
		return fmt.Errorf("connections.%w", err)
	}

	return nil
}

func (c *App) Validate() error {
	if err := c.Bot.Validate(); err != nil {
		return fmt.Errorf("bot.%w", err)
	}

	if err := c.Log.validate(); err != nil {
		return fmt.Errorf("log.%w", err)
	}

	if err := c.HTTPClient.Validate(); err != nil {
		return fmt.Errorf("http_client.%w", err)
	}

	if err := c.Processors.Validate(); err != nil {
		return fmt.Errorf("processors.%w", err)
	}

	return nil
}

func (c *Processors) Validate() error {
	if err := c.PriceChecker.Validate(); err != nil {

		return fmt.Errorf("price_checker.%w", err)
	}

	if err := c.PriceInformer.Validate(); err != nil {
		return fmt.Errorf("price_informer.%w", err)
	}

	return nil
}

func (l *log) validate() error {
	switch l.Level {
	case "panic",
		"fatal",
		"error",
		"warn", "warning",
		"info",
		"debug",
		"trace":
		return nil
	default:
		return errors.New("level is empty or invalid")
	}
}

// ParseConfig - функция получения и обработки конфига.
func ParseConfig(path string) (Config, error) {
	var config Config

	filename, err := filepath.Abs(path)
	if err != nil {
		return config, err
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return config, err
	}

	return config, nil
}
