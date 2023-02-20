package process

import (
	"github.com/robfig/cron/v3"
	"psn_discount_bot/internal/logger"
)

type Bot interface {
	SendText(chatID int64, message string)
}

type PriceInformer struct {
	config   *PriceInformerConfig
	logger   *logger.Entry
	cron     *cron.Cron
	repo     Repository
	bot      Bot
	quitChan chan struct{}
}

func NewPriceInformer(
	cfg *PriceInformerConfig,
	log *logger.Entry,
	repo Repository,
	cron *cron.Cron,
	bot Bot,
) *PriceInformer {
	return &PriceInformer{
		config:   cfg,
		logger:   log,
		cron:     cron,
		repo:     repo,
		bot:      bot,
		quitChan: make(chan struct{}),
	}
}

func (p *PriceInformer) Run() {
	if p.config.Disabled {
		return
	}

	p.logger.Infof("price_informer start with crone time: %s", p.config.CronTime)

	entryID, err := p.cron.AddFunc(p.config.CronTime, p.informAboutPriceChange)
	if err != nil {
		p.logger.WithError(err).Error("start cron price_informer")

		return
	}

	p.logger.Infof("price_informer cron %d started", entryID)

	for range p.quitChan {
		p.logger.Info("closing price_informer...")

		return
	}
}

func (p *PriceInformer) informAboutPriceChange() {
	p.logger.Debug("price_informer start")

	games, err := p.repo.GetGamesWithChangedPrice()
	if err != nil {
		p.logger.WithError(err).Error("get all games")

		return
	}

	if len(games) == 0 {
		p.logger.Debug("no games for parsing")

		return
	}

	for i := range games {
		users, err := p.repo.GetUserGamesWithDifferentPrice(games[i].ID, games[i].GetMinPrice())
		if err != nil {
			p.logger.WithError(err).WithField("game_id", games[i].ID).Error("get users for game")

			continue
		}

		if len(users) == 0 {
			continue
		}

		for u := range users {
			message := games[i].PriceChangedText()

			p.bot.SendText(int64(users[u].User.TelegramChatID), message)
		}
	}
}

func (p *PriceInformer) Close() {
	if p.config.Disabled {
		return
	}

	p.quitChan <- struct{}{}
}
