// FilePath: examples/basic/main.go

// Basic example demonstrating the usage of the verhoeff package.
package main

import (
    "fmt"
    "log"

    verhoeff "github.com/yuyudhan/verhoeff.go"
)

func main() {
    // Example 1: Generate a checksum digit
    fmt.Println("=== Generating Checksum ===")
    number := "12345"
    checksum, err := verhoeff.Generate(number)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Checksum for %s: %d\n", number, checksum)
    fmt.Printf("Complete number: %s%d\n\n", number, checksum)

    // Example 2: Validate a number with checksum
    fmt.Println("=== Validating Numbers ===")
    validNumber := "123450"
    isValid, err := verhoeff.Validate(validNumber)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s is valid: %t\n", validNumber, isValid)

    invalidNumber := "123451"
    isValid, err = verhoeff.Validate(invalidNumber)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s is valid: %t\n\n", invalidNumber, isValid)

    // Example 3: Append checksum to a number
    fmt.Println("=== Appending Checksum ===")
    original := "987654"
    withChecksum, err := verhoeff.AppendChecksum(original)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Original: %s\n", original)
    fmt.Printf("With checksum: %s\n\n", withChecksum)

    // Example 4: Working with integers
    fmt.Println("=== Working with Integers ===")
    intNumber := 54321
    intChecksum, err := verhoeff.Generate(intNumber)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Checksum for %d: %d\n", intNumber, intChecksum)

    // Example 5: Validate Aadhaar number
    fmt.Println("\n=== Validating Aadhaar Number ===")
    aadhaar := "234567890124"
    isValidAadhaar, err := verhoeff.ValidateAadhaar(aadhaar)
    if err != nil {
        fmt.Printf("Error validating Aadhaar: %v\n", err)
    } else {
        fmt.Printf("Aadhaar %s is valid: %t\n", aadhaar, isValidAadhaar)
    }
}