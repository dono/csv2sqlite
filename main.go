package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type Table struct {
	Name   string
	Scheme string
}

func New(dbName string) (*DB, error) {
	db, err := sql.Open(`sqlite3`, dbName)
	return &DB{DB: db}, err
}

func (db *DB) CreateTable(table Table) error {
	_, err := db.Exec(
		fmt.Sprintf(`CREATE TABLE %s`, table.Name),
	)
	return err
}

func main() {
	TABLE_NAME := "hoge_tb"
	SCHEME := "(col1, col2)"

	db, err := sql.Open(`sqlite3`, `./foo.db`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(
		fmt.Sprintf(`CREATE TABLE %s %s`, TABLE_NAME, SCHEME),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		fmt.Sprintf(`INSERT INTO %s (col1, col2) VALUES (?, ?)`, TABLE_NAME),
		"123",
		"text",
	)
}