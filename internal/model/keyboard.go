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

	cancelRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cancel", CancelCallbackData))

	rows = append(rows, cancelRow)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func NewGameUnsubscribeKeyboard(gameID int) tgbotapi.InlineKeyboardMarkup {
	callbackData := UnsubscribeCallbackData + CallbackDelimiter + strconv.Itoa(gameID)
	cancelRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Unsubscribe", callbackData))

	return tgbotapi.NewInlineKeyboardMarkup(cancelRow)
}

func NewSubscriptionListKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get my subscriptions", SubscriptionListCallbackData),
		),
	)
}
