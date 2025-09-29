package sections

import (
	"git-metrics/pkg/models"
	"strings"
	"testing"
)

func TestDisplayCheckoutGrowth(t *testing.T) {
	// Test with empty statistics
	output := captureOutput(func() {
		DisplayCheckoutGrowth(make(map[int]models.GrowthStatistics))
	})
	if output != "" {
		t.Errorf("expected empty output for empty statistics, got: %s", output)
	}

	// Test with sample statistics
	checkoutStats := map[int]models.GrowthStatistics{
		2023: {
			Year:              2023,
			NumberDirectories: 10,
			MaxPathDepth:      3,
			MaxPathLength:     45,
			NumberFiles:       25,
			TotalSizeFiles:    1024000,
		},
		2024: {
			Year:              2024,
			NumberDirectories: 15,
			MaxPathDepth:      4,
			MaxPathLength:     60,
			NumberFiles:       35,
			TotalSizeFiles:    2048000,
		},
	}

	output = captureOutput(func() {
		DisplayCheckoutGrowth(checkoutStats)
	})

	expectedSnippets := []string{
		"CHECKOUT GROWTH",
		"Year",
		"Directories",
		"Max depth",
		"Max path length",
		"Files",
		"Total size",
		"2023",
		"2024",
		"10", // number of directories for 2023
		"15", // number of directories for 2024
		"3",  // max depth for 2023
		"4",  // max depth for 2024
	}

	for _, expected := range expectedSnippets {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestDisplayCheckoutGrowthRow(t *testing.T) {
	stats := models.GrowthStatistics{
		Year:              2023,
		NumberDirectories: 10,
		MaxPathDepth:      3,
		MaxPathLength:     45,
		NumberFiles:       25,
		TotalSizeFiles:    1024000,
	}

	output := captureOutput(func() {
		DisplayCheckoutGrowthRow(stats)
	})

	expectedSnippets := []string{
		"2023",
		"10",   // directories
		"3",    // max depth
		"45",   // max path length  
		"25",   // files
		"1.0",  // should contain part of formatted size (1.0 MB)
	}

	for _, expected := range expectedSnippets {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q.\nOutput: %s", expected, output)
		}
	}
}