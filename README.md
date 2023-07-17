# Cursor

[![Go Report Card](https://goreportcard.com/badge/github.com/XhinLiang/cursor)](https://goreportcard.com/report/github.com/XhinLiang/cursor)
[![GoDoc](https://godoc.org/github.com/XhinLiang/cursor?status.svg)](https://godoc.org/github.com/XhinLiang/cursor)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Cursor is a library for Go that provides a builder pattern-based approach for handling cursor-based iteration over data. This abstraction layer simplifies the process of implementing cursor-based iteration and reduces boilerplate code, providing flexibility and reusability in code.

## Installation

To start using cursor, install Go and run `go get`:

```sh
$ go get -u github.com/XhinLiang/cursor
```

## Examples

To create an iterator with the Cursor library, follow these steps:

1. Create a builder
2. Set initial cursor
3. Define a data retriever function
4. Define a cursor extractor function
5. Build the iterator

```go
builder := cursor.NewBuilder().
    WithInitCursor(initCursor).
    WithDataRetriever(myDataRetrieverFunction).
    WithCursorExtractor(myCursorExtractorFunction).
    Build()
```

Then, you can use the iterator in your code:

```go
err := iterator.Iterate(context.Background(), myDataProcessorFunction)
if err != nil {
  log.Fatal(err)
}
```

Note that `myDataRetrieverFunction`, `myCursorExtractorFunction`, and `myDataProcessorFunction` need to be defined according to your specific use case.

## Error Handling

This library has defined some error types that can be returned:

- `ErrIteratorNotProperlySetUp`: Returned when data retriever or cursor extractor are not properly set up.
- `ErrDataIsNotSlice`: Returned when the data returned by the data retriever is not a slice.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://opensource.org/licenses/MIT)