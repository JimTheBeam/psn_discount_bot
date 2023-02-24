package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
	"psn_discount_bot/internal/service"
	"strconv"
	"strings"
)

func (b *TgBot) CommandRouter(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		b.startCommand(update)

		return

	case "subscription_list":
		b.getSubscriptions(chatID, userID)

		return

	case "help":
		msg.Text = `Send me url from psn store and I will tell you when the price reduces`
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "I don't know that command"
	}

	b.SendMessage(msg)
}

func (b *TgBot) CallbackRouter(update tgbotapi.Update) {
	var err error
	chatID := update.CallbackQuery.Message.Chat.ID
	fromMessageID := update.CallbackQuery.Message.MessageID
	userID := int(update.CallbackQuery.Message.Chat.ID)

	split := strings.Split(update.CallbackData(), model.CallbackDelimiter)
	method := split[0]

	switch method {
	case model.SubscribeCallbackData:
		if len(split) != 3 {
			b.SendText(chatID, service.ErrInternal.Error())
		}

		gameID, err := strconv.Atoi(split[1])
		if err != nil {
			b.SendText(chatID, service.ErrInternal.Error())
		}

		price, err := strconv.ParseFloat(split[2], 64)
		if err != nil {
			b.SendText(chatID, service.ErrInternal.Error())
		}

		data := payload.Subscribe{
			UserID: userID,
			GameID: gameID,
			Price:  price,
		}

		responseText, err := b.service.SubscribeToGame(data)
		if err != nil {
			b.SendText(chatID, err.Error())
			return
		}

		replyMarkup := model.NewGameUnsubscribeKeyboard(gameID)

		msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, responseText, replyMarkup)

		b.SendMessage(msg)

	case model.UnsubscribeCallbackData:
		var gameID int
		if len(split) > 1 {
			gameID, err = strconv.Atoi(split[1])
			if err != nil {
				b.SendText(chatID, service.ErrInternal.Error())
			}
		}

		game, err := b.service.Unsubscribe(userID, gameID)
		if err != nil {
			b.SendText(chatID, err.Error())
			return
		}

		text, replyMarkup := getTextAndMarkupForGame(game, nil)

		msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, text, replyMarkup)

		b.SendMessage(msg)

	case model.GameCallbackData:
		var gameID int
		if len(split) > 1 {
			gameID, err = strconv.Atoi(split[1])
			if err != nil {
				b.SendText(chatID, service.ErrInternal.Error())
			}
		}

		game, subscription, err := b.service.GetGameByID(userID, gameID)
		if err != nil {
			b.SendText(chatID, err.Error())
			return
		}

		text, replyMarkup := getTextAndMarkupForGame(game, subscription)

		msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, fromMessageID, text, replyMarkup)

		b.SendMessage(msg)

	case model.SubscriptionListCallbackData:
		b.getSubscriptions(chatID, int64(userID))

	case model.CancelCallbackData:
		msg := tgbotapi.NewEditMessageText(chatID, fromMessageID, update.CallbackQuery.Message.Text)
		b.SendMessage(msg)
	}
}
