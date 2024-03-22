package main

type User struct {
	Userid  int64
	Wallets []Wallet
}

func NewUserWithBtcWallet(userid int64) User {
	return User{Userid: userid, Wallets: []Wallet{{Amount: 0, Currency: BITCOIN}}}
}
