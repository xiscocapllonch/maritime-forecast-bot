package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	result, err := GetForecast(os.Getenv("XML_ID"))

	if err != nil {
		log.Panic(err)
	}

	config := tgbotapi.NewMessageToChannel(os.Getenv("CHANNEL_USER"), result)
	config.ParseMode = "html"

	_, err = bot.Send(config)

	if err != nil {
		log.Panic(err)
	}
}
