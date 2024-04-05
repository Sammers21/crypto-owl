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
		wallet.BITCOIN:    {Id: prefix + fmt.Sprint(userid) + "-btc", Currency: wallet.BITCOIN},
		wallet.ETHEREUM:   {Id: prefix + fmt.Sprint(userid) + "-eth", Currency: wallet.ETHEREUM},
		wallet.USDT_ERC20: {Id: prefix + fmt.Sprint(userid) + "-eth", Currency: wallet.USDT_ERC20},
		wallet.USDC_ERC20: {Id: prefix + fmt.Sprint(userid) + "-eth", Currency: wallet.USDC_ERC20},
		wallet.WETH_ERC20: {Id: prefix + fmt.Sprint(userid) + "-eth", Currency: wallet.WETH_ERC20},
	},
	}
}

func (u User) WalletMessage() string {
	var balances string
	balances += fmt.Sprintf("*%s*: %s\n", wallet.BITCOIN.FullName(), u.Wallets[wallet.BITCOIN].Balance())
	balances += fmt.Sprintf("*%s*: %s\n", wallet.ETHEREUM.FullName(), u.Wallets[wallet.ETHEREUM].Balance())
	balances += fmt.Sprintf("*%s*: %s\n", wallet.USDT_ERC20.FullName(), u.Wallets[wallet.USDT_ERC20].Balance())
	balances += fmt.Sprintf("*%s*: %s\n", wallet.USDC_ERC20.FullName(), u.Wallets[wallet.USDC_ERC20].Balance())
	balances += fmt.Sprintf("*%s*: %s\n", wallet.WETH_ERC20.FullName(), u.Wallets[wallet.WETH_ERC20].Balance())
	return strings.Replace(
		fmt.Sprintf("💰My Wallet \n\n%s", balances),
		".", "\\.", -1,
	)
}
