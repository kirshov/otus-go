package storage

import (
	"log"

	memorystorage "github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	MEMORY string = "memory"
	PSQL   string = "psql"
)

func GetStorage(t, dsn string, debug bool) Storage {
	switch t {
	case MEMORY:
		return memorystorage.New()
	case PSQL:
		return sqlstorage.New(dsn, debug)
	default:
		log.Fatal("unsupported storage type")
		return nil
	}
}
