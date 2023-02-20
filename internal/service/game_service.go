package service

import (
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"psn_discount_bot/internal/model"
)

var (
	ErrInvalidURL        = errors.New("url is invalid")
	ErrAlreadySubscribed = errors.New("already subscribed to this game")
	ErrInternal          = errors.New("internal error")
)

func (s *Service) SubscribeToGame(userID int, url string) (string, error) {
	game, err := s.repo.GetGameByURL(url)
	if err != nil {
		s.log.WithError(err).WithField("url", url).Error("get game")

		return "", ErrInternal
	}

	if game == nil {
		parsedGame, err := s.parser.ParseGame(url)
		if err != nil {
			s.log.WithError(err).WithField("url", url).Error("parse game")

			return "", ErrInvalidURL
		}

		game, err = s.repo.CreateGame(parsedGame)
		if err != nil {
			s.log.WithError(err).WithField("url", url).Error("create game")

			return "", ErrInternal
		}

		game.SetGameIDToPrices()

		game.Prices, err = s.repo.CreatePrices(game.Prices)
		if err != nil {
			s.log.WithError(err).WithField("game_id", game.ID).Error("create prices")

			return "", ErrInternal
		}
	}

	if game == nil {
		s.log.WithError(err).WithField("url", url).Error("game is nil")

		return "", ErrInternal
	}

	isSubscribed, err := s.repo.IsSubscribed(game.ID, userID)
	if err != nil {
		s.log.WithError(err).WithField("url", url).Error("is subscribed")

		return "", ErrInternal
	}

	if isSubscribed {
		return "", ErrAlreadySubscribed
	}

	sub := model.UsersGames{
		BaseModel:         bun.BaseModel{},
		UserTelegramID:    userID,
		GameID:            game.ID,
		SubscriptionPrice: game.GetMinPrice(),
	}

	if err := s.repo.Subscribe(sub); err != nil {
		s.log.WithError(err).WithField("url", url).Error("subscribe to a game")

		return "", ErrInternal
	}

	responseMsg := fmt.Sprintf("You have subscribed!\nGame: %s\nPrice:\n%s",
		game.Name,
		game.GetPriceText(),
	)

	return responseMsg, nil
}
