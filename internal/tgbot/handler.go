package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
)

const ErrMsg = "URL is invalid. Press /help for information"

func (b *TgBot) startCommand(update tgbotapi.Update) {
	newUser := model.User{
		TelegramID:        int(update.Message.From.ID),
		TelegramChatID:    int(update.Message.Chat.ID),
		TelegramFirstName: update.Message.From.FirstName,
		TelegramLastName:  update.Message.From.LastName,
		TelegramUserName:  update.Message.From.UserName,
	}

	b.service.CreateUser(newUser)

	text := "hello from psn_discounter_bot. I send /help for more information"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}

func (b *TgBot) subscribeToGame(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	url := update.Message.Text

	if url == "" {
		msg.Text = ErrMsg

		b.SendMessage(msg)
	}

	responseMsg, err := b.service.SubscribeToGame(int(update.Message.From.ID), url)
	if err != nil {
		responseMsg = err.Error()
	}

	msg.Text = responseMsg

	b.SendMessage(msg)
}
