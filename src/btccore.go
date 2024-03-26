package main

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"os"
)

func GetBalance(wallet string) (string, error) {
	log.Printf("Getting balance for wallet: %s", wallet)
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:                 "localhost:18332/wallet/" + wallet,
		User:                 os.Getenv("U"),
		Pass:                 os.Getenv("P"),
		HTTPPostMode:         true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:           true, // Bitcoin core does not provide TLS by default
		DisableAutoReconnect: false,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
		return "0", err
	}
	info, err := client.GetWalletInfo()
	if err != nil {
		if jerr, ok := err.(*btcjson.RPCError); ok {
			switch jerr.Code {
			case btcjson.ErrRPCWalletNotFound:
				log.Printf("Wallet not found, creating wallet: %v", wallet)
				_, err := client.CreateWallet(wallet)
				if err != nil {
					log.Printf("Error creating wallet: %v", err)
					return "0", err
				}
			default:
				log.Printf("Error getting wallet info: %v", err)
				return "0", err
			}
		}
	}
	log.Printf("Wallet info: %v", info)

	balance, err := client.GetBalance("*")
	recover()
	if err != nil {
		log.Fatal(err)
		return "0", err
	}
	defer client.Shutdown()
	return balance.String(), nil
}
