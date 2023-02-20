package process

import (
	"github.com/robfig/cron/v3"
	"psn_discount_bot/internal/logger"
	"psn_discount_bot/internal/model"
)

type Repository interface {
	GetAllGames() ([]model.Game, error)
	CreatePrices(prices []model.Price) ([]model.Price, error)
	DeleteGamePrices(gameID int) error

	//todo: price_informer
	GetGamesWithChangedPrice() ([]model.Game, error)
	GetUserGamesWithDifferentPrice(gameID int, currentPrice float64) ([]model.UsersGames, error)
}

type Parser interface {
	ParseGame(url string) (model.Game, error)
}

type PriceChecker struct {
	config   *PriceCheckerConfig
	logger   *logger.Entry
	cron     *cron.Cron
	parser   Parser
	repo     Repository
	quitChan chan struct{}
}

func NewPriceChecker(
	cfg *PriceCheckerConfig,
	log *logger.Entry,
	repo Repository,
	cron *cron.Cron,
	parser Parser,
) *PriceChecker {
	return &PriceChecker{
		config:   cfg,
		logger:   log,
		cron:     cron,
		parser:   parser,
		repo:     repo,
		quitChan: make(chan struct{}),
	}
}

func (p *PriceChecker) Run() {
	if p.config.Disabled {
		return
	}

	p.logger.Infof("price_checker start with crone time: %s", p.config.CronTime)

	entryID, err := p.cron.AddFunc(p.config.CronTime, p.parseCurrentPrices)
	if err != nil {
		p.logger.WithError(err).Error("start cron price_checker")

		return
	}

	p.logger.Infof("price_checker cron %d started", entryID)

	for range p.quitChan {
		p.logger.Info("closing price_checker...")

		return
	}
}

func (p *PriceChecker) parseCurrentPrices() {
	p.logger.Info("price_checker starting")
	defer p.logger.Debug("price_checker stopped")

	games, err := p.repo.GetAllGames()
	if err != nil {
		p.logger.WithError(err).Error("get all games")

		return
	}

	if len(games) == 0 {
		p.logger.Debug("no games for parsing")

		return
	}

	for i := range games {
		parsedGame, err := p.parser.ParseGame(games[i].Url)
		if err != nil {
			p.logger.WithError(err).WithField("url", games[i].Url).Error("parse game")

			continue
		}

		if err := p.repo.DeleteGamePrices(games[i].ID); err != nil {
			p.logger.WithError(err).WithField("game_id", games[i].ID).Error("delete prices")

			continue
		}

		parsedGame.ID = games[i].ID

		parsedGame.SetGameIDToPrices()

		if _, err := p.repo.CreatePrices(parsedGame.Prices); err != nil {
			p.logger.WithError(err).WithField("game_id", games[i].ID).Error("create prices")

			continue
		}
	}
}

func (p *PriceChecker) Close() {
	if p.config.Disabled {
		return
	}

	p.quitChan <- struct{}{}
}
