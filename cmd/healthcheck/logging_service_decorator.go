package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type LoggingServiceDecorator struct {
	backend Service
	logger  log.FieldLogger
}

func MakeLoggingServiceDecorator(backend Service, logger log.FieldLogger) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		backend: backend,
		logger:  logger,
	}
}

func (l *LoggingServiceDecorator) Check() error {
	err := l.backend.Check()
	if err == nil {
		l.logger.Debugf("%s is OK", l.backend.Print())
	} else {
		l.logger.Warnf("%s is not OK: %s", l.backend.Print(), err)
	}
	return err
}

func (l *LoggingServiceDecorator) Print() string {
	return fmt.Sprintf("logging decorator for %s", l.backend.Print())
}
