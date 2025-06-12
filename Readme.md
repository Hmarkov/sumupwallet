# Wallet Management API

A simple RESTful API that allows users to manage wallets, perform deposits and withdrawals, and view transaction history. It supports multiple wallets per user (up to 10), and each wallet maintains its own balance and transaction log.


## Run API
<pre> go run . </pre>


## Run Tests
<pre>go test -v</pre>

## Example Requests
<h3>Create a wallet</h3><br>
<pre>curl -X POST http://localhost:8080/users/u1/wallets -H "Content-Type: application/json" -d '{"name":"MyWallet"}' </pre>
<h3>Deposit to a wallet</h3><br>
<pre>curl -X PUT http://localhost:8080/users/u1/wallets/MyWallet/deposit/100</pre>
<h3>Withdraw from a wallet</h3><br>
<pre>curl -X PUT http://localhost:8080/users/u1/wallets/MyWallet/withdraw/50</pre>
<h3>Get wallet info</h3><br>
<pre>curl http://localhost:8080/users/u1/wallets/MyWallet</pre>
<h3>Get all wallets</h3><br>
<pre>curl http://localhost:8080/users/u1/wallets</pre>
<h3>Get wallet transactions</h3><br>
<pre>curl http://localhost:8080/users/u1/wallets/MyWallet/transactions</pre>

## Features
 
- Create a new wallet for a user (max 10)
- List all wallets<
- Get details of specific wallet
- Withdraw funds from a wallet
- View all transactions for a wallet


## Technologies Used

- **Gorilla Mux** – Router for managing dynamic HTTP routes
- **UUID** – To generate unique wallet IDs
- **net/http** – Standard library for HTTP server functionality
- **Testing** – `net/http/httptest` and `stretchr/testify/assert` for unit testing

## Notes

- Wallet names must be unique
- Each user can have up to 10 wallets
- Withdrawals fail if the wallet has insufficient balance
- The Application uses in memory data, no persistance storage logic implemented