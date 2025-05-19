package storage

import (
	"context"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type Storage interface {
	GetByID(ctx context.Context, id string) (domain.Event, error)
	Add(ctx context.Context, event domain.Event) error
	Update(ctx context.Context, event domain.Event) error
	Remove(ctx context.Context, id string) error
	List(ctx context.Context, days int) ([]domain.Event, error)
}
