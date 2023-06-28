package iterator

import (
	"context"
	"fmt"
)

type Entity interface {
}

type CursorIterator[T Entity] struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (data []T)
	cursorExtractor func(data []T) (nextCursor int64)
	endChecker      func(ctx context.Context, cursor int64) (shouldEnd bool)
}

func NewCursorIteratorBuilder[T Entity]() *CursorIterator[T] {
	return &CursorIterator[T]{
		endChecker: func(ctx context.Context, id int64) bool {
			return id == 0
		},
	}
}

func (c *CursorIterator[T]) WithInitCursor(id int64) *CursorIterator[T] {
	c.initCursor = id
	return c
}

func (c *CursorIterator[T]) WithDataRetriever(retriever func(ctx context.Context, cursor int64) (data []T)) *CursorIterator[T] {
	c.dataRetriever = retriever
	return c
}

func (c *CursorIterator[T]) WithCursorExtractor(extractor func(data []T) (nextCursor int64)) *CursorIterator[T] {
	c.cursorExtractor = extractor
	return c
}

func (c *CursorIterator[T]) WithEndChecker(checker func(ctx context.Context, cursor int64) (shouldEnd bool)) *CursorIterator[T] {
	c.endChecker = checker
	return c
}

// dataProcessor is a function type that processes an Entity.
// It returns a boolean indicating whether the iteration should end, and any error encountered.
type dataProcessor[T Entity] func(t T) (shouldEnd bool, handlerErr error)

func (c *CursorIterator[T]) Iterate(ctx context.Context, processor dataProcessor[T]) error {
	if c.dataRetriever == nil || c.cursorExtractor == nil {
		return fmt.Errorf("iterator not properly set up")
	}
	cursor := c.initCursor
	if c.endChecker(ctx, cursor) {
		return nil
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		data := c.dataRetriever(ctx, cursor)
		for _, e := range data {
			if shouldEnd, err := processor(e); err != nil {
				return err
			} else if shouldEnd {
				return nil
			}
		}
		cursor = c.cursorExtractor(data)
		if c.endChecker(ctx, cursor) {
			break
		}
	}
	return nil
}
