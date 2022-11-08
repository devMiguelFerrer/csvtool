# CSV TOOL
Tools to csv file in go language

## Installation
1. You can use the below Go command to install Gin.
```go
go get -u github.com/devMiguelFerrer/csvtool
```
2. Import it in your code:
```go
import "github.com/devMiguelFerrer/csvtool"
```

## Features
```go
// function to create new csv tool
func NewCSVTool(fileName string, delimiter rune) CSVTool
```
```go
// features available
type CSVTool interface {
	ToSQL(tableName, outFileName string, rowsByFile int) error
	RemoveColumn(columns ...string) CSVTool
	SplitCSV(rowsByFile int) error
}
```
## Examples & Usage

```c
// test.csv
ONE;TWO;THREE;FOUR;FIVE;SIX
11;1'2;1'3;1'4;15;16
2'1;2'2;2'3;24;2'5;2'6
31;3'2;3'3;34;3'5;36
```
```go
// initialize your new csv tool
c := NewCSVTool("test.csv", ';')
```

```go
// example to split csv
c.SplitCSV(1)
//expected result 2 csv files

// test_01.csv
// ONE;TWO;THREE;FOUR;FIVE;SIX
// 11;1'2;1'3;1'4;15;16
// 2'1;2'2;2'3;24;2'5;2'6

// test_02.csv
// ONE;TWO;THREE;FOUR;FIVE;SIX
// 31;3'2;3'3;34;3'5;36
```

```go
// examples to create sql files without some columns
c.RemoveColumn("TWO", "THREE").ToSQL("TABLE_SQL", "migration", 2)
//expected result 2 sql files without some columns

// migration_01.sql
// INSERT INTO TABLE_SQL (ONE,FOUR,FIVE,SIX) VALUES
// ('11','1''4','15','16'),
// ('2''1','24','2''5','2''6');

// migration_02.sql
// INSERT INTO TABLE_SQL (ONE,FOUR,FIVE,SIX) VALUES
// ('31','34','3''5','36');
```