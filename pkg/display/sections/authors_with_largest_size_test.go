package sections

import (
	"strings"
	"testing"
)

func TestTruncateAuthorName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Short name",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			name:     "Exact max length name",
			input:    strings.Repeat("a", maxNameLengthSize),
			expected: strings.Repeat("a", maxNameLengthSize),
		},
		{
			name:     "Long name gets truncated",
			input:    strings.Repeat("a", maxNameLengthSize+5),
			expected: strings.Repeat("a", maxNameLengthSize-3) + "...",
		},
		{
			name:     "Unicode name handling",
			input:    "José María García-López de la Cruz",
			expected: "José María García-L...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateAuthorName(tt.input)
			if result != tt.expected {
				t.Errorf("truncateAuthorName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDisplayAuthorsWithLargestSize_DataStructures(t *testing.T) {
	// Test with sample data to ensure no panics and basic functionality
	authorsByYear := map[int][][3]string{
		2023: {
			{"John Doe", "1024", ""},
			{"Jane Smith", "512", ""},
		},
		2024: {
			{"Jane Smith", "2048", ""},
		},
	}

	totalSizeByYear := map[int]int64{
		2023: 1536, // 1024 + 512
		2024: 2048,
	}

	allTimeAuthorSizes := map[string]int64{
		"John Doe":   1024,
		"Jane Smith": 2560, // 512 + 2048
	}

	// This should not panic
	DisplayAuthorsWithLargestSize(authorsByYear, totalSizeByYear, allTimeAuthorSizes)
}