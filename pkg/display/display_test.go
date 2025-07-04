package display

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

// Helper function to capture stdout during test
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

func TestPrintLargestFiles(t *testing.T) {
	files := []models.FileInformation{
		{
			Path:           "file1.txt",
			Blobs:          10,
			CompressedSize: 1024 * 1024, // 1 MB
			LastChange:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Path:           "file2.txt",
			Blobs:          5,
			CompressedSize: 1024 * 512, // 512 KB
			LastChange:     time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	totalFilesCompressedSize := int64(1024*1024 + 1024*512) // 1.5 MB
	totalBlobs := 15
	totalFiles := 2

	output := captureOutput(func() {
		PrintLargestFiles(files, totalFilesCompressedSize, totalBlobs, totalFiles)
	})

	// Check that the output contains the expected headers and files
	expectedStrings := []string{
		"LARGEST FILES",
		"File path",
		"Last commit",
		"Blobs",
		"On-disk size",
		"file1.txt",
		"file2.txt",
		"1.0 MB",
		"512.0 KB",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintTopFileExtensions(t *testing.T) {
	files := []models.FileInformation{
		{
			Path:           "file1.txt",
			Blobs:          10,
			CompressedSize: 1024 * 1024, // 1 MB
		},
		{
			Path:           "file2.txt",
			Blobs:          5,
			CompressedSize: 1024 * 512, // 512 KB
		},
		{
			Path:           "script.go",
			Blobs:          3,
			CompressedSize: 1024 * 256, // 256 KB
		},
		{
			Path:           "README", // No extension
			Blobs:          1,
			CompressedSize: 1024 * 128, // 128 KB
		},
	}

	totalBlobs := 19
	totalSize := int64(1024*1024 + 1024*512 + 1024*256 + 1024*128) // 1.875 MB

	output := captureOutput(func() {
		PrintTopFileExtensions(files, totalBlobs, totalSize)
	})

	// Check that the output contains the expected headers and extension information
	expectedStrings := []string{
		"LARGEST FILE EXTENSIONS",
		"Extension",
		"Files",
		"Blobs",
		"On-disk size",
		".txt",
		".go",
		"No Extension",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintGrowthTableHeader(t *testing.T) {
	output := captureOutput(func() {
		PrintGrowthTableHeader()
	})

	expectedStrings := []string{
		"HISTORIC & ESTIMATED GROWTH",
		"Year",
		"Commits",
		"Trees",
		"Blobs",
		"On-disk size",
		"------------------------------------------------------------------------------------------------",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintGrowthTableRow(t *testing.T) {
	statistics := models.GrowthStatistics{
		Year:       2023,
		Commits:    1000,
		Trees:      2000,
		Blobs:      3000,
		Compressed: 1024 * 1024 * 5, // 5 MB
	}

	previous := models.GrowthStatistics{
		Year:       2022,
		Commits:    800,
		Trees:      1500,
		Blobs:      2500,
		Compressed: 1024 * 1024 * 4, // 4 MB
	}

	information := models.RepositoryInformation{
		TotalCommits:   1000,
		TotalTrees:     2000,
		TotalBlobs:     3000,
		CompressedSize: 1024 * 1024 * 5, // 5 MB
	}

	output := captureOutput(func() {
		PrintGrowthTableRow(statistics, previous, information, false, 2023)
	})

	// Check that the output contains the year and some formatted numbers
	if !strings.Contains(output, "2023") {
		t.Errorf("Expected output to contain the year '2023', but it doesn't.\nOutput: %s", output)
	}

	if !strings.Contains(output, "1,000") {
		t.Errorf("Expected output to contain formatted commits '1,000', but it doesn't.\nOutput: %s", output)
	}

	if !strings.Contains(output, "3,000") {
		t.Errorf("Expected output to contain formatted blobs '3,000', but it doesn't.\nOutput: %s", output)
	}

	if !strings.Contains(output, "5.0 MB") {
		t.Errorf("Expected output to contain formatted size '5.0 MB', but it doesn't.\nOutput: %s", output)
	}
}

func TestPrintGrowthTableRowCurrentYear(t *testing.T) {
	statistics := models.GrowthStatistics{
		Year:       2025,
		Commits:    100,
		Trees:      200,
		Blobs:      300,
		Compressed: 1024 * 1024, // 1 MB
	}

	previous := models.GrowthStatistics{
		Year:       2024,
		Commits:    0,
		Trees:      0,
		Blobs:      0,
		Compressed: 0,
	}

	information := models.RepositoryInformation{
		TotalCommits:   100,
		TotalTrees:     200,
		TotalBlobs:     300,
		CompressedSize: 1024 * 1024, // 1 MB
	}

	output := captureOutput(func() {
		PrintGrowthTableRow(statistics, previous, information, false, 2025)
	})

	// Check that the output contains the separator line when it's the current year
	if !strings.Contains(output, "------------------------------------------------------------------------------------------------") {
		t.Errorf("Expected output to contain separator line for current year, but it doesn't.\nOutput: %s", output)
	}

	// Check that the output contains the current year marker
	if !strings.Contains(output, "2025^") {
		t.Errorf("Expected output to contain current year marker '2025^', but it doesn't.\nOutput: %s", output)
	}

	// Count the number of separator lines
	separatorCount := strings.Count(output, "------------------------------------------------------------------------------------------------")
	if separatorCount != 1 {
		t.Errorf("Expected exactly 1 separator line, but got %d.\nOutput: %s", separatorCount, output)
	}
}

func TestPrintGrowthTableSingleYear(t *testing.T) {
	// Test the scenario where we have only one year of data (the current year)
	// This should reproduce the issue where we get duplicate separator lines
	statistics := models.GrowthStatistics{
		Year:       2025,
		Commits:    100,
		Trees:      200,
		Blobs:      300,
		Compressed: 1024 * 1024, // 1 MB
	}

	previous := models.GrowthStatistics{
		Year:       0, // No previous year
		Commits:    0,
		Trees:      0,
		Blobs:      0,
		Compressed: 0,
	}

	information := models.RepositoryInformation{
		TotalCommits:   100,
		TotalTrees:     200,
		TotalBlobs:     300,
		CompressedSize: 1024 * 1024, // 1 MB
	}

	output := captureOutput(func() {
		PrintGrowthTableHeader()
		PrintGrowthTableRow(statistics, previous, information, false, 2025)
		// Simulate the "no estimation possible" case
		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println("No growth estimation possible: Repository is too young")
	})

	// Count the number of separator lines
	separatorCount := strings.Count(output, "------------------------------------------------------------------------------------------------")

	// We should have exactly 2 separator lines in the proper case:
	// 1. From PrintGrowthTableHeader
	// 2. From the "no estimation possible" case
	// (No separator from PrintGrowthTableRow when it's the only year)
	if separatorCount != 2 {
		t.Errorf("Expected exactly 2 separator lines, but got %d.\nOutput: %s", separatorCount, output)
	}

	// Check that we don't have consecutive separator lines (the actual issue)
	if strings.Contains(output, "------------------------------------------------------------------------------------------------\n------------------------------------------------------------------------------------------------") {
		t.Errorf("Found consecutive separator lines (duplicate separators).\nOutput: %s", output)
	}
}

func TestPrintGrowthTableMultipleYears(t *testing.T) {
	// Test the scenario where we have multiple years of data
	// This should have a separator before the current year
	information := models.RepositoryInformation{
		TotalCommits:   200,
		TotalTrees:     400,
		TotalBlobs:     600,
		CompressedSize: 1024 * 1024 * 2, // 2 MB
	}

	output := captureOutput(func() {
		PrintGrowthTableHeader()

		// Print a previous year
		previousYear := models.GrowthStatistics{
			Year:       2024,
			Commits:    50,
			Trees:      100,
			Blobs:      150,
			Compressed: 512 * 1024, // 0.5 MB
		}
		PrintGrowthTableRow(previousYear, models.GrowthStatistics{}, information, false, 2025)

		// Print the current year
		currentYear := models.GrowthStatistics{
			Year:       2025,
			Commits:    200,
			Trees:      400,
			Blobs:      600,
			Compressed: 1024 * 1024 * 2, // 2 MB
		}
		PrintGrowthTableRow(currentYear, previousYear, information, false, 2025)

		// Simulate the "no estimation possible" case
		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println("No growth estimation possible: Repository is too young")
	})

	// Count the number of separator lines
	separatorCount := strings.Count(output, "------------------------------------------------------------------------------------------------")

	// We should have exactly 3 separator lines when there are multiple years:
	// 1. From PrintGrowthTableHeader
	// 2. From PrintGrowthTableRow (current year, since previous.Year > 0)
	// 3. From the "no estimation possible" case
	if separatorCount != 3 {
		t.Errorf("Expected exactly 3 separator lines for multiple years, but got %d.\nOutput: %s", separatorCount, output)
	}

	// Check that the current year marker is present
	if !strings.Contains(output, "2025^") {
		t.Errorf("Expected output to contain current year marker '2025^', but it doesn't.\nOutput: %s", output)
	}

	// Check that we don't have consecutive separator lines
	if strings.Contains(output, "------------------------------------------------------------------------------------------------\n------------------------------------------------------------------------------------------------") {
		t.Errorf("Found consecutive separator lines (duplicate separators).\nOutput: %s", output)
	}
}
