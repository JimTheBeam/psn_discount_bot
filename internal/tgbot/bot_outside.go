package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *TgBot) SendText(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)

	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}

func (b *TgBot) SendMessage(msg tgbotapi.Chattable) {
	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}

func (b *TgBot) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(chatID, text)
	message.DisableWebPagePreview = true

	return message
}
