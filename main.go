package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

func main() {
	// This deletion/re-creation is only for testing purposes
	os.Remove("sqlite-database.db")
	log.Println("sqlite-database.db deleted")

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ = sql.Open("sqlite3", "./sqlite-database.db")
	defer sqliteDatabase.Close()
	createAddressesTable(sqliteDatabase)
	createTransactionsTable(sqliteDatabase)

	// Manually inserting addresses for testing
	insertAddress(sqliteDatabase, generateAddress("5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"), 1000)
	insertAddress(sqliteDatabase, generateAddress("0b14d501a594442a01c6859541bcb3e8164d183d32937b851835442f69d5c94e"), 100)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /addresses/{address}", getAddress)
	mux.HandleFunc("POST /transactions", createTransaction)

	log.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}

func createAddressesTable(db *sql.DB) {
	createAddressTableSQL := `CREATE TABLE addresses (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"address" TEXT,
		"balance" INTEGER
	  );`

	log.Println("Create addresses table...")
	statement, err := db.Prepare(createAddressTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("addresses table created")
}

func insertAddress(db *sql.DB, address string, balance int) {
	log.Println("insertAddress:" + address + ":" + strconv.Itoa(balance))
	insertAddressSQL := `INSERT INTO addresses(address, balance) VALUES (?, ?)`
	statement, err := db.Prepare(insertAddressSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(address, balance)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func queryAddress(db *sql.DB, address string) int {
	query := "SELECT id, address, balance FROM addresses WHERE address = ?"
	row := db.QueryRow(query, address)

	var id int
	var balance int

	err := row.Scan(&id, &address, &balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Fatal(err)
	}

	log.Println("queryBalance:"+strconv.Itoa(id)+":"+address+":", strconv.Itoa(balance))
	return balance
}

func updateAddress(db *sql.DB, address string, newBalance int) {
	log.Println("updateAddress:" + address + ":" + strconv.Itoa(newBalance))
	updateSQL := `UPDATE addresses SET balance = ? WHERE address = ?`
	statement, err := db.Prepare(updateSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(newBalance, address)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func createTransactionsTable(db *sql.DB) {
	createAddressTableSQL := `CREATE TABLE transactions (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"sender" TEXT,
		"amount" INTEGER,
		"recipient" TEXT,
		"time" INTEGER
	  );`

	log.Println("Create transactions table...")
	statement, err := db.Prepare(createAddressTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("transactions table created")
}

func insertTransaction(db *sql.DB, sender string, amount int, recipient string, time int) {
	log.Println("insertTransaction:" + sender + ":" + strconv.Itoa(amount) + ":" + recipient + ":" + strconv.Itoa(time))
	insertAddressSQL := `INSERT INTO transactions(sender, amount, recipient, time) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertAddressSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(sender, amount, recipient, time)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func generateAddress(pkey string) string {
	sum := sha256.Sum256([]byte(pkey))
	addressHex := hex.EncodeToString(sum[:])
	address := addressHex[:12]

	return address
}

func validateAddress(address string) bool {
	_, err := strconv.ParseUint(address, 16, 64)
	return err == nil
}

func getAddress(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")

	if validateAddress(address) {
		balance := queryAddress(sqliteDatabase, address)
		log.Println("getAddress:" + address + ":" + strconv.Itoa(balance))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(balance)))
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid"))
		return
	}
}

type TransactionRequest struct {
	Pkey    string `json:"pkey"`
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	senderAddress := generateAddress(req.Pkey)
	senderBalance := queryAddress(sqliteDatabase, senderAddress)

	if validateAddress(req.Address) {
		if senderBalance >= req.Amount {
			log.Println("createTransaction:" + senderAddress + ":" + req.Address + ":" + strconv.Itoa(req.Amount))
			updateAddress(sqliteDatabase, senderAddress, senderBalance-req.Amount)
			recipientBalance := queryAddress(sqliteDatabase, req.Address)

			if recipientBalance > 0 {
				updateAddress(sqliteDatabase, req.Address, recipientBalance+req.Amount)
			} else {
				insertAddress(sqliteDatabase, req.Address, req.Amount)
			}

			insertTransaction(sqliteDatabase, senderAddress, req.Amount, req.Address, int(time.Now().Unix()))
			w.WriteHeader(http.StatusCreated)
		} else {
			log.Println("Not enough funds!")
			http.Error(w, "Not enough funds", http.StatusBadRequest)
		}
	} else {
		log.Println("Recipient address is invalid!")
		http.Error(w, "Invalid recipient address", http.StatusBadRequest)
	}
}
