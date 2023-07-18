package cursor

import (
	"context"
	"errors"
	"testing"
	"time"
)

type SimpleEntity struct {
	value int
}

func TestCursorIterator(t *testing.T) {
	ctx := context.Background()

	entities := []SimpleEntity{
		{value: 1},
		{value: 2},
		{value: 3},
		{value: 4},
		{value: 5},
	}

	dataRetriever := func(ctx context.Context, cursor int) (data []SimpleEntity, err error) {
		if cursor < len(entities) {
			return entities[cursor : cursor+1], nil
		}
		return []SimpleEntity{}, nil
	}

	cursorExtractor := func(d []SimpleEntity) (shouldEnd bool, nextCursor int, err error) {
		nextCursor = d[len(d)-1].value
		if nextCursor >= len(entities) {
			return true, 0, nil
		}
		return false, nextCursor, nil
	}

	iterator := NewBuilder[SimpleEntity, int]().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		Build()

	err := iterator.Iterate(ctx, func(e SimpleEntity) error {
		t.Log("Processing entity: ", e.value)
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLargeCursorIterator(t *testing.T) {
	ctx := context.Background()

	// Create 2000 entities
	entities := make([]SimpleEntity, 2000)
	for i := 0; i < 2000; i++ {
		entities[i] = SimpleEntity{value: i + 1}
	}

	dataRetriever := func(ctx context.Context, cursor int) (data []SimpleEntity, err error) {
		time.Sleep(10 * time.Millisecond) // Simulate network latency
		if cursor < len(entities) {
			// Fetch 10 items per batch
			end := cursor + 10
			if end > len(entities) {
				end = len(entities)
			}
			return entities[cursor:end], nil
		}
		return []SimpleEntity{}, nil
	}

	cursorExtractor := func(d []SimpleEntity) (shouldEnd bool, nextCursor int, err error) {
		nextCursor = d[len(d)-1].value
		if nextCursor >= len(entities) {
			return true, 0, nil
		}
		return false, nextCursor, nil
	}

	iterator := NewBuilder[SimpleEntity, int]().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		Build()

	// Count entities to verify all entities are fetched
	count := 0
	err := iterator.Iterate(ctx, func(t SimpleEntity) error {
		count++
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if all entities were fetched
	if count != len(entities) {
		t.Errorf("Expected to process %v entities, but processed %v", len(entities), count)
	}
}

func TestCursorIteratorErrorOnDataRetriever(t *testing.T) {
	ctx := context.Background()

	entities := []SimpleEntity{
		{value: 1},
		{value: 2},
		{value: 3},
		{value: 4},
		{value: 5},
	}

	dataRetriever := func(ctx context.Context, cursor int) (data []SimpleEntity, err error) {
		if cursor < len(entities) {
			return entities[cursor : cursor+1], nil
		}
		// Simulate an error in the data retriever function
		return nil, errors.New("data retriever error")
	}

	cursorExtractor := func(d []SimpleEntity) (shouldEnd bool, nextCursor int, err error) {
		nextCursor = d[len(d)-1].value
		return false, nextCursor, nil
	}

	iterator := NewBuilder[SimpleEntity, int]().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		Build()

	err := iterator.Iterate(ctx, func(e SimpleEntity) error {
		t.Log("Processing entity: ", e.value)
		return nil
	})

	if err == nil || err.Error() != "data retriever error" {
		t.Errorf("Expected error from data retriever, got: %v", err)
	}
}

func TestCursorIteratorErrorOnCursorExtractor(t *testing.T) {
	ctx := context.Background()

	entities := []SimpleEntity{
		{value: 1},
		{value: 2},
		{value: 3},
		{value: 4},
		{value: 5},
	}

	dataRetriever := func(ctx context.Context, cursor int) (data []SimpleEntity, err error) {
		if cursor < len(entities) {
			return entities[cursor : cursor+1], nil
		}
		return []SimpleEntity{}, nil
	}

	// Simulate an error in the cursor extractor function
	cursorExtractor := func(d []SimpleEntity) (shouldEnd bool, nextCursor int, err error) {
		return false, 0, errors.New("cursor extractor error")
	}

	iterator := NewBuilder[SimpleEntity, int]().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		Build()

	err := iterator.Iterate(ctx, func(e SimpleEntity) error {
		t.Log("Processing entity: ", e.value)
		return nil
	})

	if err == nil || err.Error() != "cursor extractor error" {
		t.Errorf("Expected error from cursor extractor, got: %v", err)
	}
}

func TestCursorIteratorCanceledContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	entities := []SimpleEntity{
		{value: 1},
		{value: 2},
		{value: 3},
		{value: 4},
		{value: 5},
	}

	dataRetriever := func(ctx context.Context, cursor int) (data []SimpleEntity, err error) {
		time.Sleep(2 * time.Second) // Simulate a delay
		if cursor < len(entities) {
			return entities[cursor : cursor+1], nil
		}
		return []SimpleEntity{}, nil
	}

	cursorExtractor := func(d []SimpleEntity) (shouldEnd bool, nextCursor int, err error) {
		nextCursor = d[len(d)-1].value
		if nextCursor >= len(entities) {
			return true, 0, nil
		}
		return false, nextCursor, nil
	}

	iterator := NewBuilder[SimpleEntity, int]().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		Build()

	err := iterator.Iterate(ctx, func(e SimpleEntity) error {
		t.Log("Processing entity: ", e.value)
		return nil
	})

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context deadline exceeded error, got: %v", err)
	}
}
