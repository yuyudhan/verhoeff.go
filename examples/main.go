// FilePath: examples/main.go

// Example program demonstrating the usage of the verhoeff package.

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/yuyudhan/go-verhoeff"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  generate NUMBER        - Generate a checksum for NUMBER")
		fmt.Println("  validate NUMBER        - Validate NUMBER (with checksum as the last digit)")
		fmt.Println("  validateaadhaar NUMBER - Validate an Aadhaar number")
		fmt.Println("  append NUMBER          - Append a checksum to NUMBER")
		os.Exit(1)
	}

	command := os.Args[1]

	if len(os.Args) < 3 {
		fmt.Println("Error: NUMBER argument is required")
		os.Exit(1)
	}

	number := os.Args[2]

	switch command {
	case "generate":
		runGenerate(number)
	case "validate":
		runValidate(number)
	case "validateaadhaar":
		runValidateAadhaar(number)
	case "append":
		runAppend(number)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runGenerate(number string) {
	checksum, err := verhoeff.Generate(number)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Checksum for %s: %d\n", number, checksum)
}

func runValidate(number string) {
	valid, err := verhoeff.Validate(number)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if valid {
		fmt.Printf("%s is valid\n", number)
	} else {
		fmt.Printf("%s is NOT valid\n", number)

		// Try to suggest the correct number
		if len(number) > 1 {
			base := number[:len(number)-1]
			lastDigit := number[len(number)-1:]
			correctChecksum, err := verhoeff.Generate(base)
			if err == nil {
				lastDigitInt, err := strconv.Atoi(lastDigit)
				if err == nil {
					fmt.Printf("The correct checksum for %s would be %d (you provided %d)\n",
						base, correctChecksum, lastDigitInt)
					fmt.Printf("Correct number would be: %s%d\n", base, correctChecksum)
				}
			}
		}
	}
}

func runValidateAadhaar(number string) {
	valid, err := verhoeff.ValidateAadhaar(number)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if valid {
		fmt.Printf("Aadhaar number %s is valid\n", number)
	} else {
		fmt.Printf("Aadhaar number %s is NOT valid\n", number)
	}
}

func runAppend(number string) {
	result, err := verhoeff.AppendChecksum(number)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s with checksum: %s\n", number, result)

	// Verify it's valid
	valid, _ := verhoeff.Validate(result)
	if !valid {
		fmt.Println("Warning: The generated number failed validation!")
	}
}
