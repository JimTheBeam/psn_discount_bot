package process

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type PriceCheckerConfig struct {
	Disabled bool   `yaml:"disabled"`
	CronTime string `yaml:"cron_time"`
}

func (cfg *PriceCheckerConfig) Validate() error {
	if _, err := cron.ParseStandard(cfg.CronTime); err != nil {
		return fmt.Errorf("cron_time parse: %w", err)
	}

	return nil
}
