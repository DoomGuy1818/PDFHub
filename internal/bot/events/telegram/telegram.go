package telegram

import (
	"PDFHub/internal/bot/clients"
	"PDFHub/internal/bot/clients/telegram"
	"PDFHub/internal/bot/events"
	"PDFHub/internal/bot/lib/e"
	"errors"
)

type Processor struct {
	tg     *telegram.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(tg *telegram.Client) *Processor {
	return &Processor{
		tg: tg,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't fetch updates", err)
	}

	result := make([]events.Event, 0, len(updates))

	if len(updates) == 0 {
		return nil, nil
	}

	for _, update := range updates {
		result = append(result, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func event(upd clients.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.ProcessMessage(event)

	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) ProcessMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't meta event", ErrUnknownMetaType)
	}
	return res, nil
}

func fetchText(upd clients.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd clients.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
