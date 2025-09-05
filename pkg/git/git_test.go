package git

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"git-metrics/pkg/models"
)

func TestCalculateEstimate(t *testing.T) {
	current := models.GrowthStatistics{
		Year:       2023,
		Commits:    1000,
		Trees:      2000,
		Blobs:      3000,
		Compressed: 4000,
	}

	average := models.GrowthStatistics{
		Commits:    100,
		Trees:      200,
		Blobs:      300,
		Compressed: 400,
	}

	expected := models.GrowthStatistics{
		Year:         2024,
		Commits:      1100,
		Trees:        2200,
		Blobs:        3300,
		Compressed:   4400,
		LargestFiles: []models.FileInformation{},
	}

	result := CalculateEstimate(current, average)

	if result.Year != expected.Year {
		t.Errorf("CalculateEstimate() Year = %v, want %v", result.Year, expected.Year)
	}
	if result.Commits != expected.Commits {
		t.Errorf("CalculateEstimate() Commits = %v, want %v", result.Commits, expected.Commits)
	}
	if result.Trees != expected.Trees {
		t.Errorf("CalculateEstimate() Trees = %v, want %v", result.Trees, expected.Trees)
	}
	if result.Blobs != expected.Blobs {
		t.Errorf("CalculateEstimate() Blobs = %v, want %v", result.Blobs, expected.Blobs)
	}
	if result.Compressed != expected.Compressed {
		t.Errorf("CalculateEstimate() Compressed = %v, want %v", result.Compressed, expected.Compressed)
	}
	if len(result.LargestFiles) != len(expected.LargestFiles) {
		t.Errorf("CalculateEstimate() LargestFiles length = %v, want %v", len(result.LargestFiles), len(expected.LargestFiles))
	}
}

func TestCalculateLinearEstimation(t *testing.T) {
	current := models.GrowthStatistics{
		Year:       2023,
		Commits:    1000,
		Trees:      2000,
		Blobs:      3000,
		Compressed: 4000,
	}

	average := models.GrowthStatistics{
		Commits:    100,
		Trees:      200,
		Blobs:      300,
		Compressed: 400,
	}

	result := CalculateLinearEstimation(current, average)

	if result.Method != models.EstimationMethodLinear {
		t.Errorf("CalculateLinearEstimation() Method = %v, want %v", result.Method, models.EstimationMethodLinear)
	}
	if result.Statistics.Year != 2024 {
		t.Errorf("CalculateLinearEstimation() Year = %v, want %v", result.Statistics.Year, 2024)
	}
	if result.Statistics.Commits != 1100 {
		t.Errorf("CalculateLinearEstimation() Commits = %v, want %v", result.Statistics.Commits, 1100)
	}
}

func TestCalculateExponentialEstimation(t *testing.T) {
	yearlyData := []models.GrowthStatistics{
		{Year: 2021, Commits: 100, Trees: 200, Blobs: 300, Compressed: 1000},
		{Year: 2022, Commits: 120, Trees: 240, Blobs: 360, Compressed: 1200},
		{Year: 2023, Commits: 144, Trees: 288, Blobs: 432, Compressed: 1440},
	}

	current := yearlyData[len(yearlyData)-1]
	result := CalculateExponentialEstimation(current, yearlyData)

	if result.Method != models.EstimationMethodExponential {
		t.Errorf("CalculateExponentialEstimation() Method = %v, want %v", result.Method, models.EstimationMethodExponential)
	}
	if result.Statistics.Year != 2024 {
		t.Errorf("CalculateExponentialEstimation() Year = %v, want %v", result.Statistics.Year, 2024)
	}
	// With 20% growth rate, commits should be around 173 (144 * 1.2)
	if result.Statistics.Commits < 160 || result.Statistics.Commits > 180 {
		t.Errorf("CalculateExponentialEstimation() Commits = %v, expected around 173", result.Statistics.Commits)
	}
}

