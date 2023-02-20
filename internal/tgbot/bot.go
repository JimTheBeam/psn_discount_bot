package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/logger"
	"psn_discount_bot/internal/model"
)

type Service interface {
	SubscribeToGame(userID int, url string) (string, error)
	CreateUser(currentUser model.User)
}

type TgBot struct {
	Bot     *tgbotapi.BotAPI
	service Service
	log     *logger.Entry
	cfg     *Config
	quit    chan struct{}
}

func New(cfg *Config, log *logger.Entry, service Service) *TgBot {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		//todo:
		log.Panic(err)
	}

	bot.Debug = cfg.Debug

	log.Infof("Authorized on account %s", bot.Self.UserName)

	return &TgBot{
		Bot:     bot,
		service: service,
		log:     log,
		cfg:     cfg,
		quit:    make(chan struct{}),
	}
}

func (b *TgBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.UpdatesTimeout

	updates := b.Bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			// command
			if update.Message.IsCommand() {
				b.CommandRouter(update)

				continue
			}

			b.subscribeToGame(update)

		case <-b.quit:
			return
		}
	}
}

func (b *TgBot) Close() {
	b.quit <- struct{}{}
}
