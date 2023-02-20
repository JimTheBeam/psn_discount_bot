package model

import (
	"github.com/uptrace/bun"
	"time"
)

type UsersGames struct {
	bun.BaseModel `bun:"table:users_games,alias:users_games"`

	UserTelegramID    int        `bun:"user_telegram_id,pk"`
	GameID            int        `bun:"game_id,pk"`
	SubscriptionPrice float64    `bun:"subscription_price"`
	CreatedAt         time.Time  `bun:"created_at,nullzero"`
	DeletedAt         *time.Time `bun:"deleted_at"`

	Game *Game `bun:"rel:belongs-to,join:game_id=id"`
	User *User `bun:"rel:belongs-to,join:user_telegram_id=telegram_id"`
}
