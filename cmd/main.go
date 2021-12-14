package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

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

func main() {
	app := &cli.App{
		Name:  "csv2sqlite",
		Usage: "$ csv2sqlite -t hoge_tb -c ./hoge.csv -d ./dump.db",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "db",
				Aliases: []string{"d"},
				Usage: "-d ./dump.db",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "table",
				Aliases: []string{"t"},
				Usage: "-t hoge_tb",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "csv",
				Aliases: []string{"f"},
				Usage: "-c ./hoge.csv",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			dbName := c.String("db")
			tableNames := c.StringSlice("table")
			csvPaths := c.StringSlice("csv")

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

				scheme := cs.GenTableScheme(rows[0])
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}