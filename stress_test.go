// FilePath: stress_test.go

package verhoeff

import (
    "fmt"
    "math/rand"
    "strings"
    "sync"
    "testing"
    "time"
)

// TestVeryLongNumbers tests with increasingly long numbers
func TestVeryLongNumbers(t *testing.T) {
    lengths := []int{100, 500, 1000, 5000, 10000}

    for _, length := range lengths {
        t.Run(fmt.Sprintf("length_%d", length), func(t *testing.T) {
            // Create a long number by repeating a pattern
            pattern := "123456789"
            repetitions := length/len(pattern) + 1
            longNumber := strings.Repeat(pattern, repetitions)[:length]

            checksum, err := GenerateFromString(longNumber)
            if err != nil {
                t.Fatalf("Failed to generate checksum for %d-digit number: %v",
                    length, err)
            }

            withChecksum := fmt.Sprintf("%s%d", longNumber, checksum)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Failed to handle %d-digit number", length)
            }
        })
    }
}

// TestRandomPatterns tests with pseudo-random patterns
func TestRandomPatterns(t *testing.T) {
    rand.Seed(42) // Deterministic for reproducibility

    for i := 0; i < 100; i++ {
        // Generate a number of random length (10-100 digits)
        length := 10 + rand.Intn(91)
        var number strings.Builder
        number.Grow(length)

        for j := 0; j < length; j++ {
            digit := rand.Intn(10)
            number.WriteRune(rune('0' + digit))
        }

        numberStr := number.String()
        t.Run(fmt.Sprintf("pattern_%d", i), func(t *testing.T) {
            checksum, err := GenerateFromString(numberStr)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }

            withChecksum := fmt.Sprintf("%s%d", numberStr, checksum)
            valid, err := ValidateString(withChecksum)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if !valid {
                t.Errorf("Failed on pseudo-random pattern #%d", i)
            }

            // Also verify error detection
            runes := []rune(withChecksum)
            pos := rand.Intn(len(runes))
            oldDigit := int(runes[pos] - '0')
            newDigit := (oldDigit + 1) % 10
            runes[pos] = rune('0' + newDigit)
            modified := string(runes)

            valid, err = ValidateString(modified)
            if err != nil {
                t.Fatalf("Validation error: %v", err)
            }
            if valid {
                t.Errorf("Failed to detect error in pseudo-random pattern #%d", i)
            }
        })
    }
}

// TestBatchProcessing tests processing many numbers in sequence
func TestBatchProcessing(t *testing.T) {
    batchSize := 1000
    results := make([]struct {
        number   string
        checksum int
    }, batchSize)

    // Generate checksums for batch
    for i := 0; i < batchSize; i++ {
        number := fmt.Sprintf("%012d", i) // 12-digit number padded with zeros
        checksum, err := GenerateFromString(number)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }
        results[i] = struct {
            number   string
            checksum int
        }{number, checksum}
    }

    // Verify all results
    for i, r := range results {
        withChecksum := fmt.Sprintf("%s%d", r.number, r.checksum)
        valid, err := ValidateString(withChecksum)
        if err != nil {
            t.Fatalf("Validation error: %v", err)
        }
        if !valid {
            t.Errorf("Batch processing failed for %s", r.number)
        }

        // Verify consistency - same input should give same output
        if i < 10 {
            number := fmt.Sprintf("%012d", i)
            checksum, err := GenerateFromString(number)
            if err != nil {
                t.Fatalf("Failed to generate checksum: %v", err)
            }
            if results[i].checksum != checksum {
                t.Errorf("Inconsistent result in batch processing")
            }
        }
    }
}

// TestConcurrentSafety tests that the algorithm works correctly with concurrency
func TestConcurrentSafety(t *testing.T) {
    testNumbers := []string{
        "123456789",
        "987654321",
        "111111111",
        "999999999",
        "555555555",
    }

    numGoroutines := 10
    var wg sync.WaitGroup
    wg.Add(numGoroutines)

    results := make([][]struct {
        number   string
        checksum int
    }, numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        go func(idx int) {
            defer wg.Done()
            localResults := make([]struct {
                number   string
                checksum int
            }, len(testNumbers))

            for j, num := range testNumbers {
                checksum, err := GenerateFromString(num)
                if err != nil {
                    t.Errorf("Failed to generate checksum: %v", err)
                    return
                }
                localResults[j] = struct {
                    number   string
                    checksum int
                }{num, checksum}
            }
            results[idx] = localResults
        }(i)
    }

    wg.Wait()

    // Verify all goroutines got the same results
    firstResults := results[0]
    for i := 1; i < numGoroutines; i++ {
        for j, r := range results[i] {
            if r.checksum != firstResults[j].checksum {
                t.Errorf("Concurrent execution produced different results for %s",
                    r.number)
            }
        }
    }

    // Verify the results are correct
    for _, r := range firstResults {
        withChecksum := fmt.Sprintf("%s%d", r.number, r.checksum)
        valid, err := ValidateString(withChecksum)
        if err != nil {
            t.Fatalf("Validation error: %v", err)
        }
        if !valid {
            t.Errorf("Concurrent execution produced invalid checksum")
        }
    }
}

// TestMemoryEfficiency tests handling of very large numbers
func TestMemoryEfficiency(t *testing.T) {
    // Test with a 100,000 digit number
    hugeNumber := strings.Repeat("9", 100000)
    checksum, err := GenerateFromString(hugeNumber)
    if err != nil {
        t.Fatalf("Failed to generate checksum: %v", err)
    }

    withChecksum := fmt.Sprintf("%s%d", hugeNumber, checksum)
    valid, err := ValidateString(withChecksum)
    if err != nil {
        t.Fatalf("Validation error: %v", err)
    }
    if !valid {
        t.Errorf("Failed to handle 100,000 digit number")
    }

    // Process many large numbers to check for memory leaks
    for i := 0; i < 100; i++ {
        large := strings.Repeat("8", 10000)
        _, err := GenerateFromString(large)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }
    }
}

