package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
	"psn_discount_bot/internal/service"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) SubscribeToGameCallback(c tgbot.Context) error {
	chatID := c.Update().CallbackQuery.Message.Chat.ID
	fromMessageID := c.Update().CallbackQuery.Message.MessageID
	userID := int(c.Update().CallbackQuery.Message.Chat.ID)

	data, err := payload.BindSubscribePayload(c.Payload(), userID)
	if err != nil {
		return c.Bot().SendText(chatID, service.ErrInternal.Error())
	}

	responseText, err := h.service.SubscribeToGame(data)
	if err != nil {
		return c.Bot().SendText(chatID, err.Error())
	}

	replyMarkup := model.NewGameUnsubscribeKeyboard(data.GameID)

	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, responseText, replyMarkup)

	return c.Bot().SendMessage(msg)
}

func (h *Handler) UnsubscribeToGameCallback(c tgbot.Context) error {
	chatID := c.Update().CallbackQuery.Message.Chat.ID
	fromMessageID := c.Update().CallbackQuery.Message.MessageID
	userID := int(c.Update().CallbackQuery.Message.Chat.ID)

	gameID, err := payload.BindUnsubscribePayload(c.Payload())
	if err != nil {
		return c.Bot().SendText(chatID, service.ErrInternal.Error())
	}

	game, err := h.service.Unsubscribe(userID, gameID)
	if err != nil {
		return c.Bot().SendText(chatID, err.Error())
	}

	text, replyMarkup := getTextAndMarkupForGame(game, nil)

	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, text, replyMarkup)

	return c.Bot().SendMessage(msg)
}

func (h *Handler) GetSubscriptionsCallback(c tgbot.Context) error {
	chatID := c.Update().CallbackQuery.Message.Chat.ID
	fromMessageID := c.Update().CallbackQuery.Message.MessageID
	userID := int(c.Update().CallbackQuery.Message.Chat.ID)

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

	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, text, replyMarkup)

	return c.Bot().SendMessage(msg)
}
