package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string
	File  string
}

type StorageConf struct {
	Type    string
	Timeout int64
	DSN     string
	Debug   bool
}

type ServerConf struct {
	address string
}

func NewConfig(configFile string) Config {
	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		log.Fatal(fmt.Errorf("fatal error .env: %w", err))
	}
	viper.AutomaticEnv()

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	if err := viper.MergeInConfig(); err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}

	return Config{
		Logger: LoggerConf{
			Level: viper.GetString("logger.level"),
			File:  viper.GetString("logger.file"),
		},
		Storage: StorageConf{
			Type:    viper.GetString("storage.type"),
			Timeout: viper.GetInt64("storage.timeout"),
			DSN:     viper.GetString("app_postgres_dsn"),
			Debug:   viper.GetBool("storage.debug"),
		},
		Server: ServerConf{
			address: viper.GetString("server.host") + ":" + viper.GetString("server.port"),
		},
	}
}
