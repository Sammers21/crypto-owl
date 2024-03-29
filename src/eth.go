package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
)

func client() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:8545/")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

type EthKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
}

func (e *EthKey) Recalculate() {
	publicKey := e.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	e.publicKey = publicKeyECDSA
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Printf("Setting address to: %s", address.Hex())
	e.address = address
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func KeyForWallet(wallet string) (EthKey, error) {
	exists, _ := Exists(wallet)
	var ethKey EthKey
	if !exists {
		log.Printf("Wallet file does not exist: %s, creating it", wallet)
		// handle the case where the file doesn't exist
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
		ethKey = EthKey{privateKey: privateKey}
		ethKey.Recalculate()
		_ = crypto.FromECDSA(privateKey)
		//err = os.WriteFile(wallet, b, 0644)
		//if err != nil {
		//	log.Fatal(err)
		//}
		return ethKey, nil
	} else {
		log.Printf("Wallet file exists: %s, reading it", wallet)
		fd, err := os.ReadFile(wallet)
		if err != nil {
			log.Fatal(err)
		}
		privateKey, err := crypto.HexToECDSA(string(fd))
		if err != nil {
			log.Fatal(err)
		}
		ethKey = EthKey{privateKey: privateKey}
		ethKey.Recalculate()
		return ethKey, nil
	}
}

func GetEthBalance(wallet string) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	log.Printf("Getting balance for address: %s", key.address.Hex())
	balance, err := client.BalanceAt(context.Background(), key.address, nil)
	if err != nil {
		return "0", err
	}
	return balance.String(), err
}
