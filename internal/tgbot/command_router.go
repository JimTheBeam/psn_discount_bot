package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *TgBot) CommandRouter(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		b.startCommand(update)

		return

	case "help":
		msg.Text = `Send me url from psn store and I will tell you when the price reduces`
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}
