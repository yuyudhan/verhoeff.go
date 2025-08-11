// FilePath: verhoeff_test.go

package verhoeff

import (
    "strconv"
    "testing"
)

func TestConvertToDigits(t *testing.T) {
    tests := []struct {
        name     string
        input    any
        expected []int
        hasError bool
    }{
        {"String input", "12345", []int{1, 2, 3, 4, 5}, false},
        {"Integer input", 12345, []int{1, 2, 3, 4, 5}, false},
        {"Int64 input", int64(12345), []int{1, 2, 3, 4, 5}, false},
        {"Slice input", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}, false},
        {"Empty string", "", []int{}, false},
        {"String with non-digits", "123a45", nil, true},
        {"Float type", 123.45, nil, true},
        {"Invalid digit in slice", []int{1, 2, 10, 4}, nil, true},
        {"Negative digit in slice", []int{1, -2, 3, 4}, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ConvertToDigits(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("ConvertToDigits() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if tt.hasError {
                return
            }

            if len(got) != len(tt.expected) {
                t.Errorf("ConvertToDigits() len = %v, want %v",
                    len(got), len(tt.expected))
                return
            }

            for i := range got {
                if got[i] != tt.expected[i] {
                    t.Errorf("ConvertToDigits() at index %d = %v, want %v",
                        i, got[i], tt.expected[i])
                    return
                }
            }
        })
    }
}

func TestInvertArray(t *testing.T) {
    tests := []struct {
        name     string
        input    any
        expected []int
        hasError bool
    }{
        {"String input", "12345", []int{5, 4, 3, 2, 1}, false},
        {"Integer input", 12345, []int{5, 4, 3, 2, 1}, false},
        {"Single digit", "5", []int{5}, false},
        {"Empty string", "", []int{}, false},
        {"Invalid input", "12a34", nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := InvertArray(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("InvertArray() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if tt.hasError {
                return
            }

            if len(got) != len(tt.expected) {
                t.Errorf("InvertArray() len = %v, want %v",
                    len(got), len(tt.expected))
                return
            }

            for i := range got {
                if got[i] != tt.expected[i] {
                    t.Errorf("InvertArray() at index %d = %v, want %v",
                        i, got[i], tt.expected[i])
                    return
                }
            }
        })
    }
}

func TestGenerate(t *testing.T) {
    tests := []struct {
        name           string
        input          interface{}
        expectedDigit  int
        hasError       bool
    }{
        {"String 236", "236", 3, false},
        {"String 12345", "12345", 1, false},
        {"String 142857", "142857", 0, false},
        {"Integer 236", 236, 3, false},
        {"Integer 12345", 12345, 1, false},
        {"Single digit", "5", 8, false},
        {"Zero", "0", 4, false},
        {"Invalid input", "abc", -1, true},
        {"Empty string", "", 0, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Generate(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("Generate() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if tt.hasError {
                return
            }

            if got != tt.expectedDigit {
                t.Errorf("Generate() = %v, want %v", got, tt.expectedDigit)
            }
        })
    }
}

func TestGenerateString(t *testing.T) {
    tests := []struct {
        name     string
        input    any
        expected string
        hasError bool
    }{
        {"String input", "236", "3", false},
        {"Integer input", 12345, "1", false},
        {"Invalid input", "abc", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := GenerateString(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("GenerateString() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if got != tt.expected {
                t.Errorf("GenerateString() = %v, want %v",
                    got, tt.expected)
            }
        })
    }
}

func TestValidate(t *testing.T) {
    tests := []struct {
        name     string
        input    any
        expected bool
        hasError bool
    }{
        {"Valid string 2363", "2363", true, false},
        {"Invalid string 2364", "2364", false, false},
        {"Valid string 123451", "123451", true, false},
        {"Invalid string 123450", "123450", false, false},
        {"Valid integer 2363", 2363, true, false},
        {"Invalid integer 2364", 2364, false, false},
        {"Single digit valid", "0", true, false},
        {"Invalid input", "12a34", false, true},
        {"Empty string", "", false, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Validate(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("Validate() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if tt.hasError {
                return
            }

            if got != tt.expected {
                t.Errorf("Validate() = %v, want %v", got, tt.expected)
            }
        })
    }
}

