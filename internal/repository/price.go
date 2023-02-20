package repository

import (
	"context"
	"fmt"
	"psn_discount_bot/internal/model"
)

//todo:
func (r *Repo) CreatePrices(prices []model.Price) ([]model.Price, error) {
	_, err := r.pg.DB().NewInsert().Model(&prices).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	return prices, nil
}

func (r *Repo) DeleteGamePrices(gameID int) error {
	_, err := r.pg.DB().NewDelete().Model(&[]model.Price{}).
		Where("game_id = ?", gameID).
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("pg: %w", err)
	}

	return nil
}
