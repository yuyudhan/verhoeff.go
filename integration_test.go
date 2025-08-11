// FilePath: integration_test.go

package verhoeff

import (
    "fmt"
    "strconv"
    "testing"
)

// TestKnownValidChecksums tests against known valid checksum values
func TestKnownValidChecksums(t *testing.T) {
    testCases := []struct {
        input            string
        expectedChecksum int
    }{
        {"236", 3},
        {"12345", 1},
        {"142857", 0},
        {"123456789", 0},
        {"8473643095", 0},
        {"84736430952", 5},
        {"1", 5},
        {"0", 4},
        {"00000000", 1},
    }

    for _, tc := range testCases {
        t.Run(fmt.Sprintf("input_%s", tc.input), func(t *testing.T) {
            calculated, err := GenerateFromString(tc.input)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }

            if calculated != tc.expectedChecksum {
                t.Errorf("Failed for input '%s': expected %d, got %d",
                    tc.input, tc.expectedChecksum, calculated)
            }

            // Also verify that validation works
            withChecksum := fmt.Sprintf("%s%d", tc.input, tc.expectedChecksum)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Validation failed for '%s'", withChecksum)
            }
        })
    }
}

// TestKnownInvalidChecksums tests numbers that should fail validation
func TestKnownInvalidChecksums(t *testing.T) {
    invalidNumbers := []string{
        "2364",         // Should be 2363
        "123450",       // Should be 123451
        "1428571",      // Should be 1428570
        "123456789012", // Random 12-digit number
        "0000000001",   // Should be 0000000000 or 00000000015
    }

    for _, number := range invalidNumbers {
        t.Run(number, func(t *testing.T) {
            valid, err := ValidateString(number)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if valid {
                t.Errorf("Number '%s' should have failed validation", number)
            }
        })
    }
}

// TestAllSingleDigitErrors verifies that all single-digit errors are detected
func TestAllSingleDigitErrors(t *testing.T) {
    baseNumbers := []string{"12345", "987654321", "1111111"}

    for _, base := range baseNumbers {
        checksum, err := GenerateFromString(base)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }
        fullNumber := fmt.Sprintf("%s%d", base, checksum)

        // Try changing each digit to every other possible digit
        for pos := 0; pos < len(fullNumber); pos++ {
            runes := []rune(fullNumber)
            originalDigit := int(runes[pos] - '0')

            for newDigit := 0; newDigit < 10; newDigit++ {
                if newDigit != originalDigit {
                    runes[pos] = rune('0' + newDigit)
                    modified := string(runes)

                    valid, err := ValidateString(modified)
                    if err != nil {
                        t.Fatalf("Validation error: %v", err)
                    }
                    if valid {
                        t.Errorf("Failed to detect single-digit change at position %d "+
                            "in '%s': %d -> %d (modified: '%s')",
                            pos, fullNumber, originalDigit, newDigit, modified)
                    }
                }
            }
        }
    }
}

// TestAdjacentTranspositionDetection tests detection of adjacent digit swaps
func TestAdjacentTranspositionDetection(t *testing.T) {
    baseNumbers := []string{"12345", "987654321", "1234567890"}

    for _, base := range baseNumbers {
        checksum, err := GenerateFromString(base)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }
        fullNumber := fmt.Sprintf("%s%d", base, checksum)

        // Test all adjacent transpositions
        for i := 0; i < len(fullNumber)-1; i++ {
            runes := []rune(fullNumber)

            // Only test if digits are different
            if runes[i] != runes[i+1] {
                // Swap adjacent digits
                runes[i], runes[i+1] = runes[i+1], runes[i]
                transposed := string(runes)

                valid, err := ValidateString(transposed)
                if err != nil {
                    t.Fatalf("Validation error: %v", err)
                }
                if valid {
                    t.Errorf("Failed to detect transposition of positions %d-%d "+
                        "in '%s' (transposed: '%s')",
                        i, i+1, fullNumber, transposed)
                }
            }
        }
    }
}

