package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) CommandRouter(update tgbotapi.Update) {
	command := update.Message.Command()

	ctx := b.NewContextWithoutPayload(update)

	handlerFunc := b.getCommandHandlerFunc(command)

	if err := handlerFunc(ctx); err != nil {
		b.log.WithError(err).WithField("command", command).Error("handle command")
	}

	return
}

func (b *Bot) getCommandHandlerFunc(command string) HandlerFunc {
	switch command {
	case "start":
		return b.handler.StartCommand
	case "subscription_list":
		return b.handler.GetSubscriptionsCommand
	case "help":
		return b.handler.HelpCommand
	default:
		return b.handler.UnknownCommand
	}
}
