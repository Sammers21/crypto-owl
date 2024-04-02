package bot

import (
	"crypto-owl/wallet"
	"log"
	"math/big"
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
	Token string
	//map of user id to user
	users map[int64]User
	bot   *tgbotapi.BotAPI
}

func (t *TgBot) Start() {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}
	t.bot = bot
	t.users = make(map[int64]User)
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
					newUser := NewUser(update.Message.Chat.ID, TELEGRAM)
					t.users[newUser.Userid] = newUser
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, newUser.WalletMessage())
					msg.ReplyMarkup = numericKeyboard
				}
			default:
				log.Printf("User %d sent message: %s", update.Message.Chat.ID, update.Message.Text)
				split := strings.Split(update.Message.Text, " ")
				if len(split) != 4 {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command, please use `/send <currency> <amount> <address>`")
				} else {
					currency := split[1]
					amount := split[2]
					address := split[3]
					amountfloat, err := strconv.ParseFloat(amount, 64)
					if err != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid amount")
					} else if currency == "BTC" {
						amountBigFlt := new(big.Float).Mul(big.NewFloat(amountfloat), big.NewFloat(100000000))
						amounti64 := new(big.Int)
						amountBigFlt.Int(amounti64)
						log.Println("Amount: ", amounti64)
						user, present := t.users[update.Message.Chat.ID]
						if !present {
							log.Printf("User %d does not have wallet, creating one", update.Message.Chat.ID)
							newUser := NewUser(update.Message.Chat.ID, TELEGRAM)
							t.users[newUser.Userid] = newUser
							user = newUser
						}
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, user.Wallets[wallet.BITCOIN].Send(*amounti64, address))
					} else if currency == "ETH" {
						amountBigFlt := new(big.Float).Mul(big.NewFloat(amountfloat), big.NewFloat(1000000000000000000))
						amounti64 := new(big.Int)
						amountBigFlt.Int(amounti64)
						log.Println("Amount: ", amounti64)
						user, present := t.users[update.Message.Chat.ID]
						if !present {
							log.Printf("User %d does not have wallet, creating one", update.Message.Chat.ID)
							newUser := NewUser(update.Message.Chat.ID, TELEGRAM)
							t.users[newUser.Userid] = newUser
							user = newUser
						}
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, user.Wallets[wallet.ETHEREUM].Send(*amounti64, address))
					} else {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid currency")
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
			case "ReceiveBtc":
				user, present := t.users[update.CallbackQuery.Message.Chat.ID]
				if !present {
					log.Printf("User %d does not have wallet, creating one", update.CallbackQuery.Message.Chat.ID)
					newUser := NewUser(update.CallbackQuery.Message.Chat.ID, TELEGRAM)
					t.users[newUser.Userid] = newUser
					user = newUser
				}
				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, user.Wallets[wallet.BITCOIN].Receive())
				msg.ParseMode = "MarkdownV2"
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "SendBtc":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					"In order to send, just type `/send BTC <amount> <address>`")
				msg.ParseMode = "MarkdownV2"
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "SendEth":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					"In order to send, just type `/send ETH <amount> <address>`")
				msg.ParseMode = "MarkdownV2"
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "ReceiveEth":
				user, present := t.users[update.CallbackQuery.Message.Chat.ID]
				if !present {
					log.Printf("User %d does not have wallet, creating one", update.CallbackQuery.Message.Chat.ID)
					newUser := NewUser(update.CallbackQuery.Message.Chat.ID, TELEGRAM)
					t.users[newUser.Userid] = newUser
					user = newUser
				}
				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, user.Wallets[wallet.ETHEREUM].Receive())
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
