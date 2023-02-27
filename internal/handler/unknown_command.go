package handler

import (
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) UnknownCommand(c tgbot.Context) error {
	chatID := c.Update().Message.Chat.ID
	fromMessageID := c.Update().Message.MessageID

	responseText := "I don't know that command"

	msg := c.Bot().NewMessage(chatID, responseText)
	msg.ReplyToMessageID = fromMessageID

	return c.Bot().SendMessage(msg)
}
