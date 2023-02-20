package main

import (
	"flag"
	"fmt"
	"log"
	"psn_discount_bot/internal/config"
	"psn_discount_bot/internal/connector"
	"psn_discount_bot/internal/daemon"
	"psn_discount_bot/internal/logger"
)

func main() {
	cfgPath := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ParseConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.Validate(); err != nil {
		log.Fatal(fmt.Errorf("validate config: %w", err))
	}

	// logger
	logg := logger.New().SetService("psn_bot").SetLevel(cfg.App.Log.Level).NewEntry()

	conn, err := connector.New(&cfg.Connections)
	if err != nil {
		log.Fatal(fmt.Errorf("create connections: %w", err))
	}

	d := daemon.New(&cfg, conn, logg)
	d.Run()
}
