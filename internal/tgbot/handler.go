package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/model/payload"
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

func (b *TgBot) getSubscriptions(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	data := payload.Subscriptions{
		UserID: int(update.Message.From.ID),
		Limit:  0,
		Offset: 0,
	}

	subs, err := b.service.GetSubscriptions(data)
	if err != nil {
		b.SendText(chatID, err.Error())
		return
	}

	if len(subs) == 0 {
		b.SendText(chatID, "You haven't got subscriptions yet.")
		return
	}

	text := "Subscriptions:"

	msg := tgbotapi.NewMessage(chatID, text)

	msg.ReplyMarkup = model.NewSubscriptionsListKeyboard(subs)

	b.SendMessage(msg)
}
