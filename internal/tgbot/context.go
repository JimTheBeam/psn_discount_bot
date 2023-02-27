package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Context interface {
		Bot() *Bot
		Update() tgbotapi.Update
		Payload() string
	}

	botContext struct {
		bot     *Bot
		update  tgbotapi.Update
		payload string
	}
)

func (b *Bot) NewContextWithPayload(update tgbotapi.Update, payload string) Context {
	return &botContext{
		bot:     b,
		update:  update,
		payload: payload,
	}
}

func (b *Bot) NewContextWithoutPayload(update tgbotapi.Update) Context {
	return &botContext{
		bot:    b,
		update: update,
	}
}

func (c *botContext) Bot() *Bot {
	return c.bot
}

func (c *botContext) Update() tgbotapi.Update {
	return c.update
}

func (c *botContext) Payload() string {
	return c.payload
}
