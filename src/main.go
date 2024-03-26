package main

import "os"

func main() {
	token := os.Getenv("TELEGRAM_APITOKEN")
	bot := TgBot{token: token}
	bot.Start()
}
