package csv2sqlite

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

func castTypes(strs []string) []interface{} {
	ifaces := make([]interface{}, len(strs))

	for i, str := range strs {
		toInt, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			ifaces[i] = toInt
			continue
		}
		
		toFloat, err := strconv.ParseFloat(str, 64)
		if err == nil {
			ifaces[i] = toFloat
			continue
		}

		ifaces[i] = str
	}
	return ifaces
}

func ValidTableScheme(schemeStr string) bool {
	match, err := regexp.MatchString(`^\(.+, *.+\)$`, schemeStr)
	if err != nil {
		return false
	}
	return match
}

// TODO: escape "(", ")", ",", and so on.
func GenTableScheme(strs []string) string {
	return fmt.Sprintf("(%s)", strings.Join(strs, ","))
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
	return "(" + qParams + ")"
}

func (db *DB) CreateTable(name string, scheme string) (Table, error) {
	if !ValidTableScheme(scheme) {
		return Table{}, fmt.Errorf("invalid scheme")
	}

	_, err := db.Exec(
		fmt.Sprintf(`CREATE TABLE %s %s`, name, scheme),
	)
	return Table{DB: db, Name: name, Scheme: scheme, QParams: GenTableQParams(scheme)}, err
}


func (table Table) InsertRows(rows [][]string) error {
	for _, row := range rows {
		ifs := castTypes(row)
		_, err := table.Exec(
			fmt.Sprintf(`INSERT INTO %s %s VALUES %s`, table.Name, table.Scheme, table.QParams),
			ifs...,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
