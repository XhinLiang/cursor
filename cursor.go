package cursor

import (
	"context"
	"fmt"
)

var ErrIteratorNotProperlySetUp = fmt.Errorf("iterator not properly set up")
var ErrDataIsNotSlice = fmt.Errorf("data is not a slice")

type Builder[T any, U comparable] struct {
	initCursor      U
	dataRetriever   func(ctx context.Context, cursor U) (data []T, err error)
	cursorExtractor func(data []T) (shouldEnd bool, nextCursor U, err error)
}

func NewBuilder[T any, U comparable]() *Builder[T, U] {
	return &Builder[T, U]{}
}

func (c *Builder[T, U]) WithInitCursor(id U) *Builder[T, U] {
	c.initCursor = id
	return c
}

func (c *Builder[T, U]) WithDataRetriever(retriever func(ctx context.Context, cursor U) (data []T, err error)) *Builder[T, U] {
	c.dataRetriever = retriever
	return c
}

func (c *Builder[T, U]) WithCursorExtractor(extractor func(data []T) (shouldEnd bool, nextCursor U, err error)) *Builder[T, U] {
	c.cursorExtractor = extractor
	return c
}

type iterator[T any, U comparable] struct {
	initCursor      U
	dataRetriever   func(ctx context.Context, cursor U) (data []T, err error)
	cursorExtractor func(data []T) (shouldEnd bool, nextCursor U, err error)
}

// dataProcessor is a function type that processes the single entity from the data slice
type DataProcessor[T any] func(t T) error

type Iterator[T any, U comparable] interface {
	Iterate(ctx context.Context, processor DataProcessor[T]) error
}

func (c *Builder[T, U]) Build() Iterator[T, U] {
	return &iterator[T, U]{
		initCursor:      c.initCursor,
		dataRetriever:   c.dataRetriever,
		cursorExtractor: c.cursorExtractor,
	}
}

func (c *iterator[T, U]) Iterate(ctx context.Context, processor DataProcessor[T]) error {
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

		// TODO provide batch iterating
		for _, data := range retrievedData {
			err := processor(data)
			if err != nil {
				return err
			}
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
