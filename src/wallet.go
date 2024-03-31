package main

import (
	"fmt"
)

type Currency int

const (
	BITCOIN Currency = iota
	ETHEREUM
	TRON
	USDT
)

func (c Currency) String() string {
	switch c {
	case BITCOIN:
		return "BTC"
	case ETHEREUM:
		return "ETH"
	case TRON:
		return "TRX"
	case USDT:
		return "USDT"
	}
	return "UNKNOWN"
}

type Wallet struct {
	Id       string
	Currency Currency
}

func (w Wallet) Balance() string {
	if w.Currency == BITCOIN {
		balance, err := GetBtcBalance(w.Id)
		if err != nil {
			return "Error getting balance"
		}
		return balance
	} else if w.Currency == ETHEREUM {
		balance, err := GetEthBalance(w.Id)
		if err != nil {
			return "Error getting balance"
		}
		return balance
	}
	return "Not implemented"
}

func (w Wallet) Receive() string {
	var address string
	if w.Currency == BITCOIN {
		address = GetBtcAddress(w.Id)
	} else if w.Currency == ETHEREUM {
		address = GetEthAddress(w.Id)
	} else {
		return "Not implemented"
	}
	return fmt.Sprintf("*Receive*\n\nUse the address below to send BTC to the CryptoOwl bot wallet address\\.\nNetwork: *Bitcoin \\- BTC*\\.\n\n*Address:* `%s`\n\n Funds will be credited within 30\\-60 minutes\\.", address)
}

func (w Wallet) Send(amount int64, address string) string {
	txid, err := SendBtc(w.Id, address, amount)
	if err != nil {
		return "Error: `" + err.Error() + "`"
	}
	return fmt.Sprintf("Transaction ID: %s", txid)
}
