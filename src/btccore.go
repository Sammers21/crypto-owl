package main

import (
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"os"
)

func GetBalance(wallet string) (float64, error) {
	log.Printf("Getting balance for wallet: %s", wallet)
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18332/wallet/" + wallet,
		User:         os.Getenv("U"),
		Pass:         os.Getenv("P"),
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer client.Shutdown()
	info, err := client.GetBalance("*")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return info.ToBTC(), nil
}