// TestAadhaarValidationComprehensive tests Aadhaar number validation
func TestAadhaarValidationComprehensive(t *testing.T) {
    // Generate valid Aadhaar-like numbers for testing
    testBases := []string{
        "12345678901",
        "98765432109",
        "11111111111",
        "99999999999",
        "55555555555",
    }

    for _, base := range testBases {
        checksum, err := GenerateFromString(base)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }
        fullAadhaar := fmt.Sprintf("%s%d", base, checksum)

        valid, err := ValidateAadhaar(fullAadhaar)
        if err != nil {
            t.Fatalf("Aadhaar validation error: %v", err)
        }
        if !valid {
            t.Errorf("Valid Aadhaar '%s' failed validation", fullAadhaar)
        }
    }
}

// TestInvalidAadhaarNumbers tests rejection of invalid Aadhaar numbers
func TestInvalidAadhaarNumbers(t *testing.T) {
    testCases := []struct {
        input       string
        shouldError bool
        errorMsg    string
    }{
        {"123456789011", false, "wrong checksum"},
        {"987654321098", false, "wrong checksum"},
        {"12345", true, "wrong length"},
        {"1234567890123", true, "wrong length"},
        {"", true, "empty input"},
        {"12345678901a", true, "invalid character"},
        {"12345678901!", true, "invalid character"},
        {"123456789O12", true, "letter O instead of 0"},
        {"12-34-56-7890", true, "contains dashes"},
    }

    for _, tc := range testCases {
        t.Run(tc.errorMsg, func(t *testing.T) {
            valid, err := ValidateAadhaar(tc.input)
            if tc.shouldError {
                if err == nil {
                    t.Errorf("Expected error for '%s' but got none", tc.input)
                }
            } else {
                if err != nil {
                    t.Fatalf("Unexpected error: %v", err)
                }
                if valid {
                    t.Errorf("Invalid Aadhaar '%s' passed validation", tc.input)
                }
            }
        })
    }
}

// TestAppendChecksumVariousLengths tests append functionality with different input lengths
func TestAppendChecksumVariousLengths(t *testing.T) {
    testInputs := []string{
        "1",
        "12",
        "123",
        "1234",
        "12345",
        "123456",
        "1234567",
        "12345678",
        "123456789",
        "1234567890",
        "12345678901",
        "123456789012",
    }

    for _, input := range testInputs {
        t.Run(fmt.Sprintf("length_%d", len(input)), func(t *testing.T) {
            withChecksum, err := AppendChecksumString(input)
            if err != nil {
                t.Fatalf("Failed to append checksum: %v", err)
            }

            // Verify the appended checksum is valid
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("append_checksum failed for '%s': result '%s' is invalid",
                    input, withChecksum)
            }

            // Verify it has exactly one more character
            if len(withChecksum) != len(input)+1 {
                t.Errorf("append_checksum didn't add exactly one character")
            }

            // Verify the base part is unchanged
            if withChecksum[:len(input)] != input {
                t.Errorf("append_checksum modified the input")
            }
        })
    }
}

// TestRepeatedDigitPatterns tests checksums for repeated digit patterns
func TestRepeatedDigitPatterns(t *testing.T) {
    for digit := 0; digit <= 9; digit++ {
        digitStr := strconv.Itoa(digit)
        for length := 1; length <= 15; length++ {
            repeated := ""
            for i := 0; i < length; i++ {
                repeated += digitStr
            }

            t.Run(fmt.Sprintf("digit_%d_length_%d", digit, length), func(t *testing.T) {
                checksum, err := GenerateFromString(repeated)
                if err != nil {
                    t.Fatalf("Failed to generate checksum: %v", err)
                }

                withChecksum := fmt.Sprintf("%s%d", repeated, checksum)
                valid, err := ValidateString(withChecksum)
                if err != nil {
                    t.Fatalf("Validation error: %v", err)
                }
                if !valid {
                    t.Errorf("Failed to validate %d repetitions of digit %d",
                        length, digit)
                }
            })
        }
    }
}

// TestSequentialPatterns tests ascending and descending sequences
func TestSequentialPatterns(t *testing.T) {
    sequences := []string{
        "01234567890",
        "12345678901",
        "23456789012",
        "0123456789012345678901234567890", // Long sequence
        "9876543210",
        "8765432109",
        "7654321098",
        "9876543210987654321098765432109", // Long descending
    }

    for _, seq := range sequences {
        t.Run(seq, func(t *testing.T) {
            checksum, err := GenerateFromString(seq)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }

            withChecksum := fmt.Sprintf("%s%d", seq, checksum)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Failed to validate sequence: %s", seq)
            }
        })
    }
}

