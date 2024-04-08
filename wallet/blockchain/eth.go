package blockchain

import (
	"context"
	"crypto-owl/abis/FiatTokenV22"
	"crypto-owl/abis/SwapRouter02"
	"crypto-owl/abis/WETH9"
	"crypto-owl/abis/tether"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func GetUSDTBalance(wallet string, contractAddress string) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	contract, err := tether.NewTether(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	amount, _ := contract.BalanceOf(&bind.CallOpts{}, key.Address)
	decimals, _ := contract.Decimals(&bind.CallOpts{})
	balanceInEth := new(big.Float).SetInt(amount)
	balanceInEth = new(big.Float).Quo(balanceInEth, big.NewFloat(math.Pow10(int(decimals.Int64()))))
	return balanceInEth.String() + " USDT", nil
}
func GetWETHBalance(wallet string, contractAddress string) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	contract, err := WETH9.NewWETH9(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	amount, _ := contract.BalanceOf(&bind.CallOpts{}, key.Address)
	decimals, _ := contract.Decimals(&bind.CallOpts{})
	balanceInEth := new(big.Float).SetInt(amount)
	balanceInEth = new(big.Float).Quo(balanceInEth, big.NewFloat(math.Pow10(int(decimals))))
	return balanceInEth.String() + " WETH", nil
}

func GetUSDCBalance(wallet string, contractAddress string) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	contract, err := FiatTokenV22.NewFiatTokenV22(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	amount, _ := contract.BalanceOf(&bind.CallOpts{}, key.Address)
	decimals, _ := contract.Decimals(&bind.CallOpts{})
	balanceInEth := new(big.Float).SetInt(amount)
	balanceInEth = new(big.Float).Quo(balanceInEth, big.NewFloat(math.Pow10(int(decimals))))
	return balanceInEth.String() + " USDC", nil
}

type EthKey struct {
	privateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

func (e *EthKey) Recalculate() {
	publicKey := e.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: PublicKey is not of type *ecdsa.PublicKey")
	}
	e.PublicKey = publicKeyECDSA
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Printf("Setting Address to: %s", address.Hex())
	e.Address = address
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
	return key.Address.Hex()
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

type UniSwapParams struct {
	ContractAddress  common.Address
	TokenIn          common.Address
	TokenOut         common.Address
	Recipient        common.Address
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
}

func UniswapSwap(wallet string, params UniSwapParams) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	contract, err := SwapRouter02.NewSwapRouter02(params.ContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	singleParams := SwapRouter02.IV3SwapRouterExactInputSingleParams{
		TokenIn:           params.TokenIn,
		TokenOut:          params.TokenOut,
		Recipient:         params.Recipient,
		AmountIn:          params.AmountIn,
		Fee:               new(big.Int).SetUint64(10000),
		AmountOutMinimum:  params.AmountOutMinimum,
		SqrtPriceLimitX96: new(big.Int).SetUint64(0),
	}
	//weth, er := WETH9.NewWETH9(params.TokenIn, client)
	//if er != nil {
	//	log.Fatal(er)
	//}
	nonce, err := client.PendingNonceAt(context.Background(), key.Address)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	//x, err := weth.Approve(&bind.TransactOpts{
	//	From:     key.Address,
	//	Context:  context.Background(),
	//	GasPrice: gasPrice,
	//	Nonce:    new(big.Int).SetUint64(nonce),
	//	Signer: func(c common.Address, tx *types.Transaction) (*types.Transaction, error) {
	//		log.Printf("Signing approve transaction for Address: %s", key.Address.Hex())
	//		return types.SignTx(tx, types.NewEIP155Signer(chainID), key.privateKey)
	//	},
	//}, params.ContractAddress, params.AmountIn)
	//log.Printf("Approved tx: %s", x.Hash().Hex())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//nonce, err = client.PendingNonceAt(context.Background(), key.Address)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//gasPrice, err = client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	res, err := contract.ExactInputSingle(&bind.TransactOpts{
		From:     key.Address,
		Context:  context.Background(),
		GasPrice: gasPrice,
		Nonce:    new(big.Int).SetUint64(nonce),
		Signer: func(c common.Address, tx *types.Transaction) (*types.Transaction, error) {
			log.Printf("Signing uniswap transaction for Address: %s", key.Address.Hex())
			return types.SignTx(tx, types.NewEIP155Signer(chainID), key.privateKey)
		},
	}, singleParams)
	if err != nil {
		log.Printf("TX: %+v", res)
		log.Fatal("UNISWAP Transaction failed: ", err)
	}
	return res.Hash().Hex(), nil
}

func GetEthBalance(wallet string) (string, error) {
	key, err := KeyForWallet(wallet)
	client, err := client()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	log.Printf("Getting balance for Address: %s", key.Address.Hex())
	balance, err := client.BalanceAt(context.Background(), key.Address, nil)
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
	nonce, err := client.PendingNonceAt(context.Background(), key.Address)
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
