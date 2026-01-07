# Bank Statement CSV Processing

## Usage

```sh
go run main.go -help

Usage of main:

A simple Bank Statement CSV processing with worker pool.

Options:
  -csv string
        Comma separated CSV path
  -pool int
        Number of worker pool (default 3)
```

## Install Dependencies (if needed)

```sh
go mod tidy
```

## Run the Program

```sh
go run main.go -pool 2 -csv ./example.csv,./example_invalid.csv,./example_invalid_partial.csv
```

## Run Test

```sh
go test -v -race ./...
```

## Build the Program

```sh
go build -o <executable name> main.go
```