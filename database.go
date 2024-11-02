package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

func setupDatabase() {
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

	createAddressesTable(sqliteDatabase)
	createTransactionsTable(sqliteDatabase)
	createBlocksTable(sqliteDatabase)
}

func createAddressesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE addresses (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"address" TEXT,
		"balance" INTEGER
	  );`

	log.Println("Create addresses table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("addresses table created")
}

func createTransactionsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE transactions (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"sender" TEXT,
		"amount" INTEGER,
		"recipient" TEXT,
		"time" INTEGER
	  );`

	log.Println("Create transactions table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("transactions table created")
}

func createBlocksTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE blocks (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"block" TEXT,
		"prevBlock" TEXT,
		"address" TEXT,
		"nonce" TEXT,
		"time" INTEGER
	  );`

	log.Println("Create blocks table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("blocks table created")

	insertSQL := `INSERT INTO blocks(block, prevBlock, address, nonce, time) VALUES (?, ?, ?, ?, ?)` // genesis block
	statement, err = db.Prepare(insertSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec("0", "0", "address", "nonce", int(time.Now().Unix()))
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Created genesis block")
}

func insertAddress(db *sql.DB, address string, balance int) {
	insertSQL := `INSERT INTO addresses(address, balance) VALUES (?, ?)`
	statement, err := db.Prepare(insertSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(address, balance)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func queryAddress(db *sql.DB, address string) int {
	querySQL := "SELECT id, address, balance FROM addresses WHERE address = ?"
	row := db.QueryRow(querySQL, address)

	var id int
	var balance int

	err := row.Scan(&id, &address, &balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Fatal(err)
	}

	return balance
}

func updateAddress(db *sql.DB, address string, newBalance int) {
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

func insertTransaction(db *sql.DB, sender string, amount int, recipient string, time int) {
	insertSQL := `INSERT INTO transactions(sender, amount, recipient, time) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(sender, amount, recipient, time)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func queryBlock(db *sql.DB) (string, error) {
	querySQL := "SELECT block FROM blocks ORDER BY id DESC LIMIT 1"
	row := db.QueryRow(querySQL)

	var block string
	err := row.Scan(&block)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return block, nil
}

func insertBlock(db *sql.DB, block string, prevBlock string, address string, nonce string, time int) {
	insertSQL := `INSERT INTO blocks(block, prevBlock, address, nonce, time) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(block, prevBlock, address, nonce, time)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
