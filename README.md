# Cursor Iterator in Go

This project provides an implementation of a cursor-based iterator in Go. It offers an efficient and easy-to-use way of iterating over a large dataset chunk by chunk. 

## Getting Started

To use this package, you need to import it into your Go application:

```go
import "github.com/xhinliang/iterator"
```

## Features

- **Chainable configuration methods:** The package provides several "With" methods that can be chained together to build a `CursorIterator` instance.
- **Customizable data retrieval and processing:** You can provide your own methods for data retrieval, cursor extraction, end condition checking, and data processing.
- **Context support:** The iteration process respects the cancellation or timeout from the provided context.

## Usage

Here's a simple example of how to use this iterator:

```go
// Create a new iterator
iterator := iterator.NewCursorIteratorBuilder[Entity]().
	WithInitCursor(0).
	WithDataRetriever(myDataRetriever).
	WithCursorExtractor(myCursorExtractor).
	WithEndChecker(myEndChecker)

// Use the iterator
err := iterator.Iterate(context.Background(), myDataProcessor)
if err != nil {
	log.Fatalf("Error during iteration: %v", err)
}
```

In this example, `myDataRetriever`, `myCursorExtractor`, `myEndChecker`, and `myDataProcessor` are functions you define to retrieve, extract cursor from, check end condition for, and process your data, respectively.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the LICENSE file for details.