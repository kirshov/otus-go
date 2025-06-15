package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/queue"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)
	file, err := os.OpenFile(config.Logger.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	if err != nil {
		fmt.Println(err)
	}
	logg.SetOutput(file)

	// Storage.
	strg := storage.GetStorage(config.Storage.Type, config.Storage.DSN, config.Storage.Debug)

	// Rabbit
	r := rabbit.NewRabbit(config.Rabbit.DSN)
	closer, err := r.InitQueue(config.Rabbit.Queue)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			logg.Error(err.Error())
		}
	}(closer)

	if err != nil {
		logg.Error(err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	fmt.Println("scheduler is running...")

	ticker := time.NewTicker(config.Options.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if err := r.Stop(); err != nil {
				logg.Error("failed to stop rabbit: " + err.Error())
			}

			logg.Info("stopping service")
			return

		case <-ticker.C:
			logg.Info("run handler")
			handleItems(ctx, logg, strg, r, config)
		}
	}
}

func handleItems(ctx context.Context, logg logger.Logger, s storage.Storage, r queue.Handler, config Config) {
	items, err := s.List(ctx, config.Options.Days)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	for _, item := range items {
		jsonData, err := json.Marshal(item)
		if err != nil {
			logg.Error("failed marshal item " + item.ID + ": " + err.Error())
			continue
		}

		if err := r.Publish(config.Rabbit.Queue, jsonData); err != nil {
			logg.Error("failed marshal item " + item.ID + ": " + err.Error())
			continue
		}

		logg.Error("success publish item" + item.ID)
	}
}