func TestValidateAadhaar(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
        hasError bool
    }{
        // Note: These are example Aadhaar numbers for testing
        {"Valid Aadhaar", "234567890124", true, false},
        {"Invalid Aadhaar checksum", "234567890125", false, false},
        {"Too short", "12345678901", false, true},
        {"Too long", "1234567890123", false, true},
        {"Contains letters", "12345678901a", false, true},
        {"Contains special chars", "12345678901!", false, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateAadhaar(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("ValidateAadhaar() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if tt.hasError {
                return
            }

            if got != tt.expected {
                t.Errorf("ValidateAadhaar() = %v, want %v",
                    got, tt.expected)
            }
        })
    }
}

func TestAppendChecksum(t *testing.T) {
    tests := []struct {
        name     string
        input    any
        expected string
        hasError bool
    }{
        {"String input", "236", "2363", false},
        {"Integer input", 12345, "123451", false},
        {"Int64 input", int64(142857), "1428570", false},
        {"Slice input", []int{2, 3, 6}, "2363", false},
        {"Empty string", "", "0", false},
        {"Invalid string", "12a34", "", true},
        {"Invalid slice", []int{1, 2, 10}, "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := AppendChecksum(tt.input)

            if (err != nil) != tt.hasError {
                t.Errorf("AppendChecksum() error = %v, wantErr %v",
                    err, tt.hasError)
                return
            }

            if tt.hasError {
                return
            }

            if got != tt.expected {
                t.Errorf("AppendChecksum() = %v, want %v",
                    got, tt.expected)
            }

            // Verify the appended checksum is valid
            if !tt.hasError && got != "" {
                valid, err := Validate(got)
                if err != nil {
                    t.Errorf("Failed to validate appended checksum: %v", err)
                }
                if !valid {
                    t.Errorf("Appended checksum %v is not valid", got)
                }
            }
        })
    }
}

// Benchmark tests
func BenchmarkGenerate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = Generate("1234567890")
    }
}

func BenchmarkValidate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = Validate("12345678909")
    }
}

func BenchmarkValidateAadhaar(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = ValidateAadhaar("234567890124")
    }
}

func BenchmarkAppendChecksum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = AppendChecksum("1234567890")
    }
}

// Table-driven tests for edge cases
func TestEdgeCases(t *testing.T) {
    t.Run("Large numbers", func(t *testing.T) {
        largeNum := "999999999999999999999999999999"
        checksum, err := Generate(largeNum)
        if err != nil {
            t.Errorf("Failed to generate checksum for large number: %v", err)
        }

        withChecksum := largeNum + strconv.Itoa(checksum)
        valid, err := Validate(withChecksum)
        if err != nil {
            t.Errorf("Failed to validate large number: %v", err)
        }
        if !valid {
            t.Errorf("Large number validation failed")
        }
    })

    t.Run("All zeros", func(t *testing.T) {
        zeros := "0000000000"
        checksum, err := Generate(zeros)
        if err != nil {
            t.Errorf("Failed to generate checksum for zeros: %v", err)
        }
        if checksum != 5 {
            t.Errorf("Expected checksum 5 for all zeros, got %d", checksum)
        }
    })

    t.Run("Sequential numbers", func(t *testing.T) {
        seq := "1234567890"
        checksum, err := Generate(seq)
        if err != nil {
            t.Errorf("Failed to generate checksum for sequential: %v", err)
        }

        withChecksum := seq + strconv.Itoa(checksum)
        valid, err := Validate(withChecksum)
        if err != nil {
            t.Errorf("Failed to validate sequential: %v", err)
        }
        if !valid {
            t.Errorf("Sequential number validation failed")
        }
    })
}