<!--FilePath: README.md-->

# go-verhoeff

A Go implementation of the Verhoeff algorithm, a checksum formula for error detection.

## About the Verhoeff Algorithm

The Verhoeff algorithm is a checksum formula for error detection developed by Dutch mathematician J. Verhoeff. Unlike many other checksum formulas, it can detect all single-digit errors and most transposition errors (when adjacent digits are swapped).

For more information, see the [Wikipedia article](https://en.wikipedia.org/wiki/Verhoeff_algorithm).

## Aadhaar Numbers & Verhoeff

Did you know Aadhaar numbers (India's unique identification numbers) have the last digit as a Verhoeff checksum? The idea behind this is to quickly identify typing/data-entry errors on the entry machine.

This implementation provides a specific function for validating Aadhaar numbers using the Verhoeff algorithm, making it simple to verify their integrity.

## Installation

```bash
go get github.com/yuyudhan/go-verhoeff
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"github.com/yuyudhan/go-verhoeff"
)

func main() {
	// Generate a checksum digit for a number
	checksum, err := verhoeff.Generate("12345")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Checksum for 12345: %d\n", checksum)

	// Generate and append the checksum to a number
	numberWithChecksum, err := verhoeff.AppendChecksum("12345")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Number with checksum: %s\n", numberWithChecksum)

	// Validate a number with a checksum
	valid, err := verhoeff.Validate(numberWithChecksum)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Is valid: %t\n", valid)

	// Validate an Indian Aadhaar number
	validAadhaar, err := verhoeff.ValidateAadhaar("496858245152") // Example Aadhaar number
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Is valid Aadhaar: %t\n", validAadhaar)
}
```

### Validating Aadhaar Numbers

```go
// Check if an Aadhaar number is valid
isValid, err := verhoeff.ValidateAadhaar("496858245152")
if err != nil {
    fmt.Println("Error:", err)
} else if isValid {
    fmt.Println("The Aadhaar number is valid")
} else {
    fmt.Println("The Aadhaar number is invalid")
}
```

### API

#### `Generate(input interface{}) (int, error)`

Calculates the Verhoeff checksum digit for a given input. The input can be a string of digits or an integer.

```go
// Using a string
checksum, err := verhoeff.Generate("12345")

// Using an integer
checksum, err := verhoeff.Generate(12345)

// Using a slice
checksum, err := verhoeff.Generate([]int{1, 2, 3, 4, 5})
```

#### `GenerateString(input interface{}) (string, error)`

Same as `Generate`, but returns the checksum as a string.

#### `Validate(input interface{}) (bool, error)`

Checks if a number with its checksum digit is valid. The input can be a string of digits or an integer.

#### `ValidateAadhaar(aadhaarStr string) (bool, error)`

Checks if an Aadhaar number (Indian identification number) is valid. Aadhaar numbers must be exactly 12 digits.

#### `AppendChecksum(input interface{}) (string, error)`

Adds the calculated checksum digit to the end of input and returns the result.

#### `ConvertToDigits(input interface{}) ([]int, error)`

Converts a string or integer to a slice of digits.

#### `InvertArray(input interface{}) ([]int, error)`

Converts input to a slice of digits and reverses it.

## Running Tests

```bash
go test
```

For benchmarks:

```bash
go test -bench=.
```

## Original Implementation

This is a Go port of the [node-verhoeff](https://www.npmjs.com/package/node-verhoeff) package, which was originally created by Sergey Petushkov in 2014 and later implemented for Node.js by [@yuyudhan](https://github.com/yuyudhan).

The Node.js version can be installed via npm:

```bash
npm i node-verhoeff --save
```

## License

MIT License - See LICENSE file for details.