func TestSelectBestEstimationMethod(t *testing.T) {
	current := models.GrowthStatistics{
		Year:       2023,
		Commits:    1000,
		Trees:      2000,
		Blobs:      3000,
		Compressed: 4000,
	}

	average := models.GrowthStatistics{
		Commits:    100,
		Trees:      200,
		Blobs:      300,
		Compressed: 400,
	}

	// Test with insufficient data (should default to linear)
	yearlyData := []models.GrowthStatistics{
		{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 4000},
	}

	result := SelectBestEstimationMethod(current, average, yearlyData)
	if result.Method != models.EstimationMethodLinear {
		t.Errorf("SelectBestEstimationMethod() with insufficient data Method = %v, want %v", result.Method, models.EstimationMethodLinear)
	}

	// Test with more data
	yearlyDataMore := []models.GrowthStatistics{
		{Year: 2021, Commits: 100, Trees: 200, Blobs: 300, Compressed: 1000},
		{Year: 2022, Commits: 200, Trees: 400, Blobs: 600, Compressed: 2000},
		{Year: 2023, Commits: 300, Trees: 600, Blobs: 900, Compressed: 3000},
	}

	result2 := SelectBestEstimationMethod(current, average, yearlyDataMore)
	// The method should be either linear or exponential (we don't care which for this test)
	if result2.Method != models.EstimationMethodLinear && result2.Method != models.EstimationMethodExponential {
		t.Errorf("SelectBestEstimationMethod() returned invalid method: %v", result2.Method)
	}
}

func TestGetGitVersion(t *testing.T) {
	version := GetGitVersion()

	// We can't predict the exact version, but we can check that it's not empty
	// and follows a typical format like "2.35.1" or similar
	if version == "" || version == "Unknown" {
		t.Errorf("GetGitVersion() returned %q, expected a non-empty git version", version)
	}

	// Basic format check - shouldn't contain "git version" prefix since that's stripped
	if strings.Contains(version, "git version") {
		t.Errorf("GetGitVersion() = %q, should not contain 'git version' prefix", version)
	}
}

func TestGetGitDirectory(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		setupFunc   func() string
		cleanupFunc func(string)
		wantErr     bool
	}{
		{
			name:    "Non-existent path",
			path:    "/path/does/not/exist",
			wantErr: true,
		},
		{
			name: "Path exists but not a git repository",
			setupFunc: func() string {
				// Create a temporary directory
				tempDir, _ := os.MkdirTemp("", "not-git-repo")
				return tempDir
			},
			cleanupFunc: func(path string) {
				os.RemoveAll(path)
			},
			wantErr: true,
		},
		{
			name: "Valid git repository",
			setupFunc: func() string {
				// Create a temporary directory and initialize a git repo in it
				tempDir, _ := os.MkdirTemp("", "git-repo")
				cmd := exec.Command("git", "init", tempDir)
				cmd.Run()
				return tempDir
			},
			cleanupFunc: func(path string) {
				os.RemoveAll(path)
			},
			wantErr: false,
		},
		{
			name: "Valid bare repository",
			setupFunc: func() string {
				// Create a temporary directory and initialize a bare repo in it
				tempDir, _ := os.MkdirTemp("", "git-repo-bare")
				cmd := exec.Command("git", "init", "--bare", tempDir)
				cmd.Run()
				return tempDir
			},
			cleanupFunc: func(path string) {
				os.RemoveAll(path)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.setupFunc != nil {
				path = tt.setupFunc()
				if tt.path == "" {
					tt.path = path
				}
			}

			if tt.cleanupFunc != nil && path != "" {
				defer tt.cleanupFunc(path)
			}

			gitDir, err := GetGitDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				// If no error, verify that the path exists and is a git directory
				if _, err := os.Stat(gitDir); err != nil {
					t.Errorf("GetGitDirectory() returned path %v that does not exist", gitDir)
				}
			}
		})
	}
}

// Mock for testing
func mockRunGitCommand(_ bool, _ ...string) ([]byte, error) {
	return []byte("git version 2.35.1"), nil
}
