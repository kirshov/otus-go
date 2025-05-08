package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
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

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	strg := storage.GetStorage(config.Storage.Type, config.Storage.DSN)
	calendar := app.New(logg, strg)
	// testExecute(config, calendar, strg)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx, config.Server.address); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
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
