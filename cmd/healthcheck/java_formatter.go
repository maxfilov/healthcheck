package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type JavaFormatter struct{}

func (j JavaFormatter) Format(entry *log.Entry) ([]byte, error) {
	sprintf := fmt.Sprintf("%s [%-7s] %s\n", entry.Time.Format(time.RFC3339), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(sprintf), nil
}
