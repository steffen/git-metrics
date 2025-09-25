package sections

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"git-metrics/pkg/models"
)

// captureOutput is a helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintGrowthHistoryHeader(t *testing.T) {
	output := captureOutput(func() {
		PrintGrowthHistoryHeader()
	})

	for _, expected := range []string{"HISTORIC GROWTH", "Year", "Authors", "Commits", "On-disk size"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected header to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintGrowthHistoryRow(t *testing.T) {
	// First year statistics with delta values populated
	cumulativePrev := models.GrowthStatistics{
		Year: 2022, Authors: 8, Commits: 800, Compressed: 4 * 1000 * 1000,
		AuthorsDelta: 8, CommitsDelta: 800, CompressedDelta: 4 * 1000 * 1000,
		AuthorsPercent: 80, CommitsPercent: 80, CompressedPercent: 80,
	}

	// Second year statistics with delta values and delta percentages
	cumulative := models.GrowthStatistics{
		Year: 2023, Authors: 10, Commits: 1000, Compressed: 5 * 1000 * 1000,
		AuthorsDelta: 2, CommitsDelta: 200, CompressedDelta: 1 * 1000 * 1000,
		AuthorsPercent: 20, CommitsPercent: 20, CompressedPercent: 20,
		AuthorsDeltaPercent: 100, CommitsDeltaPercent: 100, CompressedDeltaPercent: 100,
	}

	info := models.RepositoryInformation{TotalAuthors: 10, TotalCommits: 1000, CompressedSize: 5 * 1000 * 1000}

	output := captureOutput(func() {
		// First year row (no previous delta)
		PrintGrowthHistoryRow(cumulativePrev, cumulativePrev, models.GrowthStatistics{}, info, 2023)
		// Second year row (with previous delta for Δ%)
		PrintGrowthHistoryRow(cumulative, cumulative, cumulativePrev, info, 2023)
	})

	// Check for cumulative totals and deltas
	expectedSnippets := []string{"2023^", "10", "+2", "+200", "1.0 MB", "5.0 MB", "+100 %"}
	for _, expected := range expectedSnippets {
		if !strings.Contains(output, expected) {
			t.Errorf("expected row output to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintGrowthEstimateRow(t *testing.T) {
	stats := models.GrowthStatistics{
		Year: 2024, Commits: 1100, Trees: 2200, Blobs: 3300, Compressed: 6 * 1000 * 1000,
		CommitsDeltaPercent: 25.0, // Add delta percentage for testing
	}
	prev := models.GrowthStatistics{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 5 * 1000 * 1000}
	info := models.RepositoryInformation{TotalCommits: 1000, TotalTrees: 2000, TotalBlobs: 3000, CompressedSize: 5 * 1000 * 1000}

	output := captureOutput(func() {
		PrintGrowthEstimateRow(stats, prev, info, 2023)
	})

	if !strings.Contains(output, "2024*") {
		t.Errorf("expected estimate row to contain year with * marker. Output: %s", output)
	}
}

func TestPrintLargestFiles(t *testing.T) {
	files := []models.FileInformation{
		{Path: "file1.txt", Blobs: 10, CompressedSize: 1000 * 1000, LastChange: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Path: "file2.txt", Blobs: 5, CompressedSize: 500 * 1000, LastChange: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)},
	}
	totalFilesCompressedSize := int64(1000*1000 + 500*1000)
	totalBlobs := 15
	totalFiles := 2

	output := captureOutput(func() {
		PrintLargestFiles(files, totalFilesCompressedSize, totalBlobs, totalFiles)
	})

	for _, expected := range []string{"LARGEST FILES", "file1.txt", "file2.txt", "1.0 MB", "500.0 KB"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintTopFileExtensions(t *testing.T) {
	files := []models.FileInformation{
		{Path: "file1.txt", Blobs: 10, CompressedSize: 1000 * 1000},
		{Path: "file2.txt", Blobs: 5, CompressedSize: 500 * 1000},
		{Path: "script.go", Blobs: 3, CompressedSize: 256 * 1000},
		{Path: "README", Blobs: 1, CompressedSize: 128 * 1000},
	}
	totalBlobs := 19
	totalSize := int64(1000*1000 + 500*1000 + 256*1000 + 128*1000)

	output := captureOutput(func() {
		PrintTopFileExtensions(files, totalBlobs, totalSize)
	})

	for _, expected := range []string{"LARGEST FILE EXTENSIONS", ".txt", ".go", "No Extension"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestHistoricalAndEstimateSeparation(t *testing.T) {
	// Ensure headers are distinct and not merged
	output := captureOutput(func() {
		PrintGrowthHistoryHeader()
		fmt.Println("(dummy historical rows)")
		PrintEstimatedGrowthSectionHeader()
		PrintEstimatedGrowthTableHeader()
	})

	if !strings.Contains(output, "HISTORIC GROWTH") || !strings.Contains(output, "ESTIMATED GROWTH") {
		t.Errorf("expected both historic and estimated growth headers. Output: %s", output)
	}
}

func TestPrintFileExtensionGrowth(t *testing.T) {
	// Create test data with multiple years
	yearlyStats := map[int]models.GrowthStatistics{
		2022: {
			Year: 2022,
			LargestFiles: []models.FileInformation{
				{Path: "app.go", CompressedSize: 100 * 1000},
				{Path: "README.md", CompressedSize: 50 * 1000},
				{Path: "config.json", CompressedSize: 25 * 1000},
			},
		},
		2023: {
			Year: 2023,
			LargestFiles: []models.FileInformation{
				{Path: "app.go", CompressedSize: 200 * 1000},        // +100KB .go growth
				{Path: "main.go", CompressedSize: 150 * 1000},       // +150KB .go growth (new file)
				{Path: "README.md", CompressedSize: 75 * 1000},      // +25KB .md growth
				{Path: "config.json", CompressedSize: 30 * 1000},    // +5KB .json growth
				{Path: "test.py", CompressedSize: 80 * 1000},        // +80KB .py growth (new extension)
			},
		},
	}

	output := captureOutput(func() {
		PrintFileExtensionGrowth(yearlyStats)
	})

	// Check that the function produces expected content
	expectedSnippets := []string{
		"LARGEST FILE EXTENSIONS ON-DISK SIZE GROWTH",
		"2023", // Year should be displayed
		".go",  // Should show .go extension (highest growth: 250KB)
		".py",  // Should show .py extension (second highest: 80KB)
		".md",  // Should show .md extension (third highest: 25KB)
	}

	for _, expected := range expectedSnippets {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintFileExtensionGrowthInsufficientData(t *testing.T) {
	// Test with only one year of data
	yearlyStats := map[int]models.GrowthStatistics{
		2023: {
			Year: 2023,
			LargestFiles: []models.FileInformation{
				{Path: "app.go", CompressedSize: 200 * 1000},
			},
		},
	}

	output := captureOutput(func() {
		PrintFileExtensionGrowth(yearlyStats)
	})

	// Should produce no output when there's insufficient data
	if output != "" {
		t.Errorf("expected no output with insufficient data, got: %s", output)
	}
}
