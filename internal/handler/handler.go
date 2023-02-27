package handler

import (
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

type Handler struct {
	service Service
	log     *logger.Entry
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
		log:     nil,
	}
}
