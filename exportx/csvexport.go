package exportx

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
)

type CsvExport struct {
	filename string
	rows     [][]string
	headers  []string
}

func Create(export string) (csvx *CsvExport, err error) {
	csvx = &CsvExport{
		filename: export,
		rows:     [][]string{},
		headers:  []string{},
	}
	return
}

func (csvx *CsvExport) Close() {
	filehandle, err := os.Create(csvx.filename)
	if err != nil {
		log.Fatalf("csvx - failed to create file %s, %v", csvx.filename, err)
		return
	}
	defer filehandle.Close()

	csvwriter := csv.NewWriter(filehandle)
	defer csvwriter.Flush()

	csvwriter.Write(csvx.headers)
	csvwriter.WriteAll(csvx.rows)

	log.Infof("CSV export ready")
}
func (csvx *CsvExport) Header(h ...string) {
	csvx.headers = h
}

func (csvx *CsvExport) Write(p ...interface{}) {
	ar := make([]string, len(p))
	for i, x := range p {
		switch x := x.(type) {
		case string:
			ar[i] = x
		default:
			// The default are handled quick and dirty as you can see
			s := fmt.Sprint(x)
			s = strings.TrimPrefix(s, "map")
			if s != "0" && s != "0.0" && s != "0001-01-01 00:00:00 +0000 UTC" && s != "[]" {
				ar[i] = s
			} else {
				ar[i] = ""
			}
		}
	}
	if len(ar) > len(csvx.headers) {
		ar = ar[:len(csvx.headers)]
	}
	csvx.rows = append(csvx.rows, ar)
}
