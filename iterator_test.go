package iterator

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

	dataRetriever := func(ctx context.Context, cursor int64) (data []SimpleEntity) {
		if cursor < int64(len(entities)) {
			return entities[cursor : cursor+1]
		}
		return []SimpleEntity{}
	}

	cursorExtractor := func(data []SimpleEntity) (nextCursor int64) {
		return int64(data[len(data)-1].value)
	}

	endChecker := func(ctx context.Context, cursor int64) (shouldEnd bool) {
		return cursor >= int64(len(entities))
	}

	iterator := NewCursorIteratorBuilder[SimpleEntity]().
		WithInitCursor(0).
		WithDataRetriever(dataRetriever).
		WithCursorExtractor(cursorExtractor).
		WithEndChecker(endChecker)

	err := iterator.Iterate(ctx, func(t SimpleEntity) (shouldEnd bool, handlerErr error) {
		fmt.Println("Processing entity: ", t.value)
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

	iterator := NewCursorIteratorBuilder[SimpleEntity]().
		WithInitCursor(0).
		WithDataRetriever(func(ctx context.Context, cursor int64) (data []SimpleEntity) {
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
		WithCursorExtractor(func(data []SimpleEntity) (nextCursor int64) {
			return int64(data[len(data)-1].value)
		}).
		WithEndChecker(func(ctx context.Context, cursor int64) (shouldEnd bool) {
			return cursor >= int64(len(entities))
		})

	// Count entities to verify all entities are fetched
	count := 0
	err := iterator.Iterate(ctx, func(t SimpleEntity) (shouldEnd bool, handlerErr error) {
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
