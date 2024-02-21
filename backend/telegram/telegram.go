package telegram

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
	"website/backend/telegram/utils"
)

var bot *telebot.Bot

func InitTgBot() error {
	fmt.Println("bot is running...")

	utils.LoadEnv()
	tokenBot := os.Getenv("TELEGRAM_BOT_KEY")
	var err error

	bot, err = telebot.NewBot(telebot.Settings{
		Token:  tokenBot,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	bot.Start()
	return nil
}

func SendTelegramMessage(message string) error {
	var err error

	_, err = bot.Send(telebot.ChatID(602974315), message)
	if err != nil {
		return err
	}
	return nil
}
