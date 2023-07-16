package cursor

import (
	"context"
	"fmt"
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

	dataRetriever := func(ctx context.Context, cursor int64) (data Any, err error) {
		if cursor < int64(len(entities)) {
			return entities[cursor : cursor+1], nil
		}
		return []SimpleEntity{}, nil
	}

	cursorExtractor := func(d Any) (shouldEnd bool, nextCursor int64, err error) {
		data := d.([]SimpleEntity)
		nextCursor = int64(data[len(data)-1].value)
		if nextCursor >= int64(len(entities)) {
			return true, 0, nil
		}
		return false, nextCursor, nil
	}

	iterator := NewBuilder().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		Build()

	err := iterator.Iterate(ctx, func(t Any) (shouldEnd bool, handlerErr error) {
		entity := t.(SimpleEntity)
		fmt.Println("Processing entity: ", entity.value)
		return false, nil
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

	iterator := NewBuilder().
		WithInitCursor(0).
		WithDataRetriever(func(ctx context.Context, cursor int64) (data Any, err error) {
			time.Sleep(10 * time.Millisecond) // Simulate network latency
			if cursor < int64(len(entities)) {
				// Fetch 10 items per batch
				end := cursor + 10
				if end > int64(len(entities)) {
					end = int64(len(entities))
				}
				return entities[cursor:end], nil
			}
			return []SimpleEntity{}, nil
		}).
		WithCursorExtractor(func(d Any) (shouldEnd bool, nextCursor int64, err error) {
			data := d.([]SimpleEntity)
			nextCursor = int64(data[len(data)-1].value)
			if nextCursor >= int64(len(entities)) {
				return true, 0, nil
			}
			return false, nextCursor, nil
		}).
		Build()

	// Count entities to verify all entities are fetched
	count := 0
	err := iterator.Iterate(ctx, func(t Any) (shouldEnd bool, handlerErr error) {
		count++
		return false, nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if all entities were fetched
	if count != len(entities) {
		t.Errorf("Expected to process %v entities, but processed %v", len(entities), count)
	}
}
