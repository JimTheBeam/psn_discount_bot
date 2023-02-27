package payload

import (
	"errors"
	"fmt"
	"psn_discount_bot/internal/model"
	"strconv"
	"strings"
)

type (
	Subscribe struct {
		UserID int
		GameID int
		Price  float64
	}

	Subscriptions struct {
		UserID int
		Limit  int
		Offset int
	}
)

func BindSubscribePayload(data string, userID int) (Subscribe, error) {
	splitPayload := strings.Split(data, model.CallbackDelimiter)

	if len(splitPayload) != 2 {
		return Subscribe{}, errors.New("split is not equal 2")
	}

	gameID, err := strconv.Atoi(splitPayload[0])
	if err != nil {
		return Subscribe{}, fmt.Errorf("parse game_id: %w", err)
	}

	price, err := strconv.ParseFloat(splitPayload[1], 64)
	if err != nil {
		return Subscribe{}, fmt.Errorf("parse price: %w", err)
	}

	return Subscribe{
		UserID: userID,
		GameID: gameID,
		Price:  price,
	}, nil
}

func BindUnsubscribePayload(data string) (int, error) {
	gameID, err := strconv.Atoi(data)
	if err != nil {
		return 0, fmt.Errorf("parse game_id: %w", err)
	}

	return gameID, nil
}

func BindGetGamePayload(data string) (int, error) {
	gameID, err := strconv.Atoi(data)
	if err != nil {
		return 0, fmt.Errorf("parse game_id: %w", err)
	}

	return gameID, nil
}
