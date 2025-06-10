package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setup() {
	users = make(map[string]User)
	wallets = make(map[string]Wallet)
	users["u1"] = User{ID: "u1", Name: "Test User"}
}

func TestCreateWallet(t *testing.T) {
	setup()

	body := []byte(`{"name":"TestWallet"}`)
	req := httptest.NewRequest("POST", "/users/u1/wallets", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"userid": "u1"})
	rr := httptest.NewRecorder()

	CreateWallet(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var wallet Wallet
	err := json.NewDecoder(rr.Body).Decode(&wallet)
	assert.NoError(t, err)
	assert.Equal(t, "TestWallet", wallet.Name)
	assert.Equal(t, 0, wallet.Balance)
	assert.Equal(t, "u1", wallet.UserID)
}

func TestCreateDuplicateWalletName(t *testing.T) {
	setup()

	body := []byte(`{"name":"Duplicate"}`)
	req1 := httptest.NewRequest("POST", "/users/u1/wallets", bytes.NewBuffer(body))
	req1 = mux.SetURLVars(req1, map[string]string{"userid": "u1"})
	CreateWallet(httptest.NewRecorder(), req1)

	req2 := httptest.NewRequest("POST", "/users/u1/wallets", bytes.NewBuffer(body))
	req2 = mux.SetURLVars(req2, map[string]string{"userid": "u1"})
	rr := httptest.NewRecorder()
	CreateWallet(rr, req2)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWalletLimit(t *testing.T) {
	setup()

	for i := 0; i < 10; i++ {
		body := []byte(`{"name":"W` + strconv.Itoa(i) + `"}`)
		req := httptest.NewRequest("POST", "/users/u1/wallets", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"userid": "u1"})
		CreateWallet(httptest.NewRecorder(), req)
	}

	body := []byte(`{"name":"ExtraWallet"}`)
	req := httptest.NewRequest("POST", "/users/u1/wallets", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"userid": "u1"})
	rr := httptest.NewRecorder()
	CreateWallet(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func createTestWallet(t *testing.T) Wallet {
	setup()

	body := []byte(`{"name":"TestWallet"}`)
	req := httptest.NewRequest("POST", "/users/u1/wallets", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"userid": "u1"})
	rr := httptest.NewRecorder()
	CreateWallet(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var wallet Wallet
	err := json.NewDecoder(rr.Body).Decode(&wallet)
	assert.NoError(t, err)
	return wallet
}

func TestDeposit(t *testing.T) {
	wallet := createTestWallet(t)

	depReq := httptest.NewRequest("PUT", "/wallets/"+wallet.ID+"/deposit/200", nil)
	depReq = mux.SetURLVars(depReq, map[string]string{"walletid": wallet.ID, "amount": "200"})
	depRes := httptest.NewRecorder()
	Deposit(depRes, depReq)

	assert.Equal(t, http.StatusOK, depRes.Code)

	var updated Wallet
	err := json.NewDecoder(depRes.Body).Decode(&updated)
	assert.NoError(t, err)
	assert.Equal(t, 200, updated.Balance)
}

func TestWithdraw(t *testing.T) {
	wallet := createTestWallet(t)

	// deposit first so we can withdraw
	depReq := httptest.NewRequest("PUT", "/wallets/"+wallet.ID+"/deposit/200", nil)
	depReq = mux.SetURLVars(depReq, map[string]string{"walletid": wallet.ID, "amount": "200"})
	Deposit(httptest.NewRecorder(), depReq)

	withReq := httptest.NewRequest("PUT", "/wallets/"+wallet.ID+"/withdraw/100", nil)
	withReq = mux.SetURLVars(withReq, map[string]string{"walletid": wallet.ID, "amount": "100"})
	withRes := httptest.NewRecorder()
	Withdraw(withRes, withReq)

	assert.Equal(t, http.StatusOK, withRes.Code)

	var afterWithdraw Wallet
	err := json.NewDecoder(withRes.Body).Decode(&afterWithdraw)
	assert.NoError(t, err)
	assert.Equal(t, 100, afterWithdraw.Balance)
}

func TestOverWithdraw(t *testing.T) {
	wallet := createTestWallet(t)

	// deposit 50
	depReq := httptest.NewRequest("PUT", "/wallets/"+wallet.ID+"/deposit/50", nil)
	depReq = mux.SetURLVars(depReq, map[string]string{"walletid": wallet.ID, "amount": "50"})
	Deposit(httptest.NewRecorder(), depReq)

	// try to withdraw 100 fail
	owReq := httptest.NewRequest("PUT", "/wallets/"+wallet.ID+"/withdraw/100", nil)
	owReq = mux.SetURLVars(owReq, map[string]string{"walletid": wallet.ID, "amount": "100"})
	owRes := httptest.NewRecorder()
	Withdraw(owRes, owReq)

	assert.Equal(t, http.StatusBadRequest, owRes.Code)
}
