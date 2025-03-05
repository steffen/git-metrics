package progress

import (
	"git-metrics/pkg/models"
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	spinner := NewSpinner()

	// Test that spinner returns frames in the correct sequence
	expected := []string{"|", "/", "-", "\\"}
	for i := 0; i < 8; i++ {
		frame := spinner.Next()
		if frame != expected[i%len(expected)] {
			t.Errorf("Expected frame %s at position %d, got %s", expected[i%len(expected)], i, frame)
		}
	}
}

func TestShowProgressFlag(t *testing.T) {
	// Save the original value to restore later
	originalValue := ShowProgress
	defer func() {
		ShowProgress = originalValue
	}()

	// Test that ShowProgress controls visibility
	ShowProgress = false

	// When ShowProgress is false, these functions should return without errors
	StartProgress(2023, models.GrowthStatistics{}, time.Now())
	UpdateProgress()
	StopProgress()

	// Set to true to test that it doesn't crash (can't test actual output easily)
	ShowProgress = true

	// These should execute without error but we can't easily verify console output
	// in unit tests without capturing stdout
	StartProgress(2023, models.GrowthStatistics{}, time.Now())
	time.Sleep(10 * time.Millisecond) // Small delay
	StopProgress()
}

func TestProgressState(t *testing.T) {
	// Save current state to restore later
	originalState := CurrentProgress
	originalShowProgress := ShowProgress
	defer func() {
		CurrentProgress = originalState
		ShowProgress = originalShowProgress
	}()

	ShowProgress = false // Disable actual output during test

	// Test setting and retrieving values
	testYear := 2023
	testStats := models.GrowthStatistics{
		Commits:    100,
		Trees:      50,
		Blobs:      200,
		Compressed: 1024,
	}
	startTime := time.Now()

	StartProgress(testYear, testStats, startTime)

	if CurrentProgress.Year != testYear {
		t.Errorf("Expected Year to be %d, got %d", testYear, CurrentProgress.Year)
	}

	if CurrentProgress.Statistics.Commits != testStats.Commits {
		t.Errorf("Expected Commits to be %d, got %d", testStats.Commits, CurrentProgress.Statistics.Commits)
	}

	if CurrentProgress.Statistics.Trees != testStats.Trees {
		t.Errorf("Expected Trees to be %d, got %d", testStats.Trees, CurrentProgress.Statistics.Trees)
	}

	if CurrentProgress.Statistics.Blobs != testStats.Blobs {
		t.Errorf("Expected Blobs to be %d, got %d", testStats.Blobs, CurrentProgress.Statistics.Blobs)
	}

	if CurrentProgress.Statistics.Compressed != testStats.Compressed {
		t.Errorf("Expected Compressed to be %d, got %d", testStats.Compressed, CurrentProgress.Statistics.Compressed)
	}

	if !CurrentProgress.Active {
		t.Error("Expected Active to be true, got false")
	}

	StopProgress()

	if CurrentProgress.Active {
		t.Error("Expected Active to be false after StopProgress, got true")
	}
}
