package domain

import (
	"time"
)

type Event struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	DateStart   time.Time `db:"date_start"`
	DateEnd     time.Time `db:"date_end"`
	Description string    `db:"description"`
	UserID      string    `db:"user_id"`
	NotifyDays  int       `db:"notify_days"`
}
