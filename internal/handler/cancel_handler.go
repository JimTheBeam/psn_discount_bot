package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) CancelCallback(c tgbot.Context) error {
	chatID := c.Update().CallbackQuery.Message.Chat.ID
	fromMessageID := c.Update().CallbackQuery.Message.MessageID
	text := c.Update().CallbackQuery.Message.Text

	msg := tgbotapi.NewEditMessageText(chatID, fromMessageID, text)

	return c.Bot().SendMessage(msg)
}
