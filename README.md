# csv2sqlite
CLI tool to convert csv to sqlite with type info

## Usage

Run with options (single table)

```
$ csv2sqlite -c ./foo.csv -t foo_tb -d ./dump.db
foo.csv -> dump.db (incl. foo_tb)
```

Run with options (multiple tables)

```
$ csv2sqlite -c ./foo.csv -t foo_tb -c ./bar.csv -t bar_tb -d ./dump.db
(foo.csv, bar.csv) -> dump.db (incl. foo_tb, bar_tb)
```

Run with no options (only single table)

```
$ csv2sqlite ./foo.csv
foo.csv -> foo.db (incl. foo)
```
Useful in combination with drag-and-drop or shortcuts

## Installation

```
$ go get github.com/dono/csv2sqlite
```

## License

MIT