package wallet

import (
	"crypto-owl/wallet/blockchain"
	"fmt"
	"math/big"
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

func (c Currency) FullName() string {
	switch c {
	case BITCOIN:
		return "Bitcoin"
	case ETHEREUM:
		return "Ethereum"
	case TRON:
		return "Tron"
	case USDT:
		return "USDT ERC20"
	}
	return "UNKNOWN"
}

type Wallet struct {
	Id       string
	Currency Currency
}

func (w Wallet) Balance() string {
	if w.Currency == BITCOIN {
		balance, err := blockchain.GetBtcBalance(w.Id)
		if err != nil {
			return "Error getting balance"
		}
		return balance
	} else if w.Currency == ETHEREUM {
		balance, err := blockchain.GetEthBalance(w.Id)
		if err != nil {
			return "Error getting balance"
		}
		return balance
	} else if w.Currency == USDT {
		balance, err := blockchain.GetUSDTBalance(w.Id, "0x7169D38820dfd117C3FA1f22a697dBA58d90BA06")
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
		address = blockchain.GetBtcAddress(w.Id)
	} else if w.Currency == ETHEREUM {
		address = blockchain.GetEthAddress(w.Id)
	} else {
		return "Not implemented"
	}
	return fmt.Sprintf("*Receive %s*\n\nUse the address below to send %s to the CryptoOwl bot wallet address\\.\n\n"+
		"*Address:* `%s`\n\n", w.Currency.String(), w.Currency.String(), address)
}

func (w Wallet) Send(amount big.Int, address string) string {
	if w.Currency == BITCOIN {
		txid, err := blockchain.SendBtc(w.Id, address, amount)
		if err != nil {
			return "Error: `" + err.Error() + "`"
		}
		return fmt.Sprintf("BTC Transaction ID: `%s`", txid)
	} else if w.Currency == ETHEREUM {
		txid, err := blockchain.SendEth(w.Id, address, amount)
		if err != nil {
			return "Error: `" + err.Error() + "`"
		}
		return fmt.Sprintf("ETH Transaction ID: `%s`", txid)
	}
	return "Not implemented"
}
