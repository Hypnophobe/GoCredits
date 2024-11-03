package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

func loadDatabase() {
	var err error
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
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
