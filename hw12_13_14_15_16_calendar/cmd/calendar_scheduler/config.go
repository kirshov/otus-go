package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Options OptionsConf
	Rabbit  RabbitConf
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

type OptionsConf struct {
	Interval time.Duration
	Days     int
}

type RabbitConf struct {
	DSN   string
	Queue string
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

	rabbitHost := viper.GetString("rabbit.host")
	rabbitPort := viper.GetString("rabbit.port")
	rabbitUsername := viper.GetString("rabbit_username")
	rabbitPassword := viper.GetString("rabbit_password")

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
		Options: OptionsConf{
			Interval: viper.GetDuration("options.interval"),
			Days:     viper.GetInt("options.days"),
		},
		Rabbit: RabbitConf{
			DSN:   "amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/",
			Queue: viper.GetString("rabbit.queue"),
		},
	}
}
