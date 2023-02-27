package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/config"
	"psn_discount_bot/internal/logger"
)

type HandlerFunc func(c Context) error

type Bot struct {
	Bot     *tgbotapi.BotAPI
	handler Handler
	log     *logger.Entry
	cfg     *config.BotConfig
	quit    chan struct{}
}

func New(cfg *config.BotConfig, log *logger.Entry, handler Handler) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		//todo:
		log.Panic(err)
	}

	bot.Debug = cfg.Debug

	log.Infof("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		Bot:     bot,
		handler: handler,
		log:     log,
		cfg:     cfg,
		quit:    make(chan struct{}),
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.UpdatesTimeout

	updates := b.Bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:

			if update.CallbackQuery != nil {
				b.CallbackRouter(update)

				continue
			}

			if update.Message == nil {
				continue
			}

			// command
			if update.Message.IsCommand() {
				b.CommandRouter(update)

				continue
			}

			b.getGame(update)

		case <-b.quit:
			return
		}
	}
}

func (b *Bot) Close() {
	b.quit <- struct{}{}
}
