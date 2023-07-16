package cursor

import (
	"context"
	"fmt"
	"reflect"
)

var ErrIteratorNotProperlySetUp = fmt.Errorf("iterator not properly set up")
var ErrDataIsNotSlice = fmt.Errorf("data is not a slice")

type Any interface {
}

type Builder struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (data Any, err error)
	cursorExtractor func(data Any) (shouldEnd bool, nextCursor int64, err error)
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (c *Builder) WithInitCursor(id int64) *Builder {
	c.initCursor = id
	return c
}

func (c *Builder) WithDataRetriever(retriever func(context.Context, int64) (Any, error)) *Builder {
	c.dataRetriever = retriever
	return c
}

func (c *Builder) WithCursorExtractor(extractor func(Any) (bool, int64, error)) *Builder {
	c.cursorExtractor = extractor
	return c
}

type iterator struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (data Any, err error)
	cursorExtractor func(data Any) (shouldEnd bool, nextCursor int64, err error)
}

// dataProcessor is a function type that processes an Any.
// It returns a boolean indicating whether the iteration should end, and any error encountered.
type DataProcessor func(t Any) (shouldEnd bool, handlerErr error)

type Iterator interface {
	Iterate(ctx context.Context, processor DataProcessor) error
}

func (c *Builder) Build() Iterator {
	return &iterator{
		initCursor:      c.initCursor,
		dataRetriever:   c.dataRetriever,
		cursorExtractor: c.cursorExtractor,
	}
}

func (c *iterator) Iterate(ctx context.Context, processor DataProcessor) error {
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

		retrievedData, retrievedErr := c.dataRetriever(ctx, cursor)
		if retrievedErr != nil {
			return retrievedErr
		}

		if v := reflect.ValueOf(retrievedData); v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				processor(v.Index(i).Interface())
			}
		} else {
			return ErrDataIsNotSlice
		}

		shouldEnd, nextCursor, err := c.cursorExtractor(retrievedData)
		if err != nil {
			return err
		}
		if shouldEnd {
			return nil
		}
		cursor = nextCursor
	}
}
