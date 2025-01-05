package telegram

import (
	"errors"

	"github.com/nhassl3/article-saver-bot/pkg/client"
	"github.com/nhassl3/article-saver-bot/pkg/e"
	"github.com/nhassl3/article-saver-bot/pkg/entities"
	"github.com/nhassl3/article-saver-bot/pkg/events"
	"github.com/nhassl3/article-saver-bot/pkg/storage/files"
)

const ErrProcMessage = "can't process message"

var (
	ErrUnknownMetaType  = errors.New("unknown meta type")
	ErrUnknownEventType = errors.New("unknown event type")
)

type Proc struct {
	tg      *client.Client
	offset  int
	storage *files.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

func NewProc(tg *client.Client, storage *files.Storage) *Proc {
	return &Proc{tg: tg, storage: storage}
}

func (p *Proc) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].UpdateId + 1

	return res, nil
}

func (p *Proc) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap(ErrProcMessage, ErrUnknownEventType)
	}
}

func (p *Proc) processMessage(event events.Event) error {
	meta, err := getMeta(event)
	if err != nil {
		return e.Wrap(ErrProcMessage, err)
	}

	err = p.doCmd(event.Text, meta.Username, meta.ChatId)
	return e.WrapIfErr(ErrProcMessage, err)
}

func getMeta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta) // type assertion
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(update entities.Updates) events.Event {
	updateType := fetchType(update)

	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	// chatId username
	if updateType == events.Message {
		res.Meta = Meta{
			ChatId:   update.Message.Chat.ChatId,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchText(update entities.Updates) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update entities.Updates) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
