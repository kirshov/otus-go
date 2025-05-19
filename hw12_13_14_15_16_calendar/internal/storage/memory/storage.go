package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type Storage struct {
	data map[string]domain.Event
	mu   sync.RWMutex
}

func New() *Storage {
	return &Storage{
		data: make(map[string]domain.Event),
	}
}

func (s *Storage) GetByID(ctx context.Context, id string) (domain.Event, error) {
	var event domain.Event

	select {
	case <-ctx.Done():
		return event, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.data {
		if v.ID == id {
			return v, nil
		}
	}

	return event, domain.EventExistsError{EventID: id}
}

func (s *Storage) Add(ctx context.Context, event domain.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[event.ID]; ok {
		return domain.EventExistsError{EventID: event.ID}
	}

	s.data[event.ID] = event
	return nil
}

func (s *Storage) Update(ctx context.Context, event domain.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[event.ID]; !ok {
		return domain.EventNotExistsError{EventID: event.ID}
	}

	s.data[event.ID] = event
	return nil
}

func (s *Storage) Remove(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return fmt.Errorf("event %s is not exists", id)
	}

	delete(s.data, id)
	return nil
}

func (s *Storage) List(ctx context.Context, days int) ([]domain.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	maxDate := time.Now().Add(time.Duration(days) * 24 * time.Hour)

	events := make([]domain.Event, 0, len(s.data))
	for _, v := range s.data {
		if days > 0 {
			if v.DateStart.After(maxDate) {
				continue
			}
		}

		events = append(events, v)
	}

	return events, nil
}
