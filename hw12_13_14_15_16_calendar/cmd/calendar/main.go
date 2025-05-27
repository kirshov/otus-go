package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Logg.
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
	calendar := app.New(logg, strg)
	// testExecute(config, calendar, strg)

	server := internalhttp.NewServer(calendar)
	grpcServer := internalgrpc.NewServer(calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	fmt.Println("calendar is running...")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := server.Start(config.Server.address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpcServer.Start(config.GrpcServer.address); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()
	logg.Info("stopping servers")

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
	grpcServer.Stop()

	os.Exit(1) //nolint:gocritic
}

/*
// todo проверка работы db, передавать на функциональные тесты после соответствующего урока.
func testExecute(config Config, calendar *app.App, strg storage.Storage) {
	timeout := time.Duration(config.Storage.Timeout * int64(time.Second))
	storageCtx, storageCancel := context.WithTimeout(context.Background(), timeout)
	defer storageCancel()
	event := domain.Event{
		ID:          "111-111-11",
		Title:       "title",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(5 * time.Hour),
		Description: "description",
		UserID:      "U0001",
		NotifyDays:  5,
	}
	err := calendar.CreateEvent(storageCtx, event)
	if err != nil {
		fmt.Println(err)
	}
	l, _ := strg.List(storageCtx, 0)
	fmt.Println(l)
	event.DateStart = time.Now().Add(555 * time.Hour)
	event.Title = "@SASDASDASD"
	err = strg.Update(storageCtx, event)
	if err != nil {
		fmt.Println(err)
	}
	l, err = strg.List(storageCtx, 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(l)

	l, err = strg.List(storageCtx, 25)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(l)
	err = strg.Remove(storageCtx, event.ID)
	if err != nil {
		fmt.Println(err)
	}
}*/
