package service

import (
	"psn_discount_bot/internal/logger"
	"psn_discount_bot/internal/model"
)

type Repository interface {
	// game
	GetGameByURL(url string) (*model.Game, error)
	GetGameByID(gameID int) (*model.Game, error)
	CreateGame(game model.Game) (*model.Game, error)
	IsSubscribed(gameID, userTgID int) (bool, error)
	Subscribe(sub model.UsersGames) error
	Unsubscribe(gameID, userTelegramID int) error
	GetSubscription(gameID, userTgID int) (*model.UsersGames, error)
	CreatePrices(prices []model.Price) ([]model.Price, error)

	// user
	GetUserByTgID(telegramID int) (*model.User, error)
	CreateUser(user model.User) (*model.User, error)
	UpsertUser(user model.User) (model.User, error)
}

type Parser interface {
	ParseGame(url string) (model.Game, error)
}

type Service struct {
	repo   Repository
	parser Parser
	log    *logger.Entry
}

func NewService(repo Repository, parser Parser, log *logger.Entry) *Service {
	return &Service{
		repo:   repo,
		parser: parser,
		log:    log,
	}
}
