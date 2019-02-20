package ci

import (
	"io"
	"log"
	"os"
)

type StdLogger struct {
	log *log.Logger
}

func NewStd() Logger {
	return &StdLogger{log.New(os.Stdout, "", log.Ltime)}
}

func (log *StdLogger) Named(name string) Logger {
	return log
}

func (log *StdLogger) Output() (stdout, stderr io.Writer) {
	return os.Stdout, os.Stderr
}

func (log *StdLogger) Print(v ...interface{}) {
	log.log.Print(v...)
}

func (log *StdLogger) Printf(format string, v ...interface{}) {
	log.log.Printf(format, v...)
}

func (log *StdLogger) Error(v ...interface{}) {
	log.log.Print(v...)
}

func (log *StdLogger) Errorf(format string, v ...interface{}) {
	log.log.Printf(format, v...)
}
