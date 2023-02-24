package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/logger"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
)

type Service interface {
	SubscribeToGame(data payload.Subscribe) (string, error)
	Unsubscribe(userID, gameID int) (model.Game, error)
	GetSubscriptions(data payload.Subscriptions) ([]model.UsersGames, error)
	CreateUser(currentUser model.User)

	GetGame(userID int, url string) (model.Game, *model.UsersGames, error)
	GetGameByID(userID, gameID int) (model.Game, *model.UsersGames, error)
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

			if update.CallbackQuery != nil {
				b.CallbackRouter(update)
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

func (b *TgBot) Close() {
	b.quit <- struct{}{}
}
