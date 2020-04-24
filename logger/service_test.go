package logger

import (
	"bytes"
	"testing"
)

func TestNewLogger(t *testing.T) {
	buffer := new(bytes.Buffer)

	lg:=NewLogger(buffer)

	expected:="loglevel=Info deneme=123\n"

	lg.Info("deneme","123")

	if expected!=buffer.String(){
		t.Errorf("expected:%s but returned:%s",expected,buffer.String())
	}
}

func TestNewLogger_warning(t *testing.T) {
	buffer := new(bytes.Buffer)

	lg:=NewLogger(buffer)

	expected:="loglevel=Warning deneme=123\n"

	lg.Warning("deneme","123")

	if expected!=buffer.String(){
		t.Errorf("expected:%s but returned:%s",expected,buffer.String())
	}
}
func TestNewLogger_fatal(t *testing.T) {
	buffer := new(bytes.Buffer)

	lg:=NewLogger(buffer)

	expected:="loglevel=Fatal deneme=123\n"

	lg.Fatal("deneme","123")

	if expected!=buffer.String(){
		t.Errorf("expected:%s but returned:%s",expected,buffer.String())
	}
}