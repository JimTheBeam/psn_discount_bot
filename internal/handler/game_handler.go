package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
	"psn_discount_bot/internal/service"
	"psn_discount_bot/internal/tgbot"
)

func (h *Handler) GetGameFromMessage(c tgbot.Context) error {
	chatID := c.Update().Message.Chat.ID
	userID := int(c.Update().Message.From.ID)
	messageID := c.Update().Message.MessageID
	gameURL := c.Update().Message.Text

	if gameURL == "" {
		text := "URL is empty. Press /help for information"

		return c.Bot().SendTextInReplyToMessage(chatID, messageID, text)
	}

	game, subscription, err := h.service.GetGame(userID, gameURL)
	if err != nil {
		return c.Bot().SendTextInReplyToMessage(chatID, messageID, err.Error())
	}

	msg := tgbotapi.NewMessage(chatID, "")
	msg.ReplyToMessageID = messageID

	msg.Text, msg.ReplyMarkup = getTextAndMarkupForGame(game, subscription)

	return c.Bot().SendMessage(msg)
}

func (h *Handler) GetGameCallback(c tgbot.Context) error {
	chatID := c.Update().CallbackQuery.Message.Chat.ID
	fromMessageID := c.Update().CallbackQuery.Message.MessageID
	userID := int(c.Update().CallbackQuery.Message.Chat.ID)

	gameID, err := payload.BindGetGamePayload(c.Payload())
	if err != nil {
		return c.Bot().SendText(chatID, service.ErrInternal.Error())
	}

	game, subscription, err := h.service.GetGameByID(userID, gameID)
	if err != nil {
		return c.Bot().SendText(chatID, err.Error())
	}

	text, replyMarkup := getTextAndMarkupForGame(game, subscription)

	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, text, replyMarkup)

	return c.Bot().SendMessage(msg)
}

func getTextAndMarkupForGame(game model.Game, subscription *model.UsersGames) (string, tgbotapi.InlineKeyboardMarkup) {
	var (
		text        string
		replyMarkup tgbotapi.InlineKeyboardMarkup
	)

	if subscription == nil {
		text = game.Name + "\nChoose price you want to subscribe\n\n" + game.Url
		replyMarkup = model.NewGameKeyboardWithPrices(game.Prices)
	} else {
		text = fmt.Sprintf("%s\nCurrent price:\n%s\n\nSubscription price: %.2f\n\n%s",
			game.Name, game.GetPriceText(),
			subscription.SubscriptionPrice,
			game.Url,
		)

		replyMarkup = model.NewGameUnsubscribeKeyboard(game.ID)
	}

	return text, replyMarkup
}
