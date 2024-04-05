package main

import (
	"crypto-owl/bot"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_APITOKEN")
	b := bot.TgBot{Token: token}
	b.Start()
}
