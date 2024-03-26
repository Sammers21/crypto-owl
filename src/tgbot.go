package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Send", "Send"),
		tgbotapi.NewInlineKeyboardButtonData("Receive", "Receive"),
	),
)

type TgBot struct {
	token string
	//map of user id to user
	users map[int64]User
	bot   *tgbotapi.BotAPI
}

func (t *TgBot) Start() {
	bot, err := tgbotapi.NewBotAPI(t.token)
	t.bot = bot
	t.users = make(map[int64]User)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			var msg tgbotapi.MessageConfig
			switch update.Message.Text {
			case "/wallet", "/start":
				user, present := t.users[update.Message.Chat.ID]
				if present {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, user.Wallet.WalletMessage())
					msg.ReplyMarkup = numericKeyboard
				} else {
					newUser := NewUserWithBtcWallet(update.Message.Chat.ID)
					t.users[newUser.Userid] = newUser
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, newUser.Wallet.WalletMessage())
					msg.ReplyMarkup = numericKeyboard
				}
			}
			msg.ParseMode = "MarkdownV2"
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data+" received")
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
