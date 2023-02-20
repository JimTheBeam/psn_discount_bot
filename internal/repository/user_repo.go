package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"psn_discount_bot/internal/model"
)

func (r *Repo) GetUserByTgID(telegramID int) (*model.User, error) {
	var user model.User

	err := r.pg.DB().NewSelect().Model(&user).Where("telegram_id = ?", telegramID).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("pg: %w", err)
	}

	return &user, nil
}

func (r *Repo) CreateUser(user model.User) (*model.User, error) {
	_, err := r.pg.DB().NewInsert().Model(&user).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	return &user, nil
}

func (r *Repo) UpsertUser(user model.User) (model.User, error) {
	_, err := r.pg.DB().NewInsert().Model(&user).
		On(`CONFLICT (telegram_id) DO UPDATE SET
					telegram_chat_id=EXCLUDED.telegram_chat_id, 
					telegram_first_name=EXCLUDED.telegram_first_name, 
					telegram_last_name=EXCLUDED.telegram_last_name, 
					telegram_user_name=EXCLUDED.telegram_user_name`).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return model.User{}, fmt.Errorf("pg: %w", err)
	}

	return user, nil
}
