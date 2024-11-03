package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type getAddressResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

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

func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	if validateAddress(address) {
		balance := queryAddress(sqliteDatabase, address)

		response := getAddressResponse{
			Address: address,
			Balance: balance,
		}

		writeJSONResponse(w, http.StatusOK, response)
	} else {
		response := map[string]interface{}{"ok": false, "error": "invalid address"}
		writeJSONResponse(w, http.StatusBadRequest, response)
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := map[string]interface{}{"ok": false, "error": "invalid request body"}
		writeJSONResponse(w, http.StatusBadRequest, response)
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

			response := map[string]interface{}{"ok": true}
			writeJSONResponse(w, http.StatusOK, response)
		} else {
			response := map[string]interface{}{"ok": false, "error": "insufficient funds"}
			writeJSONResponse(w, http.StatusBadRequest, response)
		}
	} else {
		response := map[string]interface{}{"ok": false, "error": "invalid address"}
		writeJSONResponse(w, http.StatusBadRequest, response)
	}
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	block, err := queryBlock(sqliteDatabase)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"block": block})
}

func submitBlock(w http.ResponseWriter, r *http.Request) {
	var req submittedBlock

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := map[string]interface{}{"ok": false, "error": "invalid request body"}
		writeJSONResponse(w, http.StatusBadRequest, response)
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
			response := map[string]interface{}{"ok": true}
			writeJSONResponse(w, http.StatusOK, response)
			return
		} else {
			response := map[string]interface{}{"ok": false, "error": "invalid block"}
			writeJSONResponse(w, http.StatusBadRequest, response)
		}
	} else {
		response := map[string]interface{}{"ok": false, "error": "previous block mismatch"}
		writeJSONResponse(w, http.StatusBadRequest, response)
	}
}