// TestChecksumDistribution tests that checksums are well-distributed
func TestChecksumDistribution(t *testing.T) {
    checksumCounts := make([]int, 10)

    // Generate checksums for sequential numbers
    for i := 0; i < 10000; i++ {
        number := fmt.Sprintf("%08d", i)
        checksum, err := GenerateFromString(number)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }
        checksumCounts[checksum]++
    }

    // Check that all digits appear as checksums
    for digit, count := range checksumCounts {
        if count == 0 {
            t.Errorf("Digit %d never appeared as checksum", digit)
        }

        // Check for reasonable distribution (within 50% of expected)
        expected := 1000.0
        ratio := float64(count) / expected
        if ratio < 0.5 || ratio > 1.5 {
            t.Errorf("Checksum digit %d has skewed distribution: %d occurrences",
                digit, count)
        }
    }
}

// TestWorstCaseScenarios tests patterns that might be challenging
func TestWorstCaseScenarios(t *testing.T) {
    patterns := map[string]string{
        "all_different": "0123456789",
        "palindrome":    "12344321",
        "binary":        "01010101010101010101",
        "high_entropy":  "31415926535897932384",
    }

    for name, pattern := range patterns {
        t.Run(name, func(t *testing.T) {
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
                t.Errorf("Failed for %s pattern", name)
            }
        })
    }
}

// TestIncrementalChanges tests that small changes produce different checksums
func TestIncrementalChanges(t *testing.T) {
    base := "123456789"
    baseChecksum, err := GenerateFromString(base)
    if err != nil {
        t.Fatalf("Failed to generate checksum: %v", err)
    }

    differentCount := 0
    totalCount := 0

    // Test incrementing each digit
    for pos := 0; pos < len(base); pos++ {
        runes := []rune(base)
        digit := int(runes[pos] - '0')
        newDigit := (digit + 1) % 10
        runes[pos] = rune('0' + newDigit)
        modified := string(runes)

        modifiedChecksum, err := GenerateFromString(modified)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }

        totalCount++
        if modifiedChecksum != baseChecksum {
            differentCount++
        }
    }

    // Most changes should produce different checksums
    if differentCount < totalCount/2 {
        t.Errorf("Too few different checksums: %d/%d", differentCount, totalCount)
    }
}

// TestExhaustiveSmallNumbers exhaustively tests all 4-digit numbers
// This test is marked as slow and should be run separately
func TestExhaustiveSmallNumbers(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping exhaustive test in short mode")
    }

    for i := 0; i < 10000; i++ {
        number := fmt.Sprintf("%04d", i)
        checksum, err := GenerateFromString(number)
        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }

        withChecksum := fmt.Sprintf("%s%d", number, checksum)
        valid, err := ValidateString(withChecksum)
        if err != nil {
            t.Fatalf("Validation error: %v", err)
        }
        if !valid {
            t.Errorf("Failed for 4-digit number: %s", number)
        }

        // Test error detection
        for wrongChecksum := 0; wrongChecksum < 10; wrongChecksum++ {
            if wrongChecksum != checksum {
                withWrong := fmt.Sprintf("%s%d", number, wrongChecksum)
                valid, err := ValidateString(withWrong)
                if err != nil {
                    t.Fatalf("Validation error: %v", err)
                }
                if valid {
                    t.Errorf("Failed to detect wrong checksum for: %s", number)
                }
            }
        }
    }
}

// TestPerformanceConsistency tests that performance is consistent
func TestPerformanceConsistency(t *testing.T) {
    testInputs := []string{
        strings.Repeat("1", 1000),
        strings.Repeat("123456789", 111),
        strings.Repeat("9876543210", 100),
        strings.Repeat("0", 1000),
    }

    var maxTime, minTime time.Duration
    minTime = time.Hour // Start with a very large value

    for _, input := range testInputs {
        start := time.Now()
        checksum, err := GenerateFromString(input)
        duration := time.Since(start)

        if err != nil {
            t.Fatalf("Failed to generate checksum: %v", err)
        }

        if duration > maxTime {
            maxTime = duration
        }
        if duration < minTime {
            minTime = duration
        }

        // Verify the result is valid
        withChecksum := fmt.Sprintf("%s%d", input, checksum)
        valid, err := ValidateString(withChecksum)
        if err != nil {
            t.Fatalf("Validation error: %v", err)
        }
        if !valid {
            t.Errorf("Invalid checksum generated")
        }
    }

    // Check that performance doesn't vary too wildly
    if maxTime > minTime*10 {
        t.Errorf("Performance varies too much: min=%v, max=%v", minTime, maxTime)
    }
}

// BenchmarkGenerateLongNumber benchmarks generation for long numbers
func BenchmarkGenerateLongNumber(b *testing.B) {
    longNumber := strings.Repeat("123456789", 100) // 900 digits

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = GenerateFromString(longNumber)
    }
}

// BenchmarkValidateLongNumber benchmarks validation for long numbers
func BenchmarkValidateLongNumber(b *testing.B) {
    longNumber := strings.Repeat("123456789", 100) + "0" // 901 digits

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = ValidateString(longNumber)
    }
}

// BenchmarkConcurrentGenerate benchmarks concurrent generation
func BenchmarkConcurrentGenerate(b *testing.B) {
    numbers := []string{
        "123456789",
        "987654321",
        "111111111",
        "999999999",
    }

    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            _, _ = GenerateFromString(numbers[i%len(numbers)])
            i++
        }
    })
}