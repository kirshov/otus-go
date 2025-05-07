package app

import (
	"context"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

type Storage interface {
	Add(ctx context.Context, event domain.Event) error
	Update(ctx context.Context, event domain.Event) error
	Remove(ctx context.Context, id string) error
	List(ctx context.Context, days int) ([]domain.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event domain.Event) error {
	err := a.storage.Add(ctx, event)
	return err
}