// TestAlternatingPatterns tests alternating digit patterns
func TestAlternatingPatterns(t *testing.T) {
    patterns := []string{
        "0101010101",
        "1010101010",
        "0123012301",
        "1234512345",
        "9090909090",
        "5555555555",
        "123123123123",
        "987987987987",
    }

    for _, pattern := range patterns {
        t.Run(pattern, func(t *testing.T) {
            checksum, err := GenerateFromString(pattern)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }

            withChecksum := fmt.Sprintf("%s%d", pattern, checksum)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Failed to validate pattern: %s", pattern)
            }

            // Also test that invalid checksums are rejected
            wrongChecksum := (checksum + 1) % 10
            withWrong := fmt.Sprintf("%s%d", pattern, wrongChecksum)
            valid, err = ValidateString(withWrong)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if valid {
                t.Errorf("Should reject invalid checksum for pattern: %s", pattern)
            }
        })
    }
}

// TestMathematicalSequences tests special mathematical sequences
func TestMathematicalSequences(t *testing.T) {
    sequences := map[string]string{
        "fibonacci_mod10": "1123581347",
        "primes":          "2357235723",
        "squares_mod10":   "1496561496",
        "pi_digits":       "3141592653",
        "e_digits":        "2718281828",
    }

    for name, seq := range sequences {
        t.Run(name, func(t *testing.T) {
            checksum, err := GenerateFromString(seq)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }

            withChecksum := fmt.Sprintf("%s%d", seq, checksum)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Failed to validate %s sequence: %s", name, seq)
            }
        })
    }
}

// TestLeadingZerosPreservation verifies that leading zeros are handled correctly
func TestLeadingZerosPreservation(t *testing.T) {
    numbersWithLeadingZeros := []string{
        "00000001",
        "00123456",
        "0000000000",
        "00000000000000000001",
    }

    for _, num := range numbersWithLeadingZeros {
        t.Run(num, func(t *testing.T) {
            checksum1, err := GenerateFromString(num)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }

            withChecksum := fmt.Sprintf("%s%d", num, checksum1)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Failed to handle leading zeros in: %s", num)
            }

            // Verify the checksum is consistent
            checksum2, err := GenerateFromString(num)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }
            if checksum1 != checksum2 {
                t.Errorf("Inconsistent checksum for number with leading zeros: %s",
                    num)
            }
        })
    }
}

// TestRealWorldPatterns tests patterns that might appear in real ID numbers
func TestRealWorldPatterns(t *testing.T) {
    examples := []struct {
        base             string
        expectedChecksum int
        description      string
    }{
        {"199812310001", 2, "date-like pattern"},
        {"202401010001", 2, "another date pattern"},
        {"100000000001", 5, "sequential ID"},
        {"999999999999", 9, "maximum value"},
    }

    for _, ex := range examples {
        t.Run(ex.description, func(t *testing.T) {
            calculated, err := GenerateFromString(ex.base)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }
            if calculated != ex.expectedChecksum {
                t.Errorf("Checksum mismatch for '%s': expected %d, got %d",
                    ex.base, ex.expectedChecksum, calculated)
            }
        })
    }
}

// TestConsistencyAcrossTypes tests that different input types give consistent results
func TestConsistencyAcrossTypes(t *testing.T) {
    testCases := []struct {
        stringInput string
        intInput    int
    }{
        {"12345", 12345},
        {"987654", 987654},
        {"1", 1},
        {"999999", 999999},
    }

    for _, tc := range testCases {
        t.Run(tc.stringInput, func(t *testing.T) {
            checksumFromString, err := GenerateFromString(tc.stringInput)
            if err != nil {
                t.Fatalf("Failed to generate checksum from string: %v", err)
            }

            checksumFromInt := GenerateInt(tc.intInput)

            if checksumFromString != checksumFromInt {
                t.Errorf("Inconsistent checksum: string gave %d, int gave %d",
                    checksumFromString, checksumFromInt)
            }

            // Also test validation
            withChecksum := fmt.Sprintf("%s%d", tc.stringInput, checksumFromString)
            validString, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("String validation error: %v", err)
            }

            intWithChecksum := tc.intInput*10 + checksumFromInt
            validInt := ValidateInt(intWithChecksum)

            if validString != validInt {
                t.Errorf("Inconsistent validation: string=%v, int=%v",
                    validString, validInt)
            }
        })
    }
}