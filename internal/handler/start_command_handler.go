package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) StartCommand(c tgbot.Context) error {
	chatID := c.Update().Message.Chat.ID
	userID := int(c.Update().Message.From.ID)

	newUser := model.User{
		TelegramID:        userID,
		TelegramChatID:    int(chatID),
		TelegramFirstName: c.Update().Message.From.FirstName,
		TelegramLastName:  c.Update().Message.From.LastName,
		TelegramUserName:  c.Update().Message.From.UserName,
	}

	h.service.CreateUser(newUser)

	text := "Hello from psn_discounter_bot. I send /help for more information"

	msg := tgbotapi.NewMessage(chatID, text)

	msg.ReplyMarkup = model.NewSubscriptionListKeyboard()

	return c.Bot().SendMessage(msg)
}
