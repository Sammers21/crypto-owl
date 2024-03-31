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
	prefix := platform.Prefix()
	return User{Userid: userid, Wallets: map[Currency]Wallet{
		BITCOIN:  {Id: prefix + fmt.Sprint(userid) + "-btc", Currency: BITCOIN},
		ETHEREUM: {Id: prefix + fmt.Sprint(userid) + "-eth", Currency: ETHEREUM},
	},
	}
}

func (u User) WalletMessage() string {
	var balances string
	for _, wallet := range u.Wallets {
		balances += fmt.Sprintf("*%s*: %s\n", wallet.Currency.FullName(), wallet.Balance())
	}
	return strings.Replace(
		fmt.Sprintf("💰My Wallet \n\n%s", balances),
		".", "\\.", -1,
	)
}
