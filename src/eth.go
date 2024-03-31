package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	"os"
)

func client() (*ethclient.Client, error) {
	client, err := ethclient.Dial("https://public.stackup.sh/api/v1/node/ethereum-sepolia")
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

func GetEthAddress(wallet string) string {
	key, err := KeyForWallet(wallet)
	if err != nil {
		log.Fatal(err)
	}
	return key.address.Hex()
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
		err = crypto.SaveECDSA(wallet, privateKey)
		if err != nil {
			log.Fatal(err)
		}
		return ethKey, nil
	} else {
		log.Printf("Wallet file exists: %s, reading it", wallet)
		privateKey, err := crypto.LoadECDSA(wallet)
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
	// convert wei to eth
	balanceInEth := new(big.Float).SetInt(balance)
	balanceInEth = new(big.Float).Quo(balanceInEth, big.NewFloat(math.Pow10(18)))
	return balanceInEth.String() + " ETH", nil
}

func SendEth(wallet string, toAddress string, amount big.Int) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	nonce, err := client.PendingNonceAt(context.Background(), key.address)
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddressEth := common.HexToAddress(toAddress)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddressEth,
		Value:    &amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	})
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key.privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	return signedTx.Hash().Hex(), nil
}
