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
	Amount   int64
	Currency Currency
}

func (w Wallet) WalletMessage() string {
	return fmt.Sprintf("💰My Wallet \n\n*%s*: %d", w.Currency.String(), w.Amount)
}
