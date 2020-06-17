package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/xiscocapllonch/aemet"
	"image/gif"
	"log"
	"os"
)

type Specification struct {
	Token     string `required:"true"`
	ChannelId int64  `split_words:"true" required:"true"`
}

func sendBotMessage(token string, config tgbotapi.Chattable) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	msg, err := bot.Send(config)

	fmt.Printf("%+v", msg.MessageID)

	if err != nil {
		return err
	}

	return nil
}

func sendHTMLMessage(html, token string, channelId int64) error {
	config := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: channelId,
			DisableNotification: true,
		},
		Text: html,
		ParseMode: "html",
	}

	err := sendBotMessage(token, config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var s Specification

	err := envconfig.Process("bot", &s)
	if err != nil {
		log.Fatal(err)
	}

	var customMsg string

	var cmdCustomMsg = &cobra.Command{
		Use:   "custom [Post custom message]",
		Short: "Post custom message on telegram channel",
		Run: func(cmd *cobra.Command, args []string) {
			err = sendHTMLMessage(customMsg, s.Token, s.ChannelId)
			if err != nil {
				log.Fatal(err)
			}

		},
	}

	cmdCustomMsg.Flags().StringVar(&customMsg, "customMsg", "", "custom message html for public on telegram channel")
	err = cmdCustomMsg.MarkFlagRequired("customMsg")

	var xmlId string

	var cmdForecast = &cobra.Command{
		Use:   "forecast [Post forecast]",
		Short: "Post maritime forecast on telegram channel",
		Run: func(cmd *cobra.Command, args []string) {
			forecastHTML, err := aemet.GetMaritimeForecast(xmlId)
			if err != nil {
				log.Fatal(err)
			}
			err = sendHTMLMessage(forecastHTML, s.Token, s.ChannelId)
			if err != nil {
				log.Fatal(err)
			}

		},
	}

	cmdForecast.Flags().StringVar(&xmlId, "xmlId", "", "xmlId config for github.com/xiscocapllonch/aemet")
	err = cmdForecast.MarkFlagRequired("xmlId")
	if err != nil {
		log.Fatal(err)
	}

	var zoneId string
	var wind bool

	var cmdForecastMap = &cobra.Command{
		Use:   "forecastMap [Post forecast map]",
		Short: "Post maritime forecast map on telegram channel",
		Run: func(cmd *cobra.Command, args []string) {
			forecastGif, err := aemet.GetMaritimeForecastMapGIF(zoneId, 24, 150, wind)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create("forecast.gif")
			if err != nil {
				log.Fatal(err)
				return
			}

			err = gif.EncodeAll(f, &forecastGif)
			if err != nil {
				log.Fatal(err)
			}

			label := "Mar combinada"

			if wind {
				label = "Mar de viento"
			}

			config := tgbotapi.AnimationConfig{
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID: s.ChannelId,
						DisableNotification: true,
					},
					File:     f.Name(),
					MimeType: "image/gif",
				},
				Caption: fmt.Sprintf(
					"%s [©AEMET](http://aemet.es/es/eltiempo/prediccion/maritima)",
					label,
				),
				ParseMode: "markdown",
			}

			err = sendBotMessage(s.Token, config)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmdForecastMap.Flags().StringVar(&zoneId, "zoneId", "", "zoneId config for github.com/xiscocapllonch/aemet")
	err = cmdForecastMap.MarkFlagRequired("zoneId")
	if err != nil {
		log.Fatal(err)
	}

	cmdForecastMap.Flags().BoolVar(&wind, "wind", false, "wind config for github.com/xiscocapllonch/aemet")
	err = cmdForecastMap.MarkFlagRequired("wind")
	if err != nil {
		log.Fatal(err)
	}

	var rootCmd = &cobra.Command{Use: "bot"}
	rootCmd.AddCommand(cmdCustomMsg, cmdForecast, cmdForecastMap)

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
