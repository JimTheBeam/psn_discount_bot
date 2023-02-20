package daemon

import (
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"psn_discount_bot/internal/config"
	"psn_discount_bot/internal/connector"
	"psn_discount_bot/internal/httpclient"
	"psn_discount_bot/internal/logger"
	"psn_discount_bot/internal/parser"
	"psn_discount_bot/internal/process"
	"psn_discount_bot/internal/repository"
	"psn_discount_bot/internal/service"
	"psn_discount_bot/internal/tgbot"
	"syscall"
)

type Daemon struct {
	bot    *tgbot.TgBot
	config *config.Config
	conn   connector.IConnector
	repo   *repository.Repo
	logger *logger.Entry
	cron   *cron.Cron
	parser *parser.Parser

	priceCheckerProcessor  *process.PriceChecker
	priceInformerProcessor *process.PriceInformer
}

func New(cfg *config.Config, conn connector.IConnector, log *logger.Entry) *Daemon {
	return &Daemon{
		config: cfg,
		conn:   conn,
		repo:   repository.New(conn),
		logger: log,
		cron:   cron.New(),
		parser: parser.NewParser(httpclient.NewClient(&cfg.App.HTTPClient)),
	}
}

func (d *Daemon) Run() {
	d.logger.Info("daemon started")
	defer d.logger.Info("daemon finished")

	serviceClient := service.NewService(d.repo, d.parser, d.logger)

	d.bot = tgbot.New(&d.config.App.Bot, d.logger, serviceClient)

	d.priceCheckerProcessor = process.NewPriceChecker(&d.config.App.Processors.PriceChecker, d.logger, d.repo, d.cron, d.parser)
	d.priceInformerProcessor = process.NewPriceInformer(&d.config.App.Processors.PriceInformer, d.logger, d.repo, d.cron, d.bot)

	// run bot
	go d.bot.Run()

	// run processes
	go d.priceCheckerProcessor.Run()
	go d.priceInformerProcessor.Run()

	d.cron.Start()

	interrupter := make(chan os.Signal, 1)
	signal.Notify(interrupter, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		//todo: other errors
		case <-interrupter:
			d.close()

			return
		}
	}
}

// close останавливает демон управления приложением.
func (d *Daemon) close() {
	d.cron.Stop()

	d.bot.Close()

	d.priceCheckerProcessor.Close()
	d.priceInformerProcessor.Close()

	d.logger.Info("close connections...")

	if err := d.conn.Close(); err != nil {
		d.logger.WithError(err).Error("close connections")
	}

	d.logger.Info("close connections success")
}
