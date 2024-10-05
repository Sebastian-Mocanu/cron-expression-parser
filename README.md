# Cron Expression Parser

This is a command-line application that parses a cron string and expands each field to show the times at which it will run.

## Prerequisites to build&contribute:

- Go 1.23
- cobra-cli (can be installed with `go install github.com/spf13/cobra-cli@latest`)

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/sebastian-mocanu/cron-expression-parser.git
   cd cron-expression-parser
   ```

2. Build the application:
   ```
   go build -o cron-parser
   ```

## Usage

Run the application with a cron expression as an argument:

```
./cron-parser parse "*/15 0 1,15 * 1-5 /usr/bin/find"
```

The output will be formatted as a table with the field name taking the first 14 columns and the times as a space-separated list following it.

## Example

Input:

```
./cron-parser parse "*/15 0 1,15 * 1-5 /usr/bin/find"
```

Output:

```
minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```

## Development

To add new commands or modify existing ones, use the cobra-cli:

```
cobra-cli add [command-name]
```

Then edit the generated file in the `cmd` directory.

## Testing

To run the tests:

```
go test ./... (add -v flag to see output)
```
