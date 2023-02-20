package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"psn_discount_bot/internal/model"
)

func (r *Repo) GetUserGamesWithDifferentPrice(gameID int, currentPrice float64) ([]model.UsersGames, error) {
	var userGames []model.UsersGames

	err := r.pg.DB().NewSelect().Model(&userGames).
		Where("?TableAlias.game_id = ?", gameID).
		Where("?TableAlias.subscription_price != ?", currentPrice).
		Relation("User").
		Scan(context.Background())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {

		return nil, fmt.Errorf("pg: %w", err)
	}

	return userGames, nil
}
