package main

import (
	tgClient "TelegramBot/clients/telegram"
	event_consumer "TelegramBot/consumer/event-consumer"
	"TelegramBot/events/telegram"
	"TelegramBot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	tgBotHost      = "api.telegram.org"
	sqlStoragePath = "data/sqlite/storage.db"
	batchSize      = 100
)

func main() {

	s, err := sqlite.New(sqlStoragePath)
	if err != nil {
		log.Fatal("can't connect to the storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage", err)
	}

	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), s)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"6628958791:AAGJq7QVHL4_vLkpwtJO5D8X4u8j0MLTQYM",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
