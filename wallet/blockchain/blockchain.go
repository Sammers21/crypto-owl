package blockchain

import "math/big"

type Blockchain interface {
	// GetBalance returns the balance of the given address
	GetBalance(address string) (string, error)
	// Send sends the given amount to the given address
	Send(wallet, to string, amount big.Int) (string, error)
	// Address returns the address of the given wallet
	Address(wallet string) (string, error)
}
