package model

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const (
	CancelCallbackData           = "cancel"
	SubscribeCallbackData        = "subscribe"
	UnsubscribeCallbackData      = "unsubscribe"
	SubscriptionListCallbackData = "subscriptions"
	GameCallbackData             = "game"
)

const CallbackDelimiter = ":"

func NewGameKeyboardWithPrices(prices []Price) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(prices))

	for i := range prices {
		callback := fmt.Sprintf("%s%s%d%s%.2f",
			SubscribeCallbackData,
			CallbackDelimiter,
			prices[i].GameID,
			CallbackDelimiter,
			prices[i].Value,
		)

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(prices[i].GetPriceText(), callback),
		)

		rows = append(rows, row)
	}

	rows = append(rows, newSubscriptionsRow())

	cancelRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cancel", CancelCallbackData))

	rows = append(rows, cancelRow)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func NewGameUnsubscribeKeyboard(gameID int) tgbotapi.InlineKeyboardMarkup {
	callbackData := UnsubscribeCallbackData + CallbackDelimiter + strconv.Itoa(gameID)
	cancelRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Unsubscribe", callbackData))

	return tgbotapi.NewInlineKeyboardMarkup(cancelRow, newSubscriptionsRow())
}

func NewSubscriptionListKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		newSubscriptionsRow(),
	)
}

func NewSubscriptionsListKeyboard(subs []UsersGames) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(subs))

	for i := range subs {
		callback := fmt.Sprintf("%s%s%d",
			GameCallbackData,
			CallbackDelimiter,
			subs[i].GameID,
		)

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subs[i].Game.Name, callback),
		)

		rows = append(rows, row)
	}

	//todo: navigation buttons
	cancelRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cancel", CancelCallbackData))

	rows = append(rows, cancelRow)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func newSubscriptionsRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("To my subscriptions", SubscriptionListCallbackData),
	)
}
