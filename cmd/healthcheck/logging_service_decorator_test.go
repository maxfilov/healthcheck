package main

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"

	"testing"
)

func TestMakeLoggingDecorator(t *testing.T) {
	decorator := MakeLoggingServiceDecorator(&ServiceStub{}, log.New())
	if decorator == nil {
		t.Errorf("Decorator should not be nil")
	}
}

type collectingFormatter struct {
	entries []*log.Entry
}

func (c *collectingFormatter) Format(entry *log.Entry) ([]byte, error) {
	c.entries = append(c.entries, entry)
	return []byte{}, nil
}

func TestLoggingDecorator_Check_OnSuccess(t *testing.T) {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	formatter := collectingFormatter{entries: []*log.Entry{}}
	logger.SetFormatter(&formatter)
	logger.SetOutput(bytes.NewBufferString(""))

	decorator := MakeLoggingServiceDecorator(&ServiceStub{}, logger)
	err := decorator.Check()
	if err != nil {
		t.Errorf("Return value should be nil when backend returns nil")
	}
	if len(formatter.entries) != 1 {
		t.Errorf("Unexpected number of log entries %d, expected=%d", len(formatter.entries), 1)
	}
	entry := formatter.entries[0]
	if entry.Message != "service stub is OK" {
		t.Errorf("Unexpected message: '%s'", entry.Message)
	}
	if entry.Level != log.DebugLevel {
		t.Errorf("Unexpected log entry level '%s', expected='%s'", entry.Level, log.DebugLevel)
	}
}

func TestLoggingDecorator_Check_OnError(t *testing.T) {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	formatter := collectingFormatter{entries: []*log.Entry{}}
	logger.SetFormatter(&formatter)
	logger.SetOutput(bytes.NewBufferString(""))

	decorator := MakeLoggingServiceDecorator(&ServiceStub{Err: errors.New("some error text")}, logger)
	err := decorator.Check()
	if err == nil {
		t.Errorf("Return value should be non nil when backend returns non nil")
	}
	if err.Error() != "some error text" {
		t.Errorf("Unexpected changes in the text")
	}
	if len(formatter.entries) != 1 {
		t.Errorf("Unexpected number of log entries %d, expected=%d", len(formatter.entries), 1)
	}
	entry := formatter.entries[0]
	if entry.Message != "service stub is not OK: some error text" {
		t.Errorf("Unexpected message: '%s'", entry.Message)
	}
	if entry.Level != log.WarnLevel {
		t.Errorf("Unexpected log entry level '%s', expected='%s'", entry.Level, log.DebugLevel)
	}
}

func TestLoggingDecorator_Print(t *testing.T) {
	s := MakeLoggingServiceDecorator(&ServiceStub{}, log.New()).Print()
	if s != "logging decorator for service stub" {
		t.Errorf("Unexpected description of the logging decorator")
	}
}
