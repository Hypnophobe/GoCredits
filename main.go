package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	loadDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /addresses/{address}", getAddress)
	mux.HandleFunc("POST /transactions", createTransaction)
	mux.HandleFunc("GET /blocks", getBlock)
	mux.HandleFunc("POST /blocks", submitBlock)

	log.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}
