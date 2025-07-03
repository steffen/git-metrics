package main

import (
	"git-metrics/pkg/utils"
	"os"
	"testing"
	"time"
)

func TestCalculateYearsMonthsDays(t *testing.T) {
	tests := []struct {
		name       string
		start      time.Time
		end        time.Time
		wantYears  int
		wantMonths int
		wantDays   int
	}{
		{
			name:       "Same day",
			start:      time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 1, 1, 15, 0, 0, 0, time.UTC),
			wantYears:  0,
			wantMonths: 0,
			wantDays:   0,
		},
		{
			name:       "Few days",
			start:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC),
			wantYears:  0,
			wantMonths: 0,
			wantDays:   4,
		},
		{
			name:       "One month",
			start:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			wantYears:  0,
			wantMonths: 1,
			wantDays:   0,
		},
		{
			name:       "One year",
			start:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantYears:  1,
			wantMonths: 0,
			wantDays:   0,
		},
		{
			name:       "Complex case",
			start:      time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC),
			wantYears:  1,
			wantMonths: 2,
			wantDays:   5,
		},
		{
			name:       "Month boundary",
			start:      time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			wantYears:  0,
			wantMonths: 1,
			wantDays:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			years, months, days := utils.CalculateYearsMonthsDays(tt.start, tt.end)
			if years != tt.wantYears || months != tt.wantMonths || days != tt.wantDays {
				t.Errorf("calculateYMD() = %v years, %v months, %v days; want %v years, %v months, %v days",
					years, months, days, tt.wantYears, tt.wantMonths, tt.wantDays)
			}
		})
	}
}

func TestIsTerminal(t *testing.T) {
	// Test case 1: Pipes are not terminals
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer reader.Close()
	defer writer.Close()

	if isTerminal(reader) {
		t.Error("Expected pipe reader not to be a terminal")
	}

	if isTerminal(writer) {
		t.Error("Expected pipe writer not to be a terminal")
	}

	// Test case 2: Regular file is not a terminal
	temporaryFile, err := os.CreateTemp("", "terminal-test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(temporaryFile.Name())
	defer temporaryFile.Close()

	if isTerminal(temporaryFile) {
		t.Error("Expected regular file not to be a terminal")
	}

	// Test case 3: Nil file handling
	var nilFile *os.File
	if isTerminal(nilFile) {
		t.Error("Expected nil file not to be a terminal")
	}

	// Test case 4: Mock os.Stdout and os.Stderr
	// Save original stdout and stderr
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		// Restore original stdout and stderr after test
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	// Create temporary files to replace stdout and stderr
	mockStdout, err := os.CreateTemp("", "mock-stdout")
	if err != nil {
		t.Fatalf("Failed to create mock stdout: %v", err)
	}
	defer os.Remove(mockStdout.Name())
	defer mockStdout.Close()

	mockStderr, err := os.CreateTemp("", "mock-stderr")
	if err != nil {
		t.Fatalf("Failed to create mock stderr: %v", err)
	}
	defer os.Remove(mockStderr.Name())
	defer mockStderr.Close()

	// Temporarily redirect stdout and stderr
	os.Stdout = mockStdout
	os.Stderr = mockStderr

	// Test the mocked stdout and stderr
	if isTerminal(os.Stdout) {
		t.Error("Mocked stdout incorrectly identified as a terminal")
	}

	if isTerminal(os.Stderr) {
		t.Error("Mocked stderr incorrectly identified as a terminal")
	}
}
