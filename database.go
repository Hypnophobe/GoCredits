package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

type Address struct {
	ID      int
	Address string
	Balance int
}

type Transaction struct {
	ID        int
	Sender    string
	Amount    int
	Recipient string
	Time      string
}

type Block struct {
	ID           int    `json:"id"`
	BlockContent string `json:"block"`
	PrevBlock    string `json:"prevBlock"`
	Address      string `json:"address"`
	Nonce        string `json:"nonce"`
	Time         int    `json:"time"`
}

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

func queryAddresses(db *sql.DB) ([]Address, error) {
	querySQL := "SELECT id, address, balance FROM addresses"
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address

	for rows.Next() {
		var addr Address
		if err := rows.Scan(&addr.ID, &addr.Address, &addr.Balance); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
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

func queryTransaction(db *sql.DB, id string) (*Transaction, error) {
	querySQL := "SELECT id, sender, amount, recipient, time FROM transactions WHERE id = ?"
	row := db.QueryRow(querySQL, id)

	var txn Transaction
	err := row.Scan(&txn.ID, &txn.Sender, &txn.Amount, &txn.Recipient, &txn.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Fatal(err)
		return nil, err
	}

	return &txn, nil
}

func queryTransactions(db *sql.DB) ([]Transaction, error) {
	querySQL := "SELECT id, sender, amount, recipient, time FROM transactions"
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var addr Transaction
		if err := rows.Scan(&addr.ID, &addr.Sender, &addr.Amount, &addr.Recipient, &addr.Time); err != nil {
			return nil, err
		}
		transactions = append(transactions, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func queryAddressTransactions(db *sql.DB, address string) ([]Transaction, error) {
	querySQL := "SELECT id, sender, amount, recipient, time FROM transactions WHERE sender = ? OR recipient = ?"
	rows, err := db.Query(querySQL, address, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var addr Transaction
		if err := rows.Scan(&addr.ID, &addr.Sender, &addr.Amount, &addr.Recipient, &addr.Time); err != nil {
			return nil, err
		}
		transactions = append(transactions, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
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

func queryBlocks(db *sql.DB) ([]Block, error) {
	querySQL := "SELECT id, block, prevBlock, address, nonce, time FROM blocks"
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []Block

	for rows.Next() {
		var blk Block
		if err := rows.Scan(&blk.ID, &blk.BlockContent, &blk.PrevBlock, &blk.Address, &blk.Nonce, &blk.Time); err != nil {
			return nil, err
		}
		blocks = append(blocks, blk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
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

func getSupply(db *sql.DB) (int, error) {
	querySQL := "SELECT SUM(balance) FROM addresses"
	var totalBalance int

	err := db.QueryRow(querySQL).Scan(&totalBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return totalBalance, nil
}
