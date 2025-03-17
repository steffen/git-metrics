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
