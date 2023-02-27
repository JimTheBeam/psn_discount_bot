package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) HelpCommand(c tgbot.Context) error {
	chatID := c.Update().Message.Chat.ID

	text := `Send me url from psn store and I will tell you when the price reduces`

	msg := tgbotapi.NewMessage(chatID, text)

	msg.ReplyMarkup = model.NewSubscriptionListKeyboard()

	return c.Bot().SendMessage(msg)
}
