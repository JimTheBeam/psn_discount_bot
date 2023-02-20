package service

import (
	"errors"
	"fmt"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
)

var (
	ErrInvalidURL = errors.New("url is invalid")
	ErrInternal   = errors.New("internal error")
)

func (s *Service) SubscribeToGame(data payload.Subscribe) (string, error) {
	game, err := s.repo.GetGameByID(data.GameID)
	if err != nil {
		s.log.WithError(err).WithField("game_id", data.GameID).Error("get game")

		return "", ErrInternal
	}

	if game == nil {
		s.log.WithField("game_id", data.GameID).Error("game is nil")

		return "", ErrInternal
	}

	isSubscribed, err := s.repo.IsSubscribed(game.ID, data.UserID)
	if err != nil {
		s.log.WithError(err).WithField("game_id", data.GameID).
			WithField("user_id", data.UserID).
			Error("is subscribed")

		return "", ErrInternal
	}

	responseMsg := fmt.Sprintf("You have subscribed!\nGame: %s\nPrice notified: %.2f",
		game.Name,
		data.Price,
	)

	if isSubscribed {
		return responseMsg, nil
	}

	sub := model.UsersGames{
		UserTelegramID:    data.UserID,
		GameID:            data.GameID,
		SubscriptionPrice: data.Price,
	}

	if err := s.repo.Subscribe(sub); err != nil {
		s.log.WithError(err).WithField("game_id", data.GameID).Error("subscribe to a game")

		return "", ErrInternal
	}

	return responseMsg, nil
}

func (s *Service) GetGame(userID int, url string) (model.Game, *model.UsersGames, error) {
	game, err := s.repo.GetGameByURL(url)
	if err != nil {
		s.log.WithError(err).WithField("url", url).Error("get game")

		return model.Game{}, nil, ErrInternal
	}

	if game == nil {
		parsedGame, err := s.parser.ParseGame(url)
		if err != nil {
			s.log.WithError(err).WithField("url", url).Error("parse game")

			return model.Game{}, nil, ErrInvalidURL
		}

		game, err = s.repo.CreateGame(parsedGame)
		if err != nil {
			s.log.WithError(err).WithField("url", url).Error("create game")

			return model.Game{}, nil, ErrInternal
		}

		game.SetGameIDToPrices()

		game.Prices, err = s.repo.CreatePrices(game.Prices)
		if err != nil {
			s.log.WithError(err).WithField("game_id", game.ID).Error("create prices")

			return model.Game{}, nil, ErrInternal
		}
	}

	if game == nil {
		s.log.WithError(err).WithField("url", url).Error("game is nil")

		return model.Game{}, nil, ErrInternal
	}

	subscription, err := s.repo.GetSubscription(game.ID, userID)
	if err != nil {
		s.log.WithError(err).WithField("game_id", game.ID).
			WithField("user_id", userID).
			Error("get subscription")

		return model.Game{}, nil, ErrInternal
	}

	return *game, subscription, nil
}

func (s *Service) Unsubscribe(userID, gameID int) (string, error) {
	game, err := s.repo.GetGameByID(gameID)
	if err != nil {
		s.log.WithError(err).WithField("game_id", gameID).Error("get game by id")

		return "", ErrInternal
	}

	if game == nil {
		s.log.WithError(err).WithField("game_id", gameID).Error("game not found")

		return "", ErrInternal
	}

	if err := s.repo.Unsubscribe(gameID, userID); err != nil {
		s.log.WithError(err).WithField("game_id", gameID).Error("unsubscribe")

		return "", ErrInternal
	}

	responseMsg := fmt.Sprintf("Unsubscribed!\nGame: %s", game.Name)

	return responseMsg, nil
}

func (s *Service) GetSubscriptions(data payload.Subscriptions) ([]model.UsersGames, error) {
	subscriptions, err := s.repo.GetSubscriptionList(data)
	if err != nil {
		s.log.WithError(err).WithField("user_id", data.UserID).Error("get subscriptions")

		return nil, ErrInternal
	}

	return subscriptions, nil
}
