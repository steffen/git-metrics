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
	// Test with stdout (should be true when running tests normally)
	result := isTerminal(os.Stdout)
	t.Logf("os.Stdout is terminal: %v", result)
	
	// Test with stderr (should be true when running tests normally)  
	result = isTerminal(os.Stderr)
	t.Logf("os.Stderr is terminal: %v", result)
	
	// We can't easily test the false case without creating pipes,
	// but the function should handle errors gracefully
	var nilFile *os.File
	if nilFile != nil {
		result = isTerminal(nilFile)
		if result {
			t.Error("Expected false for nil file, got true")
		}
	}
}
