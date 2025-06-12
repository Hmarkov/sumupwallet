package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	users["u1"] = User{
		ID: "u1", Name: "Test User", Wallets: make(map[string]Wallet),
	}
	router := mux.NewRouter()

	router.HandleFunc("/users/{userid}/wallets", GetWallets).Methods("GET")
	router.HandleFunc("/users/{userid}/wallets/{walletname}", GetWallet).Methods("GET")
	router.HandleFunc("/users/{userid}/wallets", CreateWallet).Methods("POST")
	router.HandleFunc("/users/{userid}/wallets/{walletname}/deposit/{amount}", Deposit).Methods("PUT")
	router.HandleFunc("/users/{userid}/wallets/{walletname}/withdraw/{amount}", Withdraw).Methods("PUT")
	router.HandleFunc("/users/{userid}/wallets/{walletname}/transactions", GetTransactions).Methods("GET")

	log.Println("running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
