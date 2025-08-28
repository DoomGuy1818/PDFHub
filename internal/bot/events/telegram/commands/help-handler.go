package commands

type HelpHandler struct {
	ChatID int
}

func NewHelpHandler(chatID int) HelpHandler {
	return HelpHandler{
		ChatID: chatID,
	}
}

func (h HelpHandler) Handle() error {
	return nil
}
