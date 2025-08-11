# Developer Guide

This guide contains information for developers contributing to or working with the verhoeff.go package.

## Project Structure

```
verhoeff.go/
├── verhoeff.go          # Core implementation
├── verhoeff_test.go     # Unit tests
├── integration_test.go  # Integration tests
├── stress_test.go       # Performance & stress tests
├── examples/            # Example usage
│   └── basic/
│       └── main.go
├── go.mod              # Module definition
├── README.md           # User documentation
└── DEVELOPER.md        # This file
```

## Development Setup

```bash
# Clone the repository
git clone https://github.com/yuyudhan/verhoeff.go
cd verhoeff.go

# Install dependencies (none required)
go mod download

# Run tests
go test

# Run benchmarks
go test -bench=.
```

## Testing

### Run All Tests
```bash
go test -v
```

### Test Coverage
```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Benchmarks
```bash
# Run all benchmarks
go test -bench=.

# With memory profiling
go test -bench=. -benchmem

# Specific benchmark
go test -bench=BenchmarkGenerate
```

### Race Detection
```bash
go test -race
```

### Skip Exhaustive Tests
```bash
go test -short
```

## Performance Metrics

Current benchmarks on Apple M3 Pro:

```
BenchmarkGenerate-11           7,933,540    153.0 ns/op    248 B/op    5 allocs/op
BenchmarkValidate-11           7,523,409    159.4 ns/op    248 B/op    5 allocs/op
BenchmarkValidateAadhaar-11    3,566,602    336.3 ns/op    528 B/op   12 allocs/op
BenchmarkAppendChecksum-11     6,402,806    190.8 ns/op    280 B/op    7 allocs/op
```

## Architecture

### Core Algorithm

The implementation uses three lookup tables:
- **Multiplication table (d)**: Based on dihedral group D₅
- **Permutation table (p)**: Position-dependent permutations
- **Inverse table (inv)**: Final checksum mapping

### Type Design

The package provides two API styles:

1. **Type-specific functions** (Recommended):
   - `GenerateFromString`, `GenerateInt`, `GenerateInt64`
   - No runtime type assertions
   - Better performance

2. **Generic functions** (Backward compatible):
   - `Generate`, `Validate`, `AppendChecksum`
   - Accept `interface{}`
   - Runtime type switching

## Contributing

### Code Style
- Follow standard Go conventions
- Run `gofmt` before committing
- Maintain test coverage above 95%
- Add benchmarks for new functions

### Testing Requirements
- Write unit tests for all new functions
- Add integration tests for feature changes
- Include stress tests for performance-critical code
- Ensure all tests pass with `go test -race`

### Commit Messages
Use conventional commit format:
```
feat: add new validation function
fix: correct checksum calculation for edge case
test: add benchmarks for concurrent operations
docs: update API documentation
```

## Implementation Notes

### Memory Optimization
- Pre-allocated slices where possible
- Minimal string concatenation
- Efficient digit extraction for integers

### Thread Safety
- No global mutable state
- All functions are pure
- Safe for concurrent use

### Error Handling
- Lowercase error messages (Go convention)
- Specific error types where appropriate
- Clear error context

## Debugging

### Verbose Test Output
```bash
go test -v -run TestGenerate
```

### CPU Profiling
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### Memory Profiling
```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Release Process

1. Run all tests: `go test -race`
2. Check coverage: `go test -cover`
3. Run benchmarks: `go test -bench=.`
4. Update version in go.mod if needed
5. Tag release: `git tag v1.0.0`
6. Push tags: `git push --tags`

## Support

For issues or questions:
- Open an issue on [GitHub](https://github.com/yuyudhan/verhoeff.go/issues)
- Check existing issues first
- Include Go version and minimal reproduction code