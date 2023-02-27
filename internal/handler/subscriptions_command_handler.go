package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) GetSubscriptionsCommand(c tgbot.Context) error {
	chatID := c.Update().Message.Chat.ID
	messageID := c.Update().Message.MessageID
	userID := int(c.Update().Message.Chat.ID)

	data := payload.Subscriptions{
		UserID: userID,
		Limit:  0,
		Offset: 0,
	}

	subs, err := h.service.GetSubscriptions(data)
	if err != nil {
		return c.Bot().SendText(chatID, err.Error())
	}

	if len(subs) == 0 {
		return c.Bot().SendText(chatID, "You haven't got subscriptions yet.")
	}

	text := "Subscriptions:"
	replyMarkup := model.NewSubscriptionsListKeyboard(subs)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	msg.ReplyMarkup = replyMarkup

	return c.Bot().SendMessage(msg)
}
