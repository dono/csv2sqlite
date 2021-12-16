package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	cs "github.com/dono/csv2sqlite"
	"github.com/urfave/cli/v2"
)

// check duplicate table name
func checkDupTable(strs []string) bool {
	strMap := make(map[string]struct{})
	for _, str := range strs {
		strMap[str] = struct{}{}
	}
	return len(strMap) != len(strs)
}

func action(tableNames, csvPaths []string, dbName string) error {
	if len(tableNames) != len(csvPaths) {
		return fmt.Errorf("invalid options: only one table for one csv file")
	}

	if checkDupTable(tableNames) {
		return fmt.Errorf("invalid options: each table name must be unique")
	}

	db, err := cs.NewDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	for i := range tableNames {
		csvFile, err := os.Open(csvPaths[i])
		if err != nil {
			return err
		}
		defer csvFile.Close()

		rows, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			return err
		}

		scheme := cs.GenTableScheme(rows[0], rows[1])
		table, err := db.CreateTable(tableNames[i], scheme)
		if err != nil {
			return err
		}

		err = table.InsertRows(rows[1:])
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// no options mode
	if len(os.Args) == 2 {
		csvPath := os.Args[1]
		// name of the input csv file (excluding the ext)
		csvFileName := filepath.Base(csvPath[:len(csvPath)-len(filepath.Ext(csvPath))])

		tableNames := []string{csvFileName}
		csvPaths := []string{csvPath}
		dbName := csvFileName + ".db"

		err := action(tableNames, csvPaths, dbName)
		if err != nil {
			log.Fatal(err)
		} else {
			os.Exit(0)
		}
	}

	app := &cli.App{
		Name:  "csv2sqlite",
		Usage: "$ csv2sqlite -t hoge_tb -c ./hoge.csv -d ./dump.db",
		Flags: []cli.Flag {
			&cli.StringSliceFlag{
				Name: "table",
				Aliases: []string{"t"},
				Usage: "-t hoge_tb",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name: "csv",
				Aliases: []string{"c"},
				Usage: "-c ./hoge.csv",
				Required: true,
			},
			&cli.StringFlag{
				Name: "db",
				Aliases: []string{"d"},
				Usage: "-d ./dump.db",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			dbName := c.String("db")
			tableNames := c.StringSlice("table")
			csvPaths := c.StringSlice("csv")

			return action(tableNames, csvPaths, dbName)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}