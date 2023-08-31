package logutil

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var Logfile *os.File
var Loghandler *text.Handler

func OpenLogFile(logFilePath string) (err error) {

	if logFilePath == "" {
		Logfile = os.Stderr
	} else {
		Logfile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	}
	Loghandler = text.New(Logfile)
	log.SetHandler(Loghandler)
	return err
}

func CloseLogFile() {
	if Logfile != nil && Logfile != os.Stderr && Logfile != os.Stdout {
		Logfile.Close()
	}
}
