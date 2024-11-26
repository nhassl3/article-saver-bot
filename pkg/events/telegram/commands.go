package telegram

import (
	"errors"
	"fmt"
	"github.com/nhassl3/article-saver-bot/pkg/client"
	"github.com/nhassl3/article-saver-bot/pkg/e"
	"github.com/nhassl3/article-saver-bot/pkg/storage"
	"github.com/nhassl3/article-saver-bot/pkg/storage/files"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Proc) doCmd(cmd, username string, chatId int) error {
	cmd = strings.TrimSpace(cmd)
	sendMsg := NewMessageSender(chatId, p.tg)

	log.Printf("got new command '%s' from '%s'", cmd, username)

	// add page: http://...
	// rnd page: /rnd
	// help menu: /help
	// start: /start: Welcome message and help menu
	if isURL(cmd) {
		return p.SavePage(cmd, username, chatId)
	}
	switch cmd {
	case RndCmd:
		if err := p.SendRandom(chatId, username); err != nil {
			log.Printf("failed to send random command: %s", err)
		}
	case HelpCmd:
		if err := sendMsg(msgHelp); err != nil {
			log.Printf("failed to send help menu: %s", err)
		}
	case StartCmd:
		if err := p.SendRandom(chatId, username); err != nil {
			log.Printf("failed to send start menu: %s", err)
		}
	default:
		if err := sendMsg(msgUnknownCommand); err != nil {
			log.Printf("failed to send command: %s", err)
			return err
		}
		log.Printf("unknown command: %s", cmd)
		return fmt.Errorf("unknown command: '%s'", cmd)
	}
	return nil
}

func (p *Proc) SavePage(pageURL, username string, chatId int) (err error) {
	defer func() { err = e.WrapIfErr("can't do cmd save page", err) }()

	sendMsg := NewMessageSender(chatId, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return e.Wrap("can't check exists of the link", err)
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err = p.storage.Save(page); err != nil {
		return err
	}

	return sendMsg(msgSaved)
}

func (p *Proc) SendRandom(chatId int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do send random", err) }()

	sendMsg := NewMessageSender(chatId, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, files.ErrNoSavedPage) {
		return err
	}

	if errors.Is(err, files.ErrNoSavedPage) {
		return sendMsg(msgNoSavedPages)
	}

	if err = sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func isURL(cmd string) bool {
	u, err := url.Parse(cmd)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func NewMessageSender(chatId int, tg *client.Client) func(string) error {
	return func(msg string) error {
		return e.WrapIfErr("can't send message", tg.SendMessage(chatId, msg))
	}
}
