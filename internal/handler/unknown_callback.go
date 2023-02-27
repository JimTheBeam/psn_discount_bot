package handler

import (
	"psn_discount_bot/internal/service"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) UnknownCallback(c tgbot.Context) error {
	chatID := c.Update().CallbackQuery.Message.Chat.ID
	fromMessageID := c.Update().CallbackQuery.Message.MessageID

	responseText := service.ErrInternal.Error()

	msg := c.Bot().NewMessage(chatID, responseText)
	msg.ReplyToMessageID = fromMessageID

	return c.Bot().SendMessage(msg)
}
