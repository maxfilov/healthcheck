package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ConfigLoggingLevel struct {
	Root string `mapstrucutre:"root"`
}

type LoggingConfig struct {
	Level ConfigLoggingLevel `mapstructure:"level"`
}

type ScheduleConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Delay   int  `mapstructure:"delay"`
}

type ServiceDescription struct {
	Name string `mapstructure:"service-name"`
	Port int    `mapstructure:"port"`
	Path string `mapstructure:"path"`
}

type ClientServicesConfig struct {
	Services []ServiceDescription `mapstructure:"service-list"`
}

type GeoConfig struct {
	Service string `mapstructure:"service-name"`
	Port    int    `mapstructure:"port"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type PodConfig struct {
	Namespace string `mapstructure:"namespace"`
}

type Config struct {
	Pod              PodConfig            `mapstructure:"pod"`
	Server           ServerConfig         `mapstructure:"server"`
	Logging          LoggingConfig        `mapstructure:"logging"`
	Schedule         ScheduleConfig       `mapstructure:"schedule"`
	ClientServices   ClientServicesConfig `mapstructure:"client-services"`
	Geo              *GeoConfig           `mapstructure:"geo-healthcheck"`
	FailureThreshold int                  `mapstructure:"failure-threshold"`
}

func (config *Config) AsJson() string {
	jsonString, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "could not marshal configuration into json"
	}
	return string(jsonString)
}

func (config *Config) Verify() error {
	if !config.Schedule.Enabled {
		return nil
	}
	if config.Geo == nil && config.ClientServices.Services == nil {
		return errors.New("scheduling is enabled, but no services to check are provided")
	}
	if config.FailureThreshold <= 0 {
		return fmt.Errorf("only positive values are valid for config.failure-threshold: %d", config.FailureThreshold)
	}
	if config.Schedule.Enabled && config.Schedule.Delay <= 0 {
		return fmt.Errorf("only positive values are valid for config.schedule.delay: %d", config.FailureThreshold)
	}
	return nil
}
