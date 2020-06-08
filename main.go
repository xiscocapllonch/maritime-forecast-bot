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

	result, err := getXML("http://www.aemet.es/xml/maritima/" + os.Getenv("XML_ID") + ".xml")

	if err != nil {
		log.Panic(err)
	}

	text := result.formatText()

	config := tgbotapi.NewMessageToChannel(os.Getenv("CHANNEL_USER"), text)
	config.ParseMode = "html"

	_, err = bot.Send(config)

	if err != nil {
		log.Panic(err)
	}
}
