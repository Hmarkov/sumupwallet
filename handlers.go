package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var users = make(map[string]User)
var wallets = make(map[string]Wallet)

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userid"]
	user, ok := users[userID]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if len(user.Wallets) >= 10 {
		http.Error(w, "wallet limit reached (max 10 wallets per user)", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "wallet name is required", http.StatusBadRequest)
		return
	}

	if user.Wallets == nil {
		user.Wallets = make(map[string]Wallet)
	}

	for _, wallet := range user.Wallets {
		if wallet.Name == req.Name {
			http.Error(w, "wallet name already exists", http.StatusBadRequest)
			return
		}
	}

	wallet := Wallet{
		ID:      uuid.New().String(),
		Name:    req.Name,
		Balance: 0,
		UserID:  userID,
	}

	user.Wallets[wallet.ID] = wallet
	users[userID] = user
	wallets[wallet.ID] = wallet

	json.NewEncoder(w).Encode(wallet)
}

func GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := mux.Vars(r)["walletid"]
	wallet, ok := wallets[walletID]
	if !ok {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func GetWallets(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userid"]
	user, ok := users[userID]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	walletList := []Wallet{}
	for _, wallet := range user.Wallets {
		walletList = append(walletList, wallet)
	}

	json.NewEncoder(w).Encode(walletList)
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userid"]
	walletName := mux.Vars(r)["walletname"]
	amountStr := mux.Vars(r)["amount"]

	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	user, ok := users[userID]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var wallet *Wallet
	for _, w := range user.Wallets {
		if w.Name == walletName {
			wallet = &w
			break
		}
	}

	if wallet == nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	wallet.Balance += amount
	user.Wallets[wallet.ID] = *wallet
	users[userID] = user
	wallets[wallet.ID] = *wallet

	json.NewEncoder(w).Encode(wallet)
}

func Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userid"]
	walletName := mux.Vars(r)["walletname"]
	amountStr := mux.Vars(r)["amount"]

	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	user, ok := users[userID]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var walletID string
	for id, wallet := range user.Wallets {
		if wallet.Name == walletName {
			walletID = id
			break
		}
	}
	if walletID == "" {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	wallet := user.Wallets[walletID]
	if wallet.Balance < amount {
		http.Error(w, "insufficient funds", http.StatusBadRequest)
		return
	}

	wallet.Balance -= amount
	user.Wallets[walletID] = wallet
	users[userID] = user
	wallets[walletID] = wallet

	json.NewEncoder(w).Encode(wallet)
}
