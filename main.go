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

func NewDB(dbName string) (*DB, error) {
	db, err := sql.Open(`sqlite3`, dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (*DB) CreateTable(table Table) {
	_, err = db.Exec(
		fmt.Sprintf(`CREATE TABLE %s %s`),
	)
	if err != nil {
		log.Fatal(err)
	}

}

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