package main

import (
	"fmt"
	"strings"
)

type User struct {
	Userid  int64
	Wallets map[Currency]Wallet
}

func NewUser(userid int64, platform Platform) User {
	return User{Userid: userid, Wallets: map[Currency]Wallet{
		BITCOIN:  {Id: "tg-" + fmt.Sprint(userid) + "-btc", Currency: BITCOIN},
		ETHEREUM: {Id: "tg-" + fmt.Sprint(userid) + "-eth", Currency: ETHEREUM},
		//TRON:     {Id: "tg-" + fmt.Sprint(userid) + "-trx", Currency: TRON},
		//USDT:     {Id: "tg-" + fmt.Sprint(userid) + "-usdt", Currency: USDT},
	},
	}
}

func (u User) WalletMessage() string {
	var balances string
	for _, wallet := range u.Wallets {
		balances += fmt.Sprintf("*%s*: %s\n", wallet.Currency.String(), wallet.Balance())
	}
	return strings.Replace(
		fmt.Sprintf("💰My Wallet \n\n%s", balances),
		".", "\\.", -1,
	)
}
