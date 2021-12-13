package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	cs "github.com/dono/csv2sqlite"
	"github.com/urfave/cli/v2"
)

func isExistDup(strs []string) bool {
	strMap := make(map[string]struct{})
	for _, str := range strs {
		strMap[str] = struct{}{}
	}
	return len(strMap) != len(strs)
}

func main() {
	app := &cli.App{
		Name:  "csv2sqlite",
		Usage: "$ csv2sqlite -d ./hoge.db -t fuga_tb -c ./fuga.csv",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "db, d",
				Value: "dump.db",
				Usage: "-d ./hoge.db",
			},
			&cli.StringSliceFlag{
				Name:  "table, t",
				Usage: "-t fuga_tb",
			},
			&cli.StringSliceFlag{
				Name:  "csv, c",
				Usage: "-c ./fuga.csv",
			},
		},
		Action: func(c *cli.Context) error {
			dbName := c.String("db")
			tableNames := c.StringSlice("table")
			csvPaths := c.StringSlice("csv")

			if len(tableNames) != len(csvPaths) {
				return fmt.Errorf("invalid options: only one table for one csv file")
			}

			if isExistDup(tableNames) {
				return fmt.Errorf("invalid options: each table name must be unique")
			}

			db, err := cs.NewDB(dbName)
			if err != nil {
				return err
			}
			defer db.Close()


			// loop
			csvFile, err := os.Open("")
			if err != nil {
				return err
			}
			defer csvFile.Close()

			reader := csv.NewReader(csvFile)
			rows, err := reader.ReadAll()
			if err != nil {
				return err
			}
			scheme := cs.GenTableScheme(rows[0])
			table, err := db.CreateTable("test_tb", scheme)
			if err != nil {
				return err
			}

			err = table.InsertRows(rows[1:])
			if err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}