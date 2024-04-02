package main

import (
	"crypto-owl/bot"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_APITOKEN")
	bot := bot.TgBot{Token: token}
	bot.Start()
}
