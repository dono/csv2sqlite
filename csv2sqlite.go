package csv2sqlite

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

type Table struct {
	*DB
	Name   	string
	Scheme 	string
	QParams string
}

func NewDB(dbName string) (*DB, error) {
	db, err := sql.Open(`sqlite`, dbName)
	return &DB{DB: db}, err
}

func castType(str string) interface{} {
	toInt, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return toInt
	}
	
	toFloat, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return toFloat
	}

	return str
}

func castTypes(strs []string) []interface{} {
	ifaces := make([]interface{}, len(strs))

	for i, str := range strs {
		ifaces[i] = castType(str)
	}
	return ifaces
}

// TODO: escape "(", ")", ",", and so on.
func GenTableScheme(header []string, sampleRow []string) string {
	cols := []string{}

	for i := range header {
		rt := reflect.TypeOf(castType(sampleRow[i]))
		switch rt.Kind() {
		case reflect.Int64:
			cols = append(cols, fmt.Sprintf("%s %s", header[i], "INTEGER"))
		case reflect.Float64:
			cols = append(cols, fmt.Sprintf("%s %s", header[i], "REAL"))
		case reflect.String:
			cols = append(cols, fmt.Sprintf("%s %s", header[i], "TEXT"))
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(cols, ","))
}

func GenTableQParams(scheme string) string {
	qParams := ""
	schemeArr := strings.Split(scheme, ",")
	for i := 0; i < len(schemeArr); i++ {
		if i == len(schemeArr) - 1 {
			qParams += "?"
		} else {
			qParams += "?,"
		}
	}
	return fmt.Sprintf("(%s)", qParams)
}

func (db *DB) CreateTable(name string, scheme string) (Table, error) {
	_, err := db.Exec(
		fmt.Sprintf(`CREATE TABLE %s %s`, name, scheme),
	)
	return Table{DB: db, Name: name, Scheme: scheme, QParams: GenTableQParams(scheme)}, errors.WithMessage(err, "SQLite: CREATE TABLE error")
}

func (table Table) InsertRows(rows [][]string) error {
	for _, row := range rows {
		_, err := table.Exec(
			fmt.Sprintf(`INSERT INTO %s VALUES %s`, table.Name, table.QParams),
			castTypes(row)...,
		)
		if err != nil {
			return errors.WithMessage(err, "SQLite: INSERT error")
		}
	}
	return nil
}
