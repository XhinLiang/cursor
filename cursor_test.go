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

	dataRetriever := func(ctx context.Context, cursor int64) (data Any) {
		if cursor < int64(len(entities)) {
			return entities[cursor : cursor+1]
		}
		return []SimpleEntity{}
	}

	cursorExtractor := func(d Any) (nextCursor int64) {
		data := d.([]SimpleEntity)
		return int64(data[len(data)-1].value)
	}

	endChecker := func(ctx context.Context, cursor int64) (shouldEnd bool) {
		return cursor >= int64(len(entities))
	}

	iterator := NewBuilder().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		WithEndChecker(endChecker).
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
		WithDataRetriever(func(ctx context.Context, cursor int64) (data Any) {
			time.Sleep(10 * time.Millisecond) // Simulate network latency
			if cursor < int64(len(entities)) {
				// Fetch 10 items per batch
				end := cursor + 10
				if end > int64(len(entities)) {
					end = int64(len(entities))
				}
				return entities[cursor:end]
			}
			return []SimpleEntity{}
		}).
		WithCursorExtractor(func(d Any) (nextCursor int64) {
			data := d.([]SimpleEntity)
			return int64(data[len(data)-1].value)
		}).
		WithEndChecker(func(ctx context.Context, cursor int64) (shouldEnd bool) {
			return cursor >= int64(len(entities))
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
