package exportx

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/apex/log"
)

type CsvExport struct {
	filehandle *os.File
	csvwriter  *csv.Writer
}

func Create(export string) (csvx *CsvExport, err error) {
	csvx = &CsvExport{}
	csvx.filehandle, err = os.Create(export)
	if err != nil {
		log.Fatalf("failed to open file %s, %v", export, err)
	}
	csvx.csvwriter = csv.NewWriter(csvx.filehandle)
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
	csvx.csvwriter.Write(h)
}

func (csvx *CsvExport) Write(p ...interface{}) {
	ar := make([]string, len(p))
	for i, x := range p {
		switch x := x.(type) {
		case string:
			ar[i] = x
		default:
			s := fmt.Sprint(x)
			ar[i] = s
		}
	}
	csvx.csvwriter.Write(ar)
}
