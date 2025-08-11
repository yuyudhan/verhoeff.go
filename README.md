# verhoeff.go

Go implementation of the Verhoeff checksum algorithm for error detection in numerical identifiers.

[![Go Reference](https://pkg.go.dev/badge/github.com/yuyudhan/verhoeff.go.svg)](https://pkg.go.dev/github.com/yuyudhan/verhoeff.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuyudhan/verhoeff.go)](https://goreportcard.com/report/github.com/yuyudhan/verhoeff.go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- ✅ Detects **100% of single-digit errors** and **adjacent transpositions**
- ⚡ High performance: ~150ns per operation
- 🔒 Type-safe APIs for strings, integers, and slices
- 🇮🇳 Built-in Aadhaar number validation
- 🧵 Thread-safe and zero dependencies

## Installation

```bash
go get github.com/yuyudhan/verhoeff.go
```

## Quick Start

```go
import verhoeff "github.com/yuyudhan/verhoeff.go"

// Generate checksum
checksum, _ := verhoeff.GenerateFromString("12345")  // Returns: 1

// Validate number
valid, _ := verhoeff.ValidateString("123451")  // Returns: true

// Append checksum
result, _ := verhoeff.AppendChecksumString("12345")  // Returns: "123451"

// Validate Aadhaar
valid, _ := verhoeff.ValidateAadhaar("123456789010")  // Returns: true
```

## API Reference

### String Operations
```go
GenerateFromString(s string) (int, error)
ValidateString(s string) (bool, error)
AppendChecksumString(s string) (string, error)
```

### Integer Operations
```go
GenerateInt(n int) int
ValidateInt(n int) bool
AppendChecksumInt(n int) string
```

### Aadhaar Validation
```go
ValidateAadhaar(aadhaarStr string) (bool, error)
```

### Generic Functions
```go
Generate(input interface{}) (int, error)     // Supports string, int, int64, []int
Validate(input interface{}) (bool, error)
AppendChecksum(input interface{}) (string, error)
```

## Examples

```bash
# Run example
go run examples/basic/main.go

# With arguments
go run examples/basic/main.go generate 12345
go run examples/basic/main.go validate 123451
go run examples/basic/main.go validateaadhaar 123456789010
```

## Development

See [DEVELOPER.md](DEVELOPER.md) for development setup, testing, and contributing guidelines.

## License

MIT

## Credits

- Algorithm by J. Verhoeff
- Based on [Node.js implementation](https://github.com/yuyudhan/verhoeff)