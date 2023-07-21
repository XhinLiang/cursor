package cursor

import (
	"context"
	"fmt"
	"reflect"
)

var ErrIteratorNotProperlySetUp = fmt.Errorf("iterator not properly set up")
var ErrDataIsNotSlice = fmt.Errorf("list is not a slice")

type Any interface {
}

type Builder struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (list Any, err error)
	cursorExtractor func(list Any, previousCursor int64) (shouldEnd bool, nextCursor int64, err error)
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (c *Builder) WithInitCursor(id int64) *Builder {
	c.initCursor = id
	return c
}

func (c *Builder) WithDataRetriever(retriever func(ctx context.Context, cursor int64) (list Any, err error)) *Builder {
	c.dataRetriever = retriever
	return c
}

func (c *Builder) WithCursorExtractor(extractor func(list Any, previousCursor int64) (shouldEnd bool, nextCursor int64, err error)) *Builder {
	c.cursorExtractor = extractor
	return c
}

type iterator struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (list Any, err error)
	cursorExtractor func(list Any, previousCursor int64) (shouldEnd bool, nextCursor int64, err error)
}

// SingleProcessor is a function type that processes the single entity from the data slice
type SingleProcessor func(single Any) error

// BatchProcessor is a function type that processes the entire data slice
type BatchProcessor func(list Any) error

type Iterator interface {
	Iterate(ctx context.Context, processor SingleProcessor) error
	IterateBatch(ctx context.Context, processor BatchProcessor) error
}

func (c *Builder) Build() Iterator {
	return &iterator{
		initCursor:      c.initCursor,
		dataRetriever:   c.dataRetriever,
		cursorExtractor: c.cursorExtractor,
	}
}

func (c *iterator) Iterate(ctx context.Context, processor SingleProcessor) error {
	return c.IterateBatch(ctx, func(list Any) error {
		if v := reflect.ValueOf(list); v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				if err := processor(v.Index(i).Interface()); err != nil {
					return err
				}
			}
		} else {
			return ErrDataIsNotSlice
		}
		return nil
	})
}

func (c *iterator) IterateBatch(ctx context.Context, processor BatchProcessor) error {
	if c.dataRetriever == nil || c.cursorExtractor == nil {
		return ErrIteratorNotProperlySetUp
	}
	cursor := c.initCursor
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		retrievedList, retrievedErr := c.dataRetriever(ctx, cursor)
		if retrievedErr != nil {
			return retrievedErr
		}

		iterateErr := processor(retrievedList)
		if iterateErr != nil {
			return iterateErr
		}

		shouldEnd, nextCursor, err := c.cursorExtractor(retrievedList, cursor)
		if err != nil {
			return err
		}
		if shouldEnd {
			return nil
		}
		cursor = nextCursor
	}
}
