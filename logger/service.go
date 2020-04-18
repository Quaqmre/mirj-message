package logger

import (
	"io"

	"github.com/go-kit/kit/log"
)

type Service interface {
	Info(...interface{}) error
	Warning(...interface{}) error
	Fatal(...interface{}) error
}

type logger struct {
	log.Logger
}

func NewLogger(w io.Writer) logger {
	lg := logger{}
	lg.Logger = log.NewLogfmtLogger(w)
	return lg
}

func (l logger) Info(i ...interface{}) error {
	i = append(i, "loglevel")
	i = append(i, "info")

	return l.Logger.Log(i...)
}
func (l logger) Warning(i ...interface{}) error {
	return l.Logger.Log("loglevel", "warning", i)
}
func (l logger) Fatal(i ...interface{}) error {
	return l.Logger.Log("loglevel", "fatal", i)
}
