package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	lifecycle     []Lifecycle
	logger        *logrus.Logger
	interruptions chan os.Signal
}

func MakeApplication(config *Config) (*Application, error) {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&JavaFormatter{})
	logger.Infof("started application with log level '%s'", logger.GetLevel())
	logger.Infof("configuration:\n%s", config.AsJson())
	logLevel, err := logrus.ParseLevel(config.Logging.Level.Root)
	if err != nil {
		return nil, err
	}
	logger.Infof("switching to log level '%s'", logLevel)
	logger.SetLevel(logLevel)

	var fragileServices []FragileService
	var geoService FragileService
	var scheduler Lifecycle
	if config.Schedule.Enabled {
		fragileServices, err = makeServiceList(config.FailureThreshold, config.ClientServices.Services, logger)
		if err != nil {
			return nil, err
		}
		geoService, err = makeGeoService(config, logger)
		if err != nil {
			return nil, err
		}
		if geoService != nil {
			fragileServices = append(fragileServices, geoService)
		}
		delay := time.Duration(int(time.Millisecond) * config.Schedule.Delay)
		scheduler, err = MakeScheduler(delay, toServices(fragileServices), logger)
	}

	healthHandler, err := MakeHealthHandler(
		config.Pod.Namespace,
		toFragiles(fragileServices),
		geoService)
	if err != nil {
		return nil, err
	}
	lifecycle := []Lifecycle{MakeServer(config.Server.Port, healthHandler, logger)}
	if scheduler != nil {
		lifecycle = append(lifecycle, scheduler)
	}
	return &Application{
		lifecycle:     lifecycle,
		logger:        logger,
		interruptions: make(chan os.Signal, 1),
	}, nil
}

func makeGeoService(config *Config, logger *logrus.Logger) (FragileService, error) {
	configGeo := config.Geo
	if configGeo == nil {
		return nil, nil
	}
	endpoint := fmt.Sprintf("%s:%d/health", configGeo.Service, configGeo.Port)
	service, err := MakeSimpleService(endpoint, &http.Client{})
	if err != nil {
		return nil, err
	}
	loggingDecorator := MakeLoggingServiceDecorator(service, logger)
	watchfulDecorator := MakeHopefulProxy(loggingDecorator, config.FailureThreshold)
	return watchfulDecorator, nil
}

func makeServiceList(threshold int, serviceDescriptions []ServiceDescription, logger logrus.FieldLogger) ([]FragileService, error) {
	var result []FragileService
	for _, srvDesc := range serviceDescriptions {
		//This will be used inside a service mesh, it should encrypt all communications
		//goland:noinspection HttpUrlsUsage
		endpoint := fmt.Sprintf("http://%s:%d%s", srvDesc.Name, srvDesc.Port, srvDesc.Path)
		actuatorService, err := MakeSimpleService(endpoint, &http.Client{})
		if err != nil {
			return nil, err
		}
		loggingDecorator := MakeLoggingServiceDecorator(actuatorService, logger)
		watchfulDecorator := MakeHopefulProxy(loggingDecorator, threshold)
		result = append(result, watchfulDecorator)
	}
	return result, nil
}

func toFragiles(fragileServices []FragileService) []Fragile {
	var fragiles []Fragile
	for _, fragileService := range fragileServices {
		fragiles = append(fragiles, fragileService)
	}
	return fragiles
}

func toServices(fragileServices []FragileService) []Service {
	var services []Service
	for _, fragileService := range fragileServices {
		services = append(services, fragileService)
	}
	return services
}

func (application *Application) Run() {
	for _, lifecycle := range application.lifecycle {
		lifecycle.StartAsync()
	}

	wg := application.setUpInterruptHandler()

	wg.Wait()

	for _, lifecycle := range application.lifecycle {
		err := lifecycle.AwaitShutdown()
		if err != nil {
			application.logger.Error(err)
		}
	}
}

func (application *Application) setUpInterruptHandler() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	signal.Notify(application.interruptions, syscall.SIGINT)
	go func() {
		sig := <-application.interruptions
		application.logger.Infof("received '%s' signal", sig)
		wg.Done()
		for _, lifecycle := range application.lifecycle {
			err := lifecycle.Shutdown()
			if err != nil {
				application.logger.Errorf("'%s' shutdown failed: '%s'", lifecycle, err)
			}
		}
	}()
	return wg
}
