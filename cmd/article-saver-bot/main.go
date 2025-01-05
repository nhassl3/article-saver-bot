package main

import (
	"log"

	"github.com/nhassl3/article-saver-bot/pkg/client"
	"github.com/nhassl3/article-saver-bot/pkg/config"
	eventConsumer "github.com/nhassl3/article-saver-bot/pkg/consumer/event-consumer"
	"github.com/nhassl3/article-saver-bot/pkg/events/telegram"
	"github.com/nhassl3/article-saver-bot/pkg/storage/files"
)

func main() {
	cfg := config.Config{}
	if err := cfg.MustLoad(); err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}

	eventsProcessor := telegram.NewProc(
		client.NewClient(cfg.Host, cfg.Token),
		files.NewStorage("storage"),
	)

	log.Println("[INFO] Starting server")

	consumer := eventConsumer.NewConsumer(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}
}
