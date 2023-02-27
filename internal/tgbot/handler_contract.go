package tgbot

type Handler interface {
	SubscribeToGameCallback(c Context) error
	UnsubscribeToGameCallback(c Context) error
	GetSubscriptionsCallback(c Context) error
	GetGameCallback(c Context) error
	UnknownCallback(c Context) error
	CancelCallback(c Context) error

	StartCommand(c Context) error
	GetSubscriptionsCommand(c Context) error
	HelpCommand(c Context) error
	UnknownCommand(c Context) error

	GetGameFromMessage(c Context) error
}
