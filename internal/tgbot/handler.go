package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) getGame(update tgbotapi.Update) {
	ctx := b.NewContextWithoutPayload(update)

	if err := b.handler.GetGameFromMessage(ctx); err != nil {
		b.log.WithError(err).Error("get game from message")
	}
}
