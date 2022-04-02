package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

type countingServiceMock struct {
	checkCounter int
	duration     time.Duration
}

func (c *countingServiceMock) Check() error {
	c.checkCounter++
	return nil
}

func (c *countingServiceMock) Print() string {
	return "counting service mock"
}

func TestCronScheduler_Lifecycle(t *testing.T) {
	service := countingServiceMock{duration: 0}
	logger := logrus.New()
	logger.SetOutput(bytes.NewBufferString(""))
	scheduler, err := MakeScheduler(300*time.Millisecond, []Service{&service}, logger)
	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err.Error())
	}
	scheduler.StartAsync()
	time.Sleep(500 * time.Millisecond)
	err = scheduler.Shutdown()
	if err != nil {
		t.Fatalf("Unexpected error during shutdown: '%s'", err.Error())
	}
	err = scheduler.AwaitShutdown()
	if err != nil {
		t.Fatalf("Unexpected error during shutdown awaiting: '%s'", err.Error())
	}
	if service.checkCounter != 2 {
		t.Fatalf("Unexpected service invocation number %d", service.checkCounter)
	}
}

func TestCronScheduler_OverSchedule(t *testing.T) {
	service := countingServiceMock{duration: 0}
	logger := logrus.New()
	logger.SetOutput(bytes.NewBufferString(""))
	scheduler, err := MakeScheduler(time.Second, []Service{&service}, logger)
	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err.Error())
	}
	scheduler.StartAsync()
	time.Sleep(500 * time.Millisecond)
	err = scheduler.Shutdown()
	if err != nil {
		t.Fatalf("Unexpected error during shutdown: '%s'", err.Error())
	}
	err = scheduler.AwaitShutdown()
	if err != nil {
		t.Fatalf("Unexpected error during shutdown awaiting: '%s'", err.Error())
	}
	if service.checkCounter != 1 {
		t.Fatalf("Unexpected service invocation number %d", service.checkCounter)
	}
}
