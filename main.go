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
		fmt.Sprintf(`CREATE TABLE %s %s`, table.Name, table.Scheme),
	)
	return err
}

func (db *DB) 

func main() {
	db, err := sql.Open(`sqlite3`, `./foo.db`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		fmt.Sprintf(`CREATE TABLE % %`),
	)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(

	)

}