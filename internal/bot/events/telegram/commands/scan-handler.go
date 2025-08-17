package commands

type ScanHandler struct {
	ChatID int
}

func NewScanHandler(chatID int) ScanHandler {
	return ScanHandler{
		ChatID: chatID,
	}
}

func (h ScanHandler) Handle() error {
	return nil
}
