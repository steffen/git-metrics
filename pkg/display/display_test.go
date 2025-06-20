package display

import (
	"bytes"
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
		"1.0 MB",    // 1024*1024 bytes = 1,048,576 bytes ≈ 1.0 MB in decimal
		"524.3 KB",  // 1024*512 bytes = 524,288 bytes ≈ 524.3 KB in decimal
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

	if !strings.Contains(output, "5.2 MB") {
		t.Errorf("Expected output to contain formatted size '5.2 MB', but it doesn't.\nOutput: %s", output)
	}
}
