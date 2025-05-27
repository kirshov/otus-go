package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib" // register pgx driver for sql
	"github.com/jmoiron/sqlx"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type Storage struct {
	conn *sqlx.DB
}

func New(dsn string, debug bool) *Storage {
	con, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	ctx := context.Background()
	err = con.PingContext(ctx)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	_ = debug

	return &Storage{
		conn: con,
	}
}

func (s *Storage) Add(ctx context.Context, event domain.Event) (string, error) {
	if event.ID == "" {
		event.ID = uuid.New().String()
	} else {
		e, err := s.GetByID(ctx, event.ID)
		if err != nil && errors.Is(err, domain.EventNotExistsError{}) {
			return "", err
		}

		if e.ID == event.ID {
			return "", domain.EventExistsError{EventID: event.ID}
		}
	}

	query := "INSERT INTO events VALUES (:id, :title, :date_start, :date_end, :description, :user_id, :notify_days)"
	if _, err := s.conn.NamedExecContext(ctx, query, event); err != nil {
		return "", err
	}
	return event.ID, nil
}

func (s *Storage) Update(ctx context.Context, event domain.Event) error {
	query := "UPDATE events SET " +
		"title = :title, date_start = :date_start, date_end = :date_end, " +
		"description = :description, user_id = :user_id, notify_days = :notify_days" +
		" WHERE id = :id"
	if _, err := s.conn.NamedExecContext(ctx, query, event); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Remove(ctx context.Context, id string) error {
	query := "DELETE FROM events WHERE id = :id"
	if _, err := s.conn.NamedExecContext(ctx, query, map[string]interface{}{"id": id}); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, id string) (domain.Event, error) {
	var event domain.Event

	query := "SELECT * FROM events WHERE id = $1"
	if err := s.conn.GetContext(ctx, &event, query, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return event, err
	}

	if event.ID != id {
		return event, domain.EventNotExistsError{EventID: id}
	}

	return event, nil
}

func (s *Storage) List(ctx context.Context, days int) ([]domain.Event, error) {
	query := "SELECT * FROM events WHERE date_start > NOW()"
	args := map[string]interface{}{}

	if days > 0 {
		query += " AND date_start <= NOW() + INTERVAL '1 day' * :days "
		args["days"] = days
	}

	rows, err := s.conn.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
