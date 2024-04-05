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
	USDT_ERC20
	USDC_ERC20
	WETH_ERC20
)

func CurrencyFromString(s string) Currency {
	switch s {
	case "BTC":
		return BITCOIN
	case "ETH":
		return ETHEREUM
	case "TRX":
		return TRON
	case "USDT-ERC20":
		return USDT_ERC20
	case "USDC-ERC20":
		return USDC_ERC20
	case "WETH-ERC20":
		return WETH_ERC20
	}
	return BITCOIN
}

func (c Currency) String() string {
	switch c {
	case BITCOIN:
		return "BTC"
	case ETHEREUM:
		return "ETH"
	case TRON:
		return "TRX"
	case USDT_ERC20:
		return "USDT"
	case USDC_ERC20:
		return "USDC"
	case WETH_ERC20:
		return "WETH"
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
	case USDT_ERC20:
		return "USDT ERC20"
	case USDC_ERC20:
		return "USDC ERC20"
	case WETH_ERC20:
		return "WETH ERC20"
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
	} else if w.Currency == USDT_ERC20 {
		balance, err := blockchain.GetUSDTBalance(w.Id, "0x7169D38820dfd117C3FA1f22a697dBA58d90BA06")
		if err != nil {
			return "Error getting balance"
		}
		return balance
	} else if w.Currency == USDC_ERC20 {
		balance, err := blockchain.GetUSDCBalance(w.Id, "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238")
		if err != nil {
			return "Error getting balance"
		}
		return balance
	} else if w.Currency == WETH_ERC20 {
		balance, err := blockchain.GetWETHBalance(w.Id, "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14")
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
	} else if w.Currency == ETHEREUM || w.Currency == USDT_ERC20 || w.Currency == USDC_ERC20 || w.Currency == WETH_ERC20 {
		address = blockchain.GetEthAddress(w.Id)
	} else {
		return "Not implemented"
	}
	return fmt.Sprintf("*Receive %s*\n\nUse the address below to send *%s* to the CryptoOwl bot wallet address\\.\n\n"+
		"*Address:* `%s`\n\n", w.Currency.FullName(), w.Currency.FullName(), address)
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
