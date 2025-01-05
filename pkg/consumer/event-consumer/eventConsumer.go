package event_consumer

import (
	"log"
	"time"

	"github.com/nhassl3/article-saver-bot/pkg/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{fetcher: fetcher, processor: processor, batchSize: batchSize}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err = c.HandleEvents(gotEvents); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}
	}
}

func (c *Consumer) HandleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		// Механизм ReTry, механизм backup'а, фоллбэк
		if err := c.processor.Process(event); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}
	}
	return nil
}
