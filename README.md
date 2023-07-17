# Cursor 

## Description

The `cursor` package provides an easy-to-use framework for building and managing iterators in Go. It features a Builder Pattern implementation for creating iterators, and also abstracts away the underlying logic of iterating over a set of data, given the data retriever and cursor extractor functions.

## Usage

### Building the Iterator

You can build an iterator with a Builder. First, you need to initialize a Builder with `NewBuilder()`. After that, you can set up the required functions like data retriever, cursor extractor, and initialize the cursor:

```go
iterator := cursor.NewBuilder().
    WithInitCursor(0).
    WithDataRetriever(myDataRetriever).
    WithCursorExtractor(myCursorExtractor).
    Build()
```

### Data Retriever

A data retriever is a function that retrieves data based on the current cursor position. It should be a function of type:

```go
func(ctx context.Context, cursor int64) (dataSlice Any, err error)
```

This function takes in a context and a cursor value, and returns a slice of data.

### Cursor Extractor

A cursor extractor is a function that, given a data slice, decides whether the iteration should end and what the next cursor should be. It should be a function of type:

```go
func(dataSlice Any) (shouldEnd bool, nextCursor int64, err error)
```

This function takes in a data slice, and returns a boolean indicating if iteration should end, the next cursor value, and an error, if any occurred.

### Iterating Over Data

Once you have a built iterator, you can iterate over your data with the `Iterate` function, providing a data processor:

```go
err := iterator.Iterate(ctx, myDataProcessor)
```

The `DataProcessor` is a function that processes a single entity from the data slice. It should be of type:

```go
func(t Any) error
```

## Error Handling

This package defines two errors:
- `ErrIteratorNotProperlySetUp`: Returned when the iterator is not correctly set up (data retriever or cursor extractor is missing).
- `ErrDataIsNotSlice`: Returned when the data retrieved is not a slice.

## Future Work

Batch iterating functionality is currently being developed.

## Contributing

If you want to contribute to this project, please feel free to fork the repository, create a feature branch, and open a pull request. If you encounter problems or have suggestions, please open an issue.
