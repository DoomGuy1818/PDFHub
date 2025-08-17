package telegram

type CommandHandler interface {
	Handle() error
}
