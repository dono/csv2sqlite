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
	*DB
	Name   string
	Scheme string
}

func NewDB(dbName string) (*DB, error) {
	db, err := sql.Open(`sqlite3`, dbName)
	return &DB{DB: db}, err
}

func (db *DB) CreateTable(name string, scheme string) (Table, error) {
	_, err := db.Exec(
		fmt.Sprintf(`CREATE TABLE %s %s`, name, scheme),
	)
	return Table{DB: db, Name: name, Scheme: scheme}, err
}

func castTypes(strs []string) []interface{} {
	ifaces := make([]interface{}, len(strs))

	for i, str := range strs {
		// check type
		switch
		//ifaces[i] = str
	}
	return ifaces
}

func (table Table) InsertRows(rows [][]string) error {
	for _, row := range rows {
		ifs := make([]interface{}, len(row))
		_, err := table.Exec(
			fmt.Sprintf(`INSERT INTO %s %s VALUES (?, ?)`, table.Name, table.Scheme),
			ifs...,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db, err := NewDB(`./foo.db`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	table, err := db.CreateTable("test_tb", "(col1, col2)")
	if err != nil {
		log.Fatal(err)
	}

	err = table.InsertRows([][]string{
		{"foo1", "bar1"},
		{"foo2", "bar2"},
		{"foo3", "bar3"},
	})
	if err != nil {
		log.Fatal(err)
	}
}