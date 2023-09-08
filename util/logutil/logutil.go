package logutil

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/text"
)

// LogConfig can be use to pass in a set of handlers
type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// SubLogHandler is a handler that filters out logs below a certain level.
type SubLogHandler struct {
	Handler log.Handler // the wrapped handler
	Level   log.Level   // the minimum level to pass
	Close   func()      // function to call on close
}

type MultiLogHandler struct {
	SubHandlers []SubLogHandler
	MaxLevel    log.Level
}

// NewMultiLogHandler provides a new instance of a multi log handler
func NewMultiLogHandler() (lf *MultiLogHandler) {
	lf = &MultiLogHandler{
		MaxLevel: log.FatalLevel,
	}
	return
}

// OpenFromConfig will read the logging information from a config structure passed in as parameter
func (mh *MultiLogHandler) OpenFromConfig(lc []LogConfig) {
	var err error
	if len(lc) == 0 {
		mh.Open(log.WarnLevel, "")
	} else {
		for _, c := range lc {
			level := log.WarnLevel
			if c.Level != "" {
				level, err = log.ParseLevel(c.Level)
				if err != nil {
					log.Fatalf("%v", err)
				}
			}
			mh.Open(level, c.File)
		}
	}
}

// Open an (additional) log-handler at the provided level, uses cli-handler when the filename is empty, otherwise uses a text handler
func (mh *MultiLogHandler) Open(level log.Level, filename string) {
	lf := SubLogHandler{}
	lf.Level = level
	if level < mh.MaxLevel {
		mh.MaxLevel = level
	}
	if filename == "" {
		lf.Handler = cli.New(os.Stdout)
	} else {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		lf.Close = func() { file.Close() }
		lf.Handler = text.New(file)
	}
	mh.SubHandlers = append(mh.SubHandlers, lf)
}

// Close closes the underying streams
func (mh *MultiLogHandler) Close() {
	for _, f := range mh.SubHandlers {
		if f.Close != nil {
			f.Close()
		}
	}
}

// HandleLog implements the log.Handler interface.
func (mh *MultiLogHandler) HandleLog(e *log.Entry) (err error) {
	for _, f := range mh.SubHandlers {
		if e.Level >= f.Level {
			x := f.Handler.HandleLog(e)
			if x != nil {
				err = x
			}
		}
	}
	return
}
