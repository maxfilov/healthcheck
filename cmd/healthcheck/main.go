package main

import (
	"github.com/spf13/viper"
	"strings"
)

func main() {
	config, err := assembleConfiguration()
	if err != nil {
		panic(err)
	}

	app, err := MakeApplication(config)
	if err != nil {
		panic(err)
	}
	app.Run()
}

func assembleConfiguration() (*Config, error) {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("logging.level.root", "info")
	viper.SetDefault("schedule.enabled", "false")
	viper.SetDefault("pod.namespace", "unknown")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/healthcheck")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	viper.AutomaticEnv()
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	if err := config.Verify(); err != nil {
		return nil, err
	}
	return &config, nil
}
