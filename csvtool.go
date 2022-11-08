package csvtool

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type csvTool struct {
	fileName   string
	delimiter  rune
	prefixOut  string
	table      string
	header     string
	rowsByFile int
	columns    Column
	start      time.Time
	base       string
	totalRows  int
	rows       *[][]string
}

type CSVTool interface {
	ToSQL(tableName, outFileName string, rowsByFile int) error
	RemoveColumn(columns ...string) CSVTool
	SplitCSV(rowsByFile int) error
}

type Column map[string]int

func NewCSVTool(fileName string, delimiter rune) CSVTool {
	return &csvTool{
		fileName,
		delimiter,
		"",
		"",
		"",
		0,
		make(Column),
		time.Now(),
		"INSERT INTO %s (%s) VALUES\n",
		0,
		&[][]string{},
	}
}

func (c *csvTool) RemoveColumn(columns ...string) CSVTool {
	for _, col := range columns {
		c.columns[col] = -1
	}
	return c
}

func (c *csvTool) checkHeader(firstRow []string) error {
	var headers []string
	for i, h := range firstRow {
		v, ok := c.columns[h]
		if ok && v < 0 {
			c.columns[h] = i
			continue
		}
		headers = append(headers, h)
	}
	for k, v := range c.columns {
		if v < 0 {
			return errors.New(fmt.Sprintf("ERROR: the column: \"%s\" doesn't exist in \"%s\" file.", k, c.fileName))
		}
	}
	c.header = strings.Join(headers, ",")
	return nil
}

func (c *csvTool) loadCsv() error {
	f, err := os.Open(c.fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = rune(c.delimiter)

	data, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	c.rows = &data
	c.totalRows = len(*c.rows)
	return nil
}

func (c *csvTool) ToSQL(tableName, outFileName string, rowsByFile int) error {
	fmt.Println("Loading file...")
	c.start = time.Now()

	c.table = tableName
	c.prefixOut = outFileName
	c.rowsByFile = rowsByFile
	var bucket []string

	err := c.loadCsv()
	if err != nil {
		fmt.Println("Error during load csv:", c.fileName, "Delimiter:", c.delimiter)
		return err
	}

	err = c.checkHeader((*c.rows)[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("File loaded, processing", c.totalRows, "rows")
	for i := 1; i < c.totalRows; i++ {
		var tempRow []string
	INNER:
		for j := 0; j < len((*c.rows)[i]); j++ {
			cols, ok := c.columns[(*c.rows)[0][j]]
			if ok && cols == j {
				continue INNER
			}
			tempRow = append(tempRow, (strings.Replace((*c.rows)[i][j], "'", "''", -1)))
		}
		(*c.rows)[i] = tempRow
		var endLine string
		if endLine = "'),\n"; (i%c.rowsByFile == 0) || i == c.totalRows-1 {
			endLine = "');\n"
		}

		union := strings.Join((*c.rows)[i], "','")
		bucket = append(bucket, "('"+union+endLine)
	}

	fmt.Println("Processing finished, saving file")
	c.header = fmt.Sprintf(c.base, c.table, c.header)
	err = c.saveFiles(&bucket, "sql")
	if err != nil {
		return err
	}

	fmt.Println("Files saved", time.Since(c.start))
	return nil
}

func (c *csvTool) saveFiles(rows *[]string, ex string) error {
	lap := 1
	var endRows int
	for i := 0; i < c.totalRows; i += c.rowsByFile {
		if endRows = i + c.rowsByFile; i+c.rowsByFile > c.totalRows {
			endRows = c.totalRows
		}
		err := c.saveFile(fmt.Sprintf("%s_%02d.%s", c.prefixOut, lap, ex), []byte(c.header+strings.Join((*rows)[i:endRows], "")))
		if err != nil {
			return err
		}
		lap++
	}
	return nil
}

func (c *csvTool) saveFile(filename string, b []byte) error {
	return os.WriteFile(filename, b, 0644)
}

func (c *csvTool) SplitCSV(rowsByFile int) error {
	splited := strings.Split(c.fileName, ".")
	c.prefixOut = splited[0]
	fmt.Println("Loading file...")
	c.start = time.Now()
	c.rowsByFile = rowsByFile
	var bucket []string

	err := c.loadCsv()
	if err != nil {
		fmt.Println("Error during load csv:", c.fileName, "Delimiter:", c.delimiter)
		return err
	}

	err = c.checkHeader((*c.rows)[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("File loaded, processing", c.totalRows, "rows")

	for i := 1; i < c.totalRows; i++ {
		var tempRow []string
	INNER:
		for j := 0; j < len((*c.rows)[i]); j++ {
			cols, ok := c.columns[(*c.rows)[0][j]]
			if ok && cols == j {
				continue INNER
			}
			tempRow = append(tempRow, (*c.rows)[i][j])
		}
		tempRow[len(tempRow)-1] += "\n"
		(*c.rows)[i] = tempRow
		bucket = append(bucket, strings.Join((*c.rows)[i], ","))
	}

	fmt.Println("Processing finished, saving file")
	c.header += "\n"
	err = c.saveFiles(&bucket, "csv")
	if err != nil {
		return err
	}

	fmt.Println("Files saved", time.Since(c.start))
	return nil
}
