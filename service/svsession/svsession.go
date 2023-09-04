package svsession

import (
	"fmt"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

type EdgeConfig struct {
	Edgerc    string
	Section   string
	AccountID string
	Logger    log.Interface
}

func CreateLogger(filename string, level log.Level) *log.Logger {
	// Replace "your_log_file.log" with the desired file name and path
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	logger := &log.Logger{
		Handler: text.New(file),
		Level:   log.DebugLevel, // Set log level as desired (e.g., log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel)
	}

	return logger
}

func NewSession(param *EdgeConfig) (s session.Session, err error) {
	edgercOps := []edgegrid.Option{edgegrid.WithEnv(true)}
	if param.Edgerc != "" {
		edgercOps = append(edgercOps, edgegrid.WithFile(param.Edgerc))
	}

	if param.Section != "" {
		edgercOps = append(edgercOps, edgegrid.WithSection(param.Section))
	}

	edgerc, err := edgegrid.New(edgercOps...)
	edgerc.AccountKey = param.AccountID

	if err != nil {
		log.Fatalf("edgerc error: %v", err)
		return
	}

	if edgerc.Host == "" {
		//err = fmt.Errorf("EdgeRC section not found: %s", akutil.StructToColumns(*param))
		err = fmt.Errorf("edgerc section not found, edgerc:'%s, section %s' ", param.Edgerc, param.Section)
		return
	}
	//fmt.Println(edgerc)
	s, err = session.New(
		session.WithSigner(edgerc),
		session.WithHTTPTracing(true),
		session.WithLog(param.Logger),
	)
	return
}
