package model

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:users"`

	ID                int       `bun:"id,pk,autoincrement"`
	TelegramID        int       `bun:"telegram_id"`
	TelegramChatID    int       `bun:"telegram_chat_id"`
	TelegramFirstName string    `bun:"telegram_first_name"`
	TelegramLastName  string    `bun:"telegram_last_name"`
	TelegramUserName  string    `bun:"telegram_user_name"`
	CreatedAt         time.Time `bun:"created_at,nullzero"`

	Games []Game `bun:"m2m:users_games,join:User=Game"`
}
