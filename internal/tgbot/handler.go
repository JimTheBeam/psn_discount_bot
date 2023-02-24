package tgbot

import (
	"fmt"
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

		return
	}

	game, subscription, err := b.service.GetGame(int(update.Message.From.ID), url)
	if err != nil {
		msg.Text = err.Error()
		b.SendMessage(msg)

		return
	}

	msg.Text, msg.ReplyMarkup = getTextAndMarkupForGame(game, subscription)

	b.SendMessage(msg)
}

func (b *TgBot) getSubscriptions(chatID, userID int64) {
	data := payload.Subscriptions{
		UserID: int(userID),
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

func getTextAndMarkupForGame(game model.Game, subscription *model.UsersGames) (string, tgbotapi.InlineKeyboardMarkup) {
	var (
		text        string
		replyMarkup tgbotapi.InlineKeyboardMarkup
	)

	if subscription == nil {
		text = game.Name + "\nChoose price you want to subscribe"
		replyMarkup = model.NewGameKeyboardWithPrices(game.Prices)
	} else {
		text = fmt.Sprintf("%s\nCurrent price:\n%s\n\nSubscription price: %.2f",
			game.Name, game.GetPriceText(), subscription.SubscriptionPrice,
		)

		replyMarkup = model.NewGameUnsubscribeKeyboard(game.ID)
	}

	return text, replyMarkup
}
