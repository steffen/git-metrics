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
	stats := models.GrowthStatistics{Year: 2023, Authors: 10, Commits: 1000, Compressed: 5 * 1000 * 1000}
	prev := models.GrowthStatistics{Year: 2022, Authors: 8, Commits: 800, Compressed: 4 * 1000 * 1000}
	info := models.RepositoryInformation{TotalAuthors: 10, TotalCommits: 1000, CompressedSize: 5 * 1000 * 1000}

	output := captureOutput(func() {
		PrintGrowthHistoryRow(stats, prev, info, 2023)
	})

	for _, expected := range []string{"2023^", "10", "1,000", "5.0 MB"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected row to contain %q.\nOutput: %s", expected, output)
		}
	}
}

func TestPrintGrowthEstimateRow(t *testing.T) {
	stats := models.GrowthStatistics{Year: 2024, Commits: 1100, Trees: 2200, Blobs: 3300, Compressed: 6 * 1000 * 1000}
	prev := models.GrowthStatistics{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 5 * 1000 * 1000}
	info := models.RepositoryInformation{TotalCommits: 1000, TotalTrees: 2000, TotalBlobs: 3000, CompressedSize: 5 * 1000 * 1000}

	output := captureOutput(func() {
		PrintGrowthEstimateRow(stats, prev, info, 2023)
	})

	if !strings.Contains(output, "2024*") {
		if !strings.Contains(output, "2024*") {
			t.Errorf("expected estimate row to contain year with * marker. Output: %s", output)
		}
	}
}

func TestPrintHistoricChangesPerYearHeaderAndRow(t *testing.T) {
	output := captureOutput(func() {
		PrintHistoricChangesPerYearHeader()
		// Simulate two years of deltas
		prevDelta := models.GrowthStatistics{Year: 2022, Commits: 200, Trees: 400, Blobs: 600, Compressed: 1 * 1000 * 1000}
		currentDelta := models.GrowthStatistics{Year: 2023, Commits: 300, Trees: 500, Blobs: 700, Compressed: 2 * 1000 * 1000}
		PrintHistoricChangesPerYearRow(prevDelta, models.GrowthStatistics{}, 2023)
		PrintHistoricChangesPerYearRow(currentDelta, prevDelta, 2023)
	})

	for _, expected := range []string{"HISTORIC CHANGES PER YEAR", "Year", "Commits", "2023^"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected historic changes per year output to contain %q.\nOutput: %s", expected, output)
		}
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
