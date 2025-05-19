package domain

import (
	"time"
)

type Event struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	DateStart   time.Time `db:"date_start" json:"date_start"`
	DateEnd     time.Time `db:"date_end" json:"date_end"`
	Description string    `db:"description" json:"description"`
	UserID      string    `db:"user_id" json:"user_id"`
	NotifyDays  int       `db:"notify_days" json:"notify_days"`
}
