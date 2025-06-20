package utils

import (
	"testing"
	"time"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "Kilobytes",
			bytes:    500 * 1000,
			expected: "500.0 KB",
		},
		{
			name:     "Megabytes", 
			bytes:    10 * 1000 * 1000,
			expected: " 10.0 MB",
		},
		{
			name:     "Gigabytes",
			bytes:    5 * 1000 * 1000 * 1000,
			expected: "  5.0 GB",
		},
		{
			name:     "Avoids awkward numbers",
			bytes:    1054867456, // Would be ~1006 MB in binary, now shows as clean GB
			expected: "  1.1 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		expected string
	}{
		{
			name:     "Single digit",
			number:   5,
			expected: "5",
		},
		{
			name:     "Hundreds",
			number:   123,
			expected: "123",
		},
		{
			name:     "Thousands",
			number:   1234,
			expected: "1,234",
		},
		{
			name:     "Millions",
			number:   1234567,
			expected: "1,234,567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNumber(tt.number)
			if result != tt.expected {
				t.Errorf("FormatNumber(%d) = %q, want %q", tt.number, result, tt.expected)
			}
		})
	}
}

func TestTruncatePath(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		maxLength    int
		expectedPath string
	}{
		{
			name:         "Short path",
			path:         "file.txt",
			maxLength:    10,
			expectedPath: "file.txt",
		},
		{
			name:         "Long path",
			path:         "very/long/path/to/some/file.txt",
			maxLength:    20,
			expectedPath: "very/long...file.txt",
		},
		{
			name:         "Exactly at max",
			path:         "exactly-twenty-chars",
			maxLength:    20,
			expectedPath: "exactly-twenty-chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncatePath(tt.path, tt.maxLength)
			if result != tt.expectedPath {
				t.Errorf("TruncatePath(%q, %d) = %q, want %q",
					tt.path, tt.maxLength, result, tt.expectedPath)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "Milliseconds",
			duration: 500 * time.Millisecond,
			expected: "500ms",
		},
		{
			name:     "Seconds",
			duration: 10 * time.Second,
			expected: "10s",
		},
		{
			name:     "Minutes",
			duration: 5 * time.Minute,
			expected: "5m0s",
		},
		{
			name:     "Hours",
			duration: 2 * time.Hour,
			expected: "2h0m0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestMaximum(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "First greater",
			a:        10,
			b:        5,
			expected: 10,
		},
		{
			name:     "Second greater",
			a:        5,
			b:        10,
			expected: 10,
		},
		{
			name:     "Equal values",
			a:        10,
			b:        10,
			expected: 10,
		},
		{
			name:     "Negative values",
			a:        -10,
			b:        -5,
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Maximum(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Maximum(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestLastDayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected int
	}{
		{
			name:     "January",
			time:     time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		{
			name:     "February non-leap year",
			time:     time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		{
			name:     "February leap year",
			time:     time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		{
			name:     "April",
			time:     time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastDayOfMonth(tt.time)
			if result != tt.expected {
				t.Errorf("LastDayOfMonth(%v) = %d, want %d", tt.time, result, tt.expected)
			}
		})
	}
}

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
			name:       "One day",
			start:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			wantYears:  0,
			wantMonths: 0,
			wantDays:   1,
		},
		{
			name:       "Two years, different months",
			start:      time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
			wantYears:  2,
			wantMonths: 2,
			wantDays:   5,
		},
		{
			name:       "Leap year case",
			start:      time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			wantYears:  0,
			wantMonths: 11,
			wantDays:   30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			years, months, days := CalculateYearsMonthsDays(tt.start, tt.end)
			if years != tt.wantYears || months != tt.wantMonths || days != tt.wantDays {
				t.Errorf("CalculateYearsMonthsDays() = %v years, %v months, %v days; want %v years, %v months, %v days",
					years, months, days, tt.wantYears, tt.wantMonths, tt.wantDays)
			}
		})
	}
}
