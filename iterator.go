package iterator

type Entity interface {
}

type CursorIterator[T Entity] struct {
	initCursor      int64
	dataRetriever   func(cursor int64) (data []T)
	cursorExtractor func(data []T) (nextCursor int64)
	endChecker      func(int64) (isEnd bool)
}

func NewCursorIteratorExBuilder[T Entity]() *CursorIterator[T] {
	return &CursorIterator[T]{
		endChecker: func(id int64) bool {
			return id == 0
		},
	}
}

func (c *CursorIterator[T]) WithInitCursor(id int64) *CursorIterator[T] {
	c.initCursor = id
	return c
}

func (c *CursorIterator[T]) WithDataRetriever(retriever func(int64) []T) *CursorIterator[T] {
	c.dataRetriever = retriever
	return c
}

func (c *CursorIterator[T]) WithCursorExtractor(extractor func([]T) int64) *CursorIterator[T] {
	c.cursorExtractor = extractor
	return c
}

func (c *CursorIterator[T]) WithEndChecker(checker func(int64) bool) *CursorIterator[T] {
	c.endChecker = checker
	return c
}

func (c *CursorIterator[T]) Iterate(handler func(t T) (shouldEnd bool, handlerErr error)) error {
	cursor := c.initCursor
	if c.endChecker(cursor) {
		return nil
	}
	for {
		data := c.dataRetriever(cursor)
		for _, e := range data {
			if shouldEnd, err := handler(e); err != nil {
				return err
			} else if shouldEnd {
				return nil
			}
		}
		cursor = c.cursorExtractor(data)
		if c.endChecker(cursor) {
			break
		}
	}
	return nil
}
