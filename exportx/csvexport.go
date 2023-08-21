package exportx

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
)

type CsvExport struct {
	filehandle *os.File
	csvwriter  *csv.Writer
	headers    []string
}

func Create(export string) (csvx *CsvExport, err error) {
	csvx = &CsvExport{}
	csvx.filehandle, err = os.Create(export)
	if err != nil {
		log.Fatalf("csvx - failed to create file %s, %v", export, err)
		return
	}
	csvx.csvwriter = csv.NewWriter(csvx.filehandle)
	log.Infof("CSV export created: %s", export)
	return
}

func (csvx *CsvExport) Close() {
	if csvx.csvwriter != nil {
		csvx.csvwriter.Flush()
	}
	if csvx.filehandle != nil {
		csvx.filehandle.Close()
	}
}
func (csvx *CsvExport) Header(h ...string) {
	csvx.headers = h
	csvx.csvwriter.Write(h)
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
	csvx.csvwriter.Write(ar)
}
