package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"psn_discount_bot/internal/model"
)

func (r *Repo) GetGameByURL(url string) (*model.Game, error) {
	var game model.Game

	err := r.pg.DB().NewSelect().Model(&game).
		Where("url = ?", url).
		Relation("Prices").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("pg: %w", err)
	}

	return &game, nil
}

func (r *Repo) GetGamesWithChangedPrice() ([]model.Game, error) {
	var games []model.Game

	err := r.pg.DB().NewSelect().Model(&games).
		Distinct().
		Join("JOIN users_games AS ug ON games.id =ug.game_id").
		Join("JOIN prices AS p ON games.id =p.game_id").
		Where("p.value < ug.subscription_price").
		Relation("Prices").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("pg: %w", err)
	}

	return games, nil
}

func (r *Repo) UpsertGame(game model.Game) (model.Game, error) {
	_, err := r.pg.DB().NewInsert().Model(&game).
		On(`CONFLICT (url) DO UPDATE SET current_price=EXCLUDED.current_price`).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return model.Game{}, fmt.Errorf("pg: %w", err)
	}

	return game, nil
}

func (r *Repo) CreateGame(game model.Game) (*model.Game, error) {
	_, err := r.pg.DB().NewInsert().Model(&game).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	return &game, nil
}

func (r *Repo) GetAllGames() ([]model.Game, error) {
	var games []model.Game

	err := r.pg.DB().NewSelect().Model(&games).
		Where("deleted_at IS NULL").
		Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	return games, nil
}

func (r *Repo) IsSubscribed(gameID, userTgID int) (bool, error) {
	exist, err := r.pg.DB().NewSelect().Model(&model.UsersGames{}).
		Where("user_telegram_id = ?", userTgID).
		Where("game_id = ?", gameID).
		Where("deleted_at IS NULL").
		Exists(context.Background())
	if err != nil {
		return false, fmt.Errorf("pg: %w", err)
	}

	return exist, nil
}

func (r *Repo) Subscribe(sub model.UsersGames) error {
	err := r.pg.DB().NewInsert().Model(&sub).Scan(context.Background())
	if err != nil {
		return fmt.Errorf("pg: %w", err)
	}

	return nil
}
