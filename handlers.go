package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TransactionRequest struct {
	Pkey    string `json:"pkey"`
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type submittedBlock struct {
	Block         string `json:"block"`
	PreviousBlock string `json:"prevBlock"`
	Address       string `json:"address"`
	Nonce         string `json:"nonce"`
}

func getAddress(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")

	if validateAddress(address) {
		balance := queryAddress(sqliteDatabase, address)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(balance)))
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid"))
		return
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	senderAddress := generateAddress(req.Pkey)
	senderBalance := queryAddress(sqliteDatabase, senderAddress)

	if validateAddress(req.Address) {
		if senderBalance >= req.Amount {
			updateAddress(sqliteDatabase, senderAddress, senderBalance-req.Amount)
			recipientBalance := queryAddress(sqliteDatabase, req.Address)

			if recipientBalance > 0 {
				updateAddress(sqliteDatabase, req.Address, recipientBalance+req.Amount)
			} else {
				insertAddress(sqliteDatabase, req.Address, req.Amount)
			}

			insertTransaction(sqliteDatabase, senderAddress, req.Amount, req.Address, int(time.Now().Unix()))
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "insufficient", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "invalid", http.StatusBadRequest)
	}
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	block, err := queryBlock(sqliteDatabase)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(block))
}

func submitBlock(w http.ResponseWriter, r *http.Request) {
	var req submittedBlock

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	qBlock, err := queryBlock(sqliteDatabase)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if qBlock == req.PreviousBlock {
		if req.Block == genBlock(req.PreviousBlock, req.Address, req.Nonce) {
			insertBlock(sqliteDatabase, req.Block, req.PreviousBlock, req.Address, req.Nonce, int(time.Now().Unix()))
			oldBalance := queryAddress(sqliteDatabase, req.Address)

			if oldBalance > 0 {
				updateAddress(sqliteDatabase, req.Address, oldBalance+1)
			} else {
				insertAddress(sqliteDatabase, req.Address, oldBalance+1)
			}

			insertTransaction(sqliteDatabase, "null", 1, req.Address, int(time.Now().Unix()))
			return
		} else {
			http.Error(w, "invalid", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "invalid", http.StatusBadRequest)
	}
}
