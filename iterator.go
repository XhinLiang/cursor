package iterator

import (
	"context"
	"fmt"
	"reflect"
)

type Any interface {
}

type CursorIterator struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (data Any)
	cursorExtractor func(data Any) (nextCursor int64)
	endChecker      func(ctx context.Context, cursor int64) (shouldEnd bool)
}

func NewCursorIteratorBuilder() *CursorIterator {
	return &CursorIterator{
		endChecker: func(ctx context.Context, id int64) bool {
			return id == 0
		},
	}
}

func (c *CursorIterator) WithInitCursor(id int64) *CursorIterator {
	c.initCursor = id
	return c
}

func (c *CursorIterator) WithDataRetriever(retriever func(ctx context.Context, cursor int64) (data Any)) *CursorIterator {
	c.dataRetriever = retriever
	return c
}

func (c *CursorIterator) WithCursorExtractor(extractor func(data Any) (nextCursor int64)) *CursorIterator {
	c.cursorExtractor = extractor
	return c
}

func (c *CursorIterator) WithEndChecker(checker func(ctx context.Context, cursor int64) (shouldEnd bool)) *CursorIterator {
	c.endChecker = checker
	return c
}

// dataProcessor is a function type that processes an Any.
// It returns a boolean indicating whether the iteration should end, and any error encountered.
type dataProcessor func(t Any) (shouldEnd bool, handlerErr error)

func (c *CursorIterator) Iterate(ctx context.Context, processor dataProcessor) error {
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

		v := reflect.ValueOf(data)

		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				// Convert reflect.Value to interface{} and pass it to processor
				processor(v.Index(i).Interface())
			}
		} else {
			return fmt.Errorf("data is not a slice")
		}

		cursor = c.cursorExtractor(data)
		if c.endChecker(ctx, cursor) {
			break
		}
	}
	return nil
}
