// FilePath: verhoeff.go

// Package verhoeff implements the Verhoeff algorithm for checksum
// calculation and validation. The Verhoeff algorithm is a checksum formula
// for error detection developed by Dutch mathematician J. Verhoeff.
// It can detect all single-digit errors and most transposition errors.
// For more information: https://en.wikipedia.org/wiki/Verhoeff_algorithm
package verhoeff

import (
    "errors"
    "fmt"
    "strconv"
    "unicode"
)

// Lookup tables for the Verhoeff algorithm
var (
    // multiplication table d
    d = [][]int{
        {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
        {1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
        {2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
        {3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
        {4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
        {5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
        {6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
        {7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
        {8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
        {9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
    }

    // permutation table p
    p = [][]int{
        {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
        {1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
        {5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
        {8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
        {9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
        {4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
        {2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
        {7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
    }

    // inverse table inv
    inv = []int{0, 4, 3, 2, 1, 5, 6, 7, 8, 9}
)

// stringToDigits converts a string to a slice of digits.
// It returns an error if the string contains non-digit characters.
func stringToDigits(s string) ([]int, error) {
    if s == "" {
        return []int{}, nil
    }
    
    digits := make([]int, 0, len(s))
    for _, char := range s {
        if !unicode.IsDigit(char) {
            return nil, errors.New("input contains non-digit characters")
        }
        digit, _ := strconv.Atoi(string(char))
        digits = append(digits, digit)
    }
    return digits, nil
}

// intToDigits converts an integer to a slice of digits.
func intToDigits(n int) []int {
    if n == 0 {
        return []int{0}
    }
    
    if n < 0 {
        n = -n // Handle negative numbers by taking absolute value
    }
    
    // Count digits
    temp := n
    count := 0
    for temp > 0 {
        count++
        temp /= 10
    }
    
    // Extract digits in reverse order then reverse
    digits := make([]int, count)
    for i := count - 1; i >= 0; i-- {
        digits[i] = n % 10
        n /= 10
    }
    
    return digits
}

// int64ToDigits converts an int64 to a slice of digits.
func int64ToDigits(n int64) []int {
    if n == 0 {
        return []int{0}
    }
    
    if n < 0 {
        n = -n // Handle negative numbers by taking absolute value
    }
    
    // Count digits
    temp := n
    count := 0
    for temp > 0 {
        count++
        temp /= 10
    }
    
    // Extract digits in reverse order then reverse
    digits := make([]int, count)
    for i := count - 1; i >= 0; i-- {
        digits[i] = int(n % 10)
        n /= 10
    }
    
    return digits
}

// sliceToDigits validates a slice of integers as digits.
func sliceToDigits(slice []int) ([]int, error) {
    result := make([]int, len(slice))
    for i, digit := range slice {
        if digit < 0 || digit > 9 {
            return nil, errors.New("input contains invalid digit")
        }
        result[i] = digit
    }
    return result, nil
}

// reverseDigits reverses a slice of digits in place.
func reverseDigits(digits []int) {
    for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
        digits[i], digits[j] = digits[j], digits[i]
    }
}

// calculateChecksum calculates the Verhoeff checksum for a slice of digits.
func calculateChecksum(digits []int) int {
    // Work with a copy to avoid modifying the input
    reversed := make([]int, len(digits))
    copy(reversed, digits)
    reverseDigits(reversed)
    
    c := 0
    for i, digit := range reversed {
        c = d[c][p[(i+1)%8][digit]]
    }
    
    return inv[c]
}

// validateChecksum validates a number with its checksum digit.
func validateChecksum(digits []int) bool {
    if len(digits) == 0 {
        return false
    }
    
    // Work with a copy to avoid modifying the input
    reversed := make([]int, len(digits))
    copy(reversed, digits)
    reverseDigits(reversed)
    
    c := 0
    for i, digit := range reversed {
        c = d[c][p[i%8][digit]]
    }
    
    return c == 0
}

// GenerateFromString calculates the Verhoeff checksum digit for a string of digits.
func GenerateFromString(s string) (int, error) {
    digits, err := stringToDigits(s)
    if err != nil {
        return -1, err
    }
    return calculateChecksum(digits), nil
}

// GenerateInt calculates the Verhoeff checksum digit for an integer.
func GenerateInt(n int) int {
    digits := intToDigits(n)
    return calculateChecksum(digits)
}

// GenerateInt64 calculates the Verhoeff checksum digit for an int64.
func GenerateInt64(n int64) int {
    digits := int64ToDigits(n)
    return calculateChecksum(digits)
}

// GenerateSlice calculates the Verhoeff checksum digit for a slice of digits.
func GenerateSlice(digits []int) (int, error) {
    validDigits, err := sliceToDigits(digits)
    if err != nil {
        return -1, err
    }
    return calculateChecksum(validDigits), nil
}

// Generate calculates the Verhoeff checksum digit for various input types.
// Supported types: string, int, int64, []int
// This function is kept for backward compatibility but using the type-specific
// functions (GenerateFromString, GenerateInt, etc.) is recommended.
func Generate(input interface{}) (int, error) {
    switch v := input.(type) {
    case string:
        return GenerateFromString(v)
    case int:
        return GenerateInt(v), nil
    case int64:
        return GenerateInt64(v), nil
    case []int:
        return GenerateSlice(v)
    default:
        return -1, fmt.Errorf("unsupported input type: %T", input)
    }
}

// ValidateString checks if a string number with its checksum digit is valid.
func ValidateString(s string) (bool, error) {
    digits, err := stringToDigits(s)
    if err != nil {
        return false, err
    }
    if len(digits) == 0 {
        return false, errors.New("empty input")
    }
    return validateChecksum(digits), nil
}

// ValidateInt checks if an integer with its checksum digit is valid.
func ValidateInt(n int) bool {
    digits := intToDigits(n)
    return validateChecksum(digits)
}

// ValidateInt64 checks if an int64 with its checksum digit is valid.
func ValidateInt64(n int64) bool {
    digits := int64ToDigits(n)
    return validateChecksum(digits)
}

// ValidateSlice checks if a slice of digits with its checksum is valid.
func ValidateSlice(digits []int) (bool, error) {
    validDigits, err := sliceToDigits(digits)
    if err != nil {
        return false, err
    }
    if len(validDigits) == 0 {
        return false, errors.New("empty input")
    }
    return validateChecksum(validDigits), nil
}

// Validate checks if a number with its checksum digit is valid.
// Supported types: string, int, int64, []int
// This function is kept for backward compatibility but using the type-specific
// functions (ValidateString, ValidateInt, etc.) is recommended.
func Validate(input interface{}) (bool, error) {
    switch v := input.(type) {
    case string:
        return ValidateString(v)
    case int:
        return ValidateInt(v), nil
    case int64:
        return ValidateInt64(v), nil
    case []int:
        return ValidateSlice(v)
    default:
        return false, fmt.Errorf("unsupported input type: %T", input)
    }
}

// ValidateAadhaar checks if an Aadhaar number (Indian identification
// number) is valid. Aadhaar numbers must be exactly 12 digits, and the
// last digit is a checksum.
func ValidateAadhaar(aadhaarStr string) (bool, error) {
    if len(aadhaarStr) != 12 {
        return false, errors.New("aadhaar numbers should be 12 digits in length")
    }

    digits, err := stringToDigits(aadhaarStr)
    if err != nil {
        return false, errors.New("aadhaar numbers must contain only numbers")
    }

    // Extract the checksum digit (the last digit)
    checksum := digits[len(digits)-1]

    // Calculate the expected checksum for the first 11 digits
    expectedChecksum := calculateChecksum(digits[:len(digits)-1])

    return checksum == expectedChecksum, nil
}

// AppendChecksumString adds the calculated checksum digit to a string.
func AppendChecksumString(s string) (string, error) {
    checksum, err := GenerateFromString(s)
    if err != nil {
        return "", err
    }
    return s + strconv.Itoa(checksum), nil
}

// AppendChecksumInt adds the calculated checksum digit to an integer.
func AppendChecksumInt(n int) string {
    checksum := GenerateInt(n)
    return strconv.Itoa(n) + strconv.Itoa(checksum)
}

// AppendChecksumInt64 adds the calculated checksum digit to an int64.
func AppendChecksumInt64(n int64) string {
    checksum := GenerateInt64(n)
    return strconv.FormatInt(n, 10) + strconv.Itoa(checksum)
}

// AppendChecksumSlice adds the calculated checksum digit to a slice of digits.
func AppendChecksumSlice(digits []int) (string, error) {
    checksum, err := GenerateSlice(digits)
    if err != nil {
        return "", err
    }
    
    result := ""
    for _, d := range digits {
        result += strconv.Itoa(d)
    }
    return result + strconv.Itoa(checksum), nil
}

// AppendChecksum adds the calculated checksum digit to the input.
// Supported types: string, int, int64, []int
// This function is kept for backward compatibility but using the type-specific
// functions (AppendChecksumString, AppendChecksumInt, etc.) is recommended.
func AppendChecksum(input interface{}) (string, error) {
    switch v := input.(type) {
    case string:
        return AppendChecksumString(v)
    case int:
        return AppendChecksumInt(v), nil
    case int64:
        return AppendChecksumInt64(v), nil
    case []int:
        return AppendChecksumSlice(v)
    default:
        return "", fmt.Errorf("unsupported input type: %T", input)
    }
}

// ConvertToDigits converts various input types to a slice of digits.
// Supported types: string, int, int64, []int
// This function provides compatibility with the original API.
func ConvertToDigits(input interface{}) ([]int, error) {
    switch v := input.(type) {
    case string:
        return stringToDigits(v)
    case int:
        return intToDigits(v), nil
    case int64:
        return int64ToDigits(v), nil
    case []int:
        return sliceToDigits(v)
    default:
        return nil, fmt.Errorf("unsupported input type: %T", input)
    }
}

// InvertArray converts input to a slice of digits and reverses it.
// This function provides compatibility with the original API.
func InvertArray(input interface{}) ([]int, error) {
    digits, err := ConvertToDigits(input)
    if err != nil {
        return nil, err
    }
    
    // Create a copy to avoid modifying the input
    result := make([]int, len(digits))
    copy(result, digits)
    reverseDigits(result)
    
    return result, nil
}

// GenerateString is an alias for Generate that returns a string.
// Deprecated: Use Generate instead.
func GenerateString(input interface{}) (string, error) {
    checksum, err := Generate(input)
    if err != nil {
        return "", err
    }
    return strconv.Itoa(checksum), nil
}