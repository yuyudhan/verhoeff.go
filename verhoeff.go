// FilePath: verhoeff.go

// Package verhoeff implements the Verhoeff algorithm for checksum calculation and validation.
// The Verhoeff algorithm is a checksum formula for error detection developed by Dutch mathematician J. Verhoeff.
// It can detect all single-digit errors and most transposition errors.
// For more information: https://en.wikipedia.org/wiki/Verhoeff_algorithm

package verhoeff

import (
	"errors"
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

// ConvertToDigits converts a string or integer to a slice of digits.
// It returns an error if the input contains non-digit characters.
func ConvertToDigits(input interface{}) ([]int, error) {
	var digits []int

	switch v := input.(type) {
	case int:
		return ConvertToDigits(strconv.Itoa(v))
	case string:
		for _, char := range v {
			if !unicode.IsDigit(char) {
				return nil, errors.New("input contains non-digit characters")
			}
			digit, _ := strconv.Atoi(string(char))
			digits = append(digits, digit)
		}
	default:
		return nil, errors.New("unsupported input type, must be string or int")
	}

	return digits, nil
}

// InvertArray converts input to a slice of digits and reverses it.
// This is a helper function for the Verhoeff algorithm.
func InvertArray(input interface{}) ([]int, error) {
	digits, err := ConvertToDigits(input)
	if err != nil {
		return nil, err
	}

	// Reverse the array
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return digits, nil
}

// Generate calculates the Verhoeff checksum digit for a given input.
// The input can be a string of digits or an integer.
func Generate(input interface{}) (int, error) {
	invertedArray, err := InvertArray(input)
	if err != nil {
		return -1, err
	}

	c := 0
	for i, digit := range invertedArray {
		c = d[c][p[(i+1)%8][digit]]
	}

	return inv[c], nil
}

// GenerateString is a convenience function that returns the checksum as a string.
func GenerateString(input interface{}) (string, error) {
	checksum, err := Generate(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(checksum), nil
}

// Validate checks if a number with its checksum digit is valid.
// The input can be a string of digits or an integer.
func Validate(input interface{}) (bool, error) {
	invertedArray, err := InvertArray(input)
	if err != nil {
		return false, err
	}

	c := 0
	for i, digit := range invertedArray {
		c = d[c][p[i%8][digit]]
	}

	return c == 0, nil
}

// ValidateAadhaar checks if an Aadhaar number (Indian identification number) is valid.
// Aadhaar numbers must be exactly 12 digits, and the last digit is a checksum.
func ValidateAadhaar(aadhaarStr string) (bool, error) {
	if len(aadhaarStr) != 12 {
		return false, errors.New("Aadhaar numbers should be 12 digits in length")
	}

	digits, err := ConvertToDigits(aadhaarStr)
	if err != nil {
		return false, errors.New("Aadhaar numbers must contain only numbers")
	}

	// Extract the checksum digit (the last digit)
	checksum := digits[len(digits)-1]

	// Remove the checksum digit for validation
	number := aadhaarStr[:len(aadhaarStr)-1]

	// Calculate the expected checksum
	expectedChecksum, err := Generate(number)
	if err != nil {
		return false, err
	}

	return checksum == expectedChecksum, nil
}

// AppendChecksum adds the calculated checksum digit to the end of input and returns the result.
func AppendChecksum(input interface{}) (string, error) {
	var inputStr string

	switch v := input.(type) {
	case int:
		inputStr = strconv.Itoa(v)
	case string:
		inputStr = v
	default:
		return "", errors.New("unsupported input type, must be string or int")
	}

	checksum, err := Generate(inputStr)
	if err != nil {
		return "", err
	}

	return inputStr + strconv.Itoa(checksum), nil
}
