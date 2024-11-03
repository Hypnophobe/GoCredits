package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	loadDatabase()
	defer sqliteDatabase.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /address/{address}", getAddress)                  // Get a single address
	mux.HandleFunc("GET /addresses", getAddresses)                        // Get all addresses
	mux.HandleFunc("POST /transaction", createTransaction)                // Create a transaction
	mux.HandleFunc("GET /transaction/{id}", getTransaction)               // Get single transaction by ID
	mux.HandleFunc("GET /transactions/{address}", getAddressTransactions) // Get all transactions relating to an address
	mux.HandleFunc("GET /transactions", getTransactions)                  // Get all transactions from database
	mux.HandleFunc("POST /block", submitBlock)                            // Submit a block
	mux.HandleFunc("GET /block", getBlock)                                // Get last block
	mux.HandleFunc("GET /blocks", getBlocks)                              // Get all blocks

	log.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}
