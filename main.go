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
	router.HandleFunc("/wallets/{walletid}", GetWallet).Methods("GET")
	router.HandleFunc("/users/{userid}/wallets", CreateWallet).Methods("POST")
	router.HandleFunc("/wallets/{walletid}/deposit/{amount}", Deposit).Methods("PUT")
	router.HandleFunc("/wallets/{walletid}/withdraw/{amount}", Withdraw).Methods("PUT")

	log.Println("running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
