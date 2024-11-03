package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

func main() {
	os.Remove("sqlite-database.db")
	log.Println("sqlite-database.db deleted")

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sqliteDatabase.Close() // Ensure the database is closed when done

	createAddressesTable(sqliteDatabase)
	createTransactionsTable(sqliteDatabase)
	createBlocksTable(sqliteDatabase)

	log.Println("Done")
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
