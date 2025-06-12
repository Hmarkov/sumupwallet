package main

import "time"

/*
Structures:
User of the system - id, name and a map of wallets where each key is the wallet id.
Wallet - id, name for easier user access, balance, link to the owner, Transactions list tracking depostis, withdrawals.
Transaction for each financial operation - type (deposit/withdraw), Amount, Date (time the transaction occured).
*/
type User struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Wallets map[string]Wallet `json:"wallets"`
}

type Wallet struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Balance      int           `json:"balance"`
	UserID       string        `json:"userId"`
	Transactions []Transaction `json:"transactions"`
}
type Transaction struct {
	Type   string
	Amount int
	Date   time.Time
}
