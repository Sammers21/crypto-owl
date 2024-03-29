package main

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Send BTC", "SendBtc"),
		tgbotapi.NewInlineKeyboardButtonData("Receive BTC", "ReceiveBtc"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Send ETH", "SendEth"),
		tgbotapi.NewInlineKeyboardButtonData("Receive ETH", "ReceiveEth"),
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
				log.Printf("User %d requested wallet", update.Message.Chat.ID)
				user, present := t.users[update.Message.Chat.ID]
				if present {
					log.Printf("User %d already has wallet", update.Message.Chat.ID)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, user.WalletMessage())
					msg.ReplyMarkup = numericKeyboard
				} else {
					log.Printf("User %d does not have wallet, creating one", update.Message.Chat.ID)
					newUser := NewUserWithBtcWallet(update.Message.Chat.ID)
					t.users[newUser.Userid] = newUser
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, newUser.WalletMessage())
					msg.ReplyMarkup = numericKeyboard
				}
			default:
				log.Printf("User %d sent message: %s", update.Message.Chat.ID, update.Message.Text)
				split := strings.Split(update.Message.Text, " ")
				if len(split) != 3 {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command, please use `/send <amount> <address>`")
				} else {
					amount, err := strconv.ParseFloat(split[1], 64)
					amounti64 := int64(amount * 100000000)
					log.Println("Amount: ", amounti64)
					if err != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid amount")
					} else {
						address := split[2]
						user, present := t.users[update.Message.Chat.ID]
						if !present {
							log.Printf("User %d does not have wallet, creating one", update.Message.Chat.ID)
							newUser := NewUserWithBtcWallet(update.Message.Chat.ID)
							t.users[newUser.Userid] = newUser
							user = newUser
						}
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, user.Wallets[BITCOIN].Send(amounti64, address))
					}
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
			switch update.CallbackQuery.Data {
			case "Receive":
				user, present := t.users[update.CallbackQuery.Message.Chat.ID]
				if !present {
					log.Printf("User %d does not have wallet, creating one", update.CallbackQuery.Message.Chat.ID)
					newUser := NewUserWithBtcWallet(update.CallbackQuery.Message.Chat.ID)
					t.users[newUser.Userid] = newUser
					user = newUser
				}
				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, user.Wallets[BITCOIN].Receive())
				msg.ParseMode = "MarkdownV2"
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "SendBtc":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "In order to send, just type `/send <amount> <address>`")
				msg.ParseMode = "MarkdownV2"
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Unknown command")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}
}
