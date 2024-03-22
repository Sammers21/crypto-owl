package main

type Currency int

const (
	BITCOIN Currency = iota
	ETHEREUM
	TRON
	USDT
)

type Wallet struct {
	Amount   int64
	Currency Currency
}
