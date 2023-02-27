package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"strings"
)

func (b *Bot) CallbackRouter(update tgbotapi.Update) {
	method := getMethodFromCallback(update.CallbackData())
	payloadStr := getPayloadFromCallback(update.CallbackData(), method)

	ctx := b.NewContextWithPayload(update, payloadStr)

	handlerFunc := b.getCallbackHandlerFunc(method)

	if err := handlerFunc(ctx); err != nil {
		b.log.WithError(err).WithField("method", method).Error("handle callback")
	}

	return
}

func (b *Bot) getCallbackHandlerFunc(method string) HandlerFunc {
	switch method {
	case model.SubscribeCallbackData:
		return b.handler.SubscribeToGameCallback

	case model.UnsubscribeCallbackData:
		return b.handler.UnsubscribeToGameCallback

	case model.GameCallbackData:
		return b.handler.GetGameCallback

	case model.SubscriptionListCallbackData:
		return b.handler.GetSubscriptionsCallback

	case model.CancelCallbackData:
		return b.handler.CancelCallback

	default:
		return b.handler.UnknownCallback
	}
}

func getMethodFromCallback(callback string) string {
	split := strings.Split(callback, model.CallbackDelimiter)
	method := split[0]

	return method
}

func getPayloadFromCallback(callback, method string) string {
	return strings.TrimPrefix(callback, method+model.CallbackDelimiter)
}
