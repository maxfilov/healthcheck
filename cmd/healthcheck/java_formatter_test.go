package main

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestJavaFormatter_Format(t *testing.T) {
	formatter := JavaFormatter{}
	format, err := formatter.Format(&log.Entry{
		Time:    time.Time{},
		Level:   log.InfoLevel,
		Message: "hello message",
	})
	if err != nil {
		t.Errorf("Unexpected error occured during formatting '%s'", err.Error())
	}
	if string(format) != "0001-01-01T00:00:00Z [INFO   ] hello message\n" {
		t.Errorf(string(format))
	}
}
