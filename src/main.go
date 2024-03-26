package main

import "os"

func main() {
	println(GetBalance("test01"))
	token := os.Getenv("TELEGRAM_APITOKEN")
	bot := TgBot{token: token}
	bot.Start()
}
