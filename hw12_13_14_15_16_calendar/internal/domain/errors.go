package domain

import "fmt"

type EventExistsError struct {
	EventID string
}

type EventNotExistsError struct {
	EventID string
}

func (e EventExistsError) Error() string {
	return fmt.Sprintf("event %s already exists", e.EventID)
}

func (e EventNotExistsError) Error() string {
	return fmt.Sprintf("event %s is not exists", e.EventID)
}
