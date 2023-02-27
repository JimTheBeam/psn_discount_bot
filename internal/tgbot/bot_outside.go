package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) SendText(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)

	_, err := b.Bot.Send(msg)

	return err
}

func (b *Bot) SendMessage(msg tgbotapi.Chattable) error {
	_, err := b.Bot.Send(msg)

	return err
}

func (b *Bot) SendTextInReplyToMessage(chatID int64, messageID int, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID

	_, err := b.Bot.Send(msg)

	return err
}

func (b *Bot) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(chatID, text)
	message.DisableWebPagePreview = true

	return message
}
