package main

type User struct {
	Userid int64
	Wallet Wallet
}

func NewUserWithBtcWallet(userid int64) User {
	return User{Userid: userid, Wallet: Wallet{Amount: 0, Currency: BITCOIN}}
}
