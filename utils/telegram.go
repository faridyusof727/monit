package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"mon-tool-be/models"
	"os"
	"strconv"
	"strings"
)

func SendMessage(alerts []models.Alert, message string) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	for _, alert := range alerts {
		telegramID, err := strconv.ParseInt(alert.Key, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		msg := tgbotapi.NewMessage(telegramID, message)
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func InitTelegram(db *gorm.DB) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Hi "+update.Message.Chat.UserName+",\n"+
					"\n"+
					"Thank you for using "+os.Getenv("APP_NAME")+" Telegram Integration. To start integrating, please ENTER your integration key. You can get the key in your Monitor Settings.",
			)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages

			var alert models.Alert
			var monitor models.Monitor

			splitted := strings.Split(update.Message.Text, "___")

			r := db.Where("owner = ?", splitted[0]).Where("id = ?", splitted[1]).First(&monitor)

			if r.RowsAffected > 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You key is correct. You will receive alert as a message. You also may invite our bot to a group.")
				msg.ReplyToMessageID = update.Message.MessageID
				alert.MonitorID = monitor.ID
				alert.Key = strconv.Itoa(int(update.Message.Chat.ID))
				db.Create(&alert)
				if _, err := bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			}

			if r.RowsAffected == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You key is incorrect. Please try again. If you have any problem, you can email us a support ticket.")
				if _, err := bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			}
		}

	}
}
