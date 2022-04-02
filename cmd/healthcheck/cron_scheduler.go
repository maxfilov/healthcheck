package main

import (
	cron "github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

type CronScheduler struct {
	scheduler *cron.Scheduler
	logger    log.FieldLogger
}

func MakeScheduler(delay time.Duration, services []Service, logger log.FieldLogger) (Lifecycle, error) {
	scheduler := cron.NewScheduler(time.UTC)
	for _, service := range services {
		logger.Infof("polling the status of %s", service.Print())
		_, err := scheduler.Every(int(delay / time.Millisecond)).Milliseconds().Do(service.Check)
		if err != nil {
			return nil, err
		}
	}
	return &CronScheduler{
		scheduler: scheduler,
		logger:    logger,
	}, nil
}

func (ss *CronScheduler) StartAsync() {
	ss.logger.Info("starting jobs scheduler")
	ss.scheduler.StartAsync()
}

func (ss *CronScheduler) Shutdown() error {
	ss.logger.Info("stopping scheduler")
	ss.scheduler.Stop()
	return nil
}

func (ss *CronScheduler) AwaitShutdown() error {
	ss.logger.Info("scheduler stopped")
	return nil
}
