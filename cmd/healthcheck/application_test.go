package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

type fragileServiceStub struct{}

func (stub *fragileServiceStub) IsOk() bool    { return true }
func (stub *fragileServiceStub) Check() error  { return nil }
func (stub *fragileServiceStub) Print() string { return "fragile service" }

func TestMakeApplication(t *testing.T) {
	application, err := MakeApplication(&Config{
		Logging: LoggingConfig{Level: ConfigLoggingLevel{Root: "info"}},
		Server:  ServerConfig{Port: 8080},
	})
	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err.Error())
	}
	if len(application.lifecycle) != 1 {
		t.Errorf("Unexpected number of lifecycle components %d", len(application.lifecycle))
	}
	http.DefaultServeMux = new(http.ServeMux)
}

func TestMakeApplication_schedulingEnabled(t *testing.T) {
	config := Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Logging: LoggingConfig{
			Level: ConfigLoggingLevel{
				Root: "info",
			},
		},
		Schedule: ScheduleConfig{
			Enabled: true,
			Delay:   2000,
		},
		ClientServices: ClientServicesConfig{
			Services: []ServiceDescription{
				{
					Name: "name",
					Port: 80,
					Path: "/path",
				},
			},
		},
		Geo: &GeoConfig{
			Service: "http://something",
			Port:    80,
		},
		FailureThreshold: 0,
	}

	application, err := MakeApplication(&config)
	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err.Error())
	}
	if len(application.lifecycle) != 2 {
		t.Errorf("Unexpected number of lifecycle components %d", len(application.lifecycle))
	}
	http.DefaultServeMux = new(http.ServeMux)
}

func TestMakeApplication_invalidGeo(t *testing.T) {
	config := Config{
		Server:           ServerConfig{Port: 8080},
		Logging:          LoggingConfig{Level: ConfigLoggingLevel{Root: "info"}},
		Schedule:         ScheduleConfig{Enabled: true, Delay: 2000},
		ClientServices:   ClientServicesConfig{},
		Geo:              &GeoConfig{Service: "%$", Port: 80},
		FailureThreshold: 0,
	}

	_, err := MakeApplication(&config)
	defer func() {
		http.DefaultServeMux = new(http.ServeMux)
	}()
	if err == nil {
		t.Fatal("Unexpected nil error")
	}
}

func TestMakeApplication_invalidServices(t *testing.T) {
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))

	config := Config{
		Server:   ServerConfig{Port: 8080},
		Logging:  LoggingConfig{Level: ConfigLoggingLevel{Root: "info"}},
		Schedule: ScheduleConfig{Enabled: true, Delay: 2000},
		ClientServices: ClientServicesConfig{
			Services: []ServiceDescription{
				{
					Name: "%$",
					Port: 80,
					Path: "/health",
				},
			},
		},
		Geo:              nil,
		FailureThreshold: 0,
	}

	_, err := MakeApplication(&config)
	defer func() {
		http.DefaultServeMux = new(http.ServeMux)
	}()
	if err == nil {
		t.Fatal("Unexpected nil error")
	}
}

func TestMakeApplication_invalidLogLevelInConfiguration(t *testing.T) {
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))

	application, err := MakeApplication(&Config{
		Logging: LoggingConfig{
			Level: ConfigLoggingLevel{
				Root: "invalid",
			},
		},
	})
	if err == nil {
		t.Fatalf("Unexpected nil as error")
	}
	if application != nil {
		t.Errorf("Unexpected non nil application")
	}
	http.DefaultServeMux = new(http.ServeMux)
}

func TestApplication_toFragiles(t *testing.T) {
	fragileServices := []FragileService{&fragileServiceStub{}}
	fragiles := toFragiles(fragileServices)
	if len(fragiles) != 1 {
		t.Errorf("Unexpected number of fragiles in result: %d", len(fragiles))
	}
}

func TestApplication_toServices(t *testing.T) {
	fragileServices := []FragileService{&fragileServiceStub{}}
	services := toServices(fragileServices)
	if len(services) != 1 {
		t.Errorf("Unexpected number of fragiles in result: %d", len(services))
	}
}

type lifecycleMock struct {
}

func (l *lifecycleMock) StartAsync() {
}

func (l *lifecycleMock) Shutdown() error {
	return fmt.Errorf("test shutdown error")
}

func (l *lifecycleMock) AwaitShutdown() error {
	return fmt.Errorf("test await shutdown error")
}

func TestApplication_setUpInterruptionHandler(t *testing.T) {
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))
	sigs := make(chan os.Signal, 1)
	application := Application{
		lifecycle:     []Lifecycle{&lifecycleMock{}},
		logger:        logger,
		interruptions: sigs,
	}
	wg := application.setUpInterruptHandler()
	sigs <- syscall.SIGINT
	wg.Wait()
	http.DefaultServeMux = new(http.ServeMux)
}

func TestApplication_makeServiceList(t *testing.T) {
	serviceDescriptions := []ServiceDescription{
		{
			Name: "service",
			Port: 8080,
			Path: "/health",
		},
	}
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))
	list, err := makeServiceList(3, serviceDescriptions, logger)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if len(list) != 1 {
		t.Errorf("Unexpected number of services: %d", len(list))
	}
	service := list[0]
	watchfulDecorator := service.(*HopefulProxy)
	if watchfulDecorator.threshold != 3 {
		t.Errorf("Unexpected threshold value: %d", watchfulDecorator.threshold)
	}
	loggingDecorator := watchfulDecorator.backend.(*LoggingServiceDecorator)
	if loggingDecorator.logger != logger {
		t.Errorf("Unexpected logger address")
	}
	simpleService := loggingDecorator.backend.(*SimpleService)
	if simpleService.endpoint != "http://service:8080/health" {
		t.Errorf("Unexpected endpoint: %s", simpleService.endpoint)
	}
}

func TestApplication_Lifecycle(t *testing.T) {
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))
	application := &Application{
		lifecycle:     []Lifecycle{&lifecycleMock{}},
		logger:        logger,
		interruptions: make(chan os.Signal, 1),
	}
	go func() {
		time.Sleep(time.Second)
		application.interruptions <- syscall.SIGINT
	}()
	application.Run()
	http.DefaultServeMux = new(http.ServeMux)
}

func TestApplication_makeGeoService(t *testing.T) {
	logger := log.New()
	logger.SetOutput(bytes.NewBufferString(""))
	config := Config{
		Geo: &GeoConfig{
			Service: "http://service",
			Port:    80,
		},
		FailureThreshold: 3}
	service, err := makeGeoService(&config, logger)
	if err != nil {
		t.Fatalf("Unexpecgted error: %s", err.Error())
	}
	watchfulService := service.(*HopefulProxy)
	if watchfulService.threshold != 3 {
		t.Errorf("Unexpected threshold: %d", watchfulService.threshold)
	}
	loggingDecorator := watchfulService.backend.(*LoggingServiceDecorator)
	if loggingDecorator.logger != logger {
		t.Error("Unexpected logger address")
	}
	simpleService := loggingDecorator.backend.(*SimpleService)
	simpleServiceDescription := simpleService.Print()
	if simpleServiceDescription != "service at 'http://service:80/health'" {
		t.Errorf("Unexpected description: %s", simpleServiceDescription)
	}
}
