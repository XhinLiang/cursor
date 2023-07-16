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
	dataRetriever   func(ctx context.Context, cursor int64) (data Any)
	cursorExtractor func(data Any) (nextCursor int64)
	endChecker      func(ctx context.Context, cursor int64) (shouldEnd bool)
}

func NewBuilder() *Builder {
	return &Builder{
		endChecker: func(ctx context.Context, id int64) bool {
			return id == 0
		},
	}
}

func (c *Builder) WithInitCursor(id int64) *Builder {
	c.initCursor = id
	return c
}

func (c *Builder) WithDataRetriever(retriever func(ctx context.Context, cursor int64) (data Any)) *Builder {
	c.dataRetriever = retriever
	return c
}

func (c *Builder) WithCursorExtractor(extractor func(data Any) (nextCursor int64)) *Builder {
	c.cursorExtractor = extractor
	return c
}

func (c *Builder) WithEndChecker(checker func(ctx context.Context, cursor int64) (shouldEnd bool)) *Builder {
	c.endChecker = checker
	return c
}

type iterator struct {
	initCursor      int64
	dataRetriever   func(ctx context.Context, cursor int64) (data Any)
	cursorExtractor func(data Any) (nextCursor int64)
	endChecker      func(ctx context.Context, cursor int64) (shouldEnd bool)
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
		endChecker:      c.endChecker,
	}
}

func (c *iterator) Iterate(ctx context.Context, processor DataProcessor) error {
	if c.dataRetriever == nil || c.cursorExtractor == nil {
		return ErrIteratorNotProperlySetUp
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

		if v := reflect.ValueOf(data); v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				processor(v.Index(i).Interface())
			}
		} else {
			return ErrDataIsNotSlice
		}

		cursor = c.cursorExtractor(data)
		if c.endChecker(ctx, cursor) {
			break
		}
	}
	return nil
}
