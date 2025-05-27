package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	event1 := domain.Event{ID: "1", Title: "title1"}
	event2 := domain.Event{ID: "2", Title: "title2"}
	updatedEvent := domain.Event{ID: event1.ID, Title: "updatedTitle"}

	storage := New()
	ctx := context.Background()

	t.Run("Add", func(t *testing.T) {
		id, err := storage.Add(ctx, event1)
		require.NotEmpty(t, id)
		require.NoError(t, err)

		l, err := storage.List(ctx, 0)
		require.NoError(t, err)
		require.ElementsMatch(t, []domain.Event{event1}, l)
	})

	t.Run("Update", func(t *testing.T) {
		err := storage.Update(ctx, updatedEvent)
		require.NoError(t, err)
		l, err := storage.List(ctx, 0)
		require.NoError(t, err)
		require.ElementsMatch(t, []domain.Event{updatedEvent}, l)
	})

	t.Run("Add next", func(t *testing.T) {
		id, err := storage.Add(ctx, event2)
		require.NotEmpty(t, id)
		require.NoError(t, err)
		l, err := storage.List(ctx, 0)
		require.NoError(t, err)
		require.ElementsMatch(t, []domain.Event{updatedEvent, event2}, l)
	})

	t.Run("Delete", func(t *testing.T) {
		err := storage.Remove(ctx, event1.ID)
		require.NoError(t, err)
		l, err := storage.List(ctx, 0)
		require.NoError(t, err)
		require.ElementsMatch(t, []domain.Event{event2}, l)

		err = storage.Remove(ctx, event2.ID)
		require.NoError(t, err)
		l, err = storage.List(ctx, 0)
		require.NoError(t, err)
		require.Empty(t, l)
	})
}

func TestAlreadyExists(t *testing.T) {
	event := domain.Event{ID: "1", Title: "title1", DateStart: time.Now()}
	storage := New()
	ctx := context.Background()

	id, err := storage.Add(ctx, event)
	require.NotEmpty(t, id)
	require.NoError(t, err)

	_, err = storage.Add(ctx, event)
	var eventExistsErr domain.EventExistsError
	require.ErrorAs(t, err, &eventExistsErr)
}

func TestStorageListByDays(t *testing.T) {
	event1 := domain.Event{ID: "1", Title: "title1", DateStart: time.Now()}
	event2 := domain.Event{ID: "2", Title: "title2", DateStart: time.Now().Add(128 * time.Hour)}

	storage := New()
	ctx := context.Background()
	id, err := storage.Add(ctx, event1)
	require.NotEmpty(t, id)
	require.NoError(t, err)

	id, err = storage.Add(ctx, event2)
	require.NotEmpty(t, id)
	require.NoError(t, err)

	l, err := storage.List(ctx, 0)
	require.NoError(t, err)
	require.Len(t, l, 2)

	l, err = storage.List(ctx, 2)
	require.NoError(t, err)
	require.ElementsMatch(t, []domain.Event{event1}, l)
}
