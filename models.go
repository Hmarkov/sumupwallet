package main

type User struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Wallets map[string]Wallet `json:"wallets"`
}

type Wallet struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	UserID  string `json:"userId"`
}
