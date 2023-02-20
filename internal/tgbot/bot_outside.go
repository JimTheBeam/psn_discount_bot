package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *TgBot) SendText(chatID int, message string) {
	msg := tgbotapi.NewMessage(int64(chatID), message)

	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}

func (b *TgBot) SendMessage(msg tgbotapi.Chattable) {
	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}
