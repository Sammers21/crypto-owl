package bot

import (
	"crypto-owl/wallet"
	"fmt"
	"strings"
)

type User struct {
	Userid  int64
	Wallets map[wallet.Currency]wallet.Wallet
}

func NewUser(userid int64, platform Platform) User {
	prefix := platform.Prefix()
	return User{Userid: userid, Wallets: map[wallet.Currency]wallet.Wallet{
		wallet.BITCOIN:  {Id: prefix + fmt.Sprint(userid) + "-btc", Currency: wallet.BITCOIN},
		wallet.ETHEREUM: {Id: prefix + fmt.Sprint(userid) + "-eth", Currency: wallet.ETHEREUM},
		wallet.USDT:     {Id: prefix + fmt.Sprint(userid) + "-usdt", Currency: wallet.USDT},
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
