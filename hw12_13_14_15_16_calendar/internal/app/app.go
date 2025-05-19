package app

import (
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  logger.Logger
	Storage storage.Storage
}

func New(logger logger.Logger, storage storage.Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) GetLogger() logger.Logger {
	return a.Logger
}

func (a *App) GetStorage() storage.Storage {
	return a.Storage
}
