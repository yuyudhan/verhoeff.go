# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a Go implementation of the Verhoeff algorithm for checksum calculation and validation. The Verhoeff algorithm can detect all single-digit errors and most transposition errors, making it particularly useful for validating identification numbers like India's Aadhaar numbers.

## Common Development Commands

### Building
```bash
go build -v ./...
```

### Testing
```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run benchmarks
go test -bench=. -benchmem

# Run a specific test
go test -v -run TestValidateAadhaar
```

### Running the Example
```bash
# Build and run the example program
go run examples/main.go generate 12345
go run examples/main.go validate 123450
go run examples/main.go validateaadhaar 496858245152
go run examples/main.go append 12345
```

## Architecture & Key Components

### Core Package Structure
The library is a single-package implementation with clear separation of concerns:

- **verhoeff.go**: Core algorithm implementation with lookup tables (d, p, inv) and all public API functions
- **verhoeff_test.go**: Comprehensive test suite including unit tests and benchmarks
- **examples/main.go**: CLI tool demonstrating all library functions

### Key Functions
- `Generate(input)`: Calculates checksum digit for input
- `Validate(input)`: Validates a number with its checksum
- `ValidateAadhaar(aadhaarStr)`: Validates 12-digit Aadhaar numbers
- `AppendChecksum(input)`: Generates and appends checksum to input
- `ConvertToDigits(input)`: Converts string/int to digit slice
- `InvertArray(input)`: Reverses digit array for algorithm processing

### Algorithm Implementation
The Verhoeff algorithm uses three lookup tables:
- **d table**: 10x10 multiplication table for the dihedral group D5
- **p table**: 8x10 permutation table
- **inv table**: Inverse values for checksum calculation

The algorithm processes digits in reverse order, applying permutations and group operations to detect both single-digit errors and transposition errors.

## Testing Approach

The test suite covers:
- Basic functionality for all public functions
- Edge cases (empty strings, non-digit characters)
- Aadhaar-specific validation (12-digit requirement)
- Performance benchmarks for optimization
- Input type flexibility (string, int)

## Dependencies

This is a zero-dependency library - it only uses Go's standard library (errors, strconv, unicode).

## Module Information

- Module: `github.com/yuyudhan/go-verhoeff`
- Go version: 1.21+
- License: MIT