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
	t := []interface{}{"log", "Info"}
	t = append(t, i...)
	return l.Logger.Log(t...)
}
func (l logger) Warning(i ...interface{}) error {
	t := []interface{}{"log", "Warning"}
	t = append(t, i...)
	return l.Logger.Log(t...)
}
func (l logger) Fatal(i ...interface{}) error {
	t := []interface{}{"log", "Fatal"}
	t = append(t, i...)
	return l.Logger.Log(t...)
}
