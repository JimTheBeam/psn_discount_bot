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

	msg.ReplyMarkup = model.NewSubscriptionListKeyboard()

	if _, err := b.Bot.Send(msg); err != nil {
		b.log.WithError(err).Error("send message")
	}
}

func (b *TgBot) getGame(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	url := update.Message.Text

	if url == "" {
		msg.Text = ErrMsg

		b.SendMessage(msg)
	}

	game, subscription, err := b.service.GetGame(int(update.Message.From.ID), url)
	if err != nil {
		msg.Text = err.Error()
		b.SendMessage(msg)
	}

	if subscription == nil {
		msg.Text = game.Name + "\nChoose price you want to subscribe"
		msg.ReplyMarkup = model.NewGameKeyboardWithPrices(game.Prices)
	} else {
		msg.Text = game.Name + "\n" + game.GetPriceText() + "\n"
		msg.ReplyMarkup = model.NewGameUnsubscribeKeyboard(game.ID)
	}

	b.SendMessage(msg)
}
