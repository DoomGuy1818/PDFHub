package bot

import (
	tgClient "PDFHub/internal/bot/clients/telegram"
	eventConsumer "PDFHub/internal/bot/consumer/event-consumer"
	"PDFHub/internal/bot/events/telegram"
	"log"
	"os"
)

func main() {
	eventsProcessor := telegram.New(tgClient.New(mustHost(), mustToken()))

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustHost() string {
	return os.Getenv("BOT_HOST")
}

func mustToken() string {
	return os.Getenv("BOT_TOKEN")
}
