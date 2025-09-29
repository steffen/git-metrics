package git

import (
	"git-metrics/pkg/models"
	"os"
	"os/exec"
	"strings"
	"testing"
)

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

func TestGetGrowthStatsCheckoutData(t *testing.T) {
	// Test that GetGrowthStats now includes checkout growth data
	stats, err := GetGrowthStats(2025, models.GrowthStatistics{}, false)
	if err != nil {
		t.Fatalf("GetGrowthStats() returned error: %v", err)
	}

	if stats.Year != 2025 {
		t.Errorf("expected Year to be 2025, got %d", stats.Year)
	}

	// In a working git repo, we should have checkout growth data
	if stats.NumberFiles == 0 {
		t.Errorf("expected NumberFiles to be greater than 0, got %d", stats.NumberFiles)
	}

	if stats.NumberDirectories == 0 {
		t.Errorf("expected NumberDirectories to be greater than 0, got %d", stats.NumberDirectories)
	}

	if stats.MaxPathDepth < 0 {
		t.Errorf("expected MaxPathDepth to be non-negative, got %d", stats.MaxPathDepth)
	}

	if stats.MaxPathLength <= 0 {
		t.Errorf("expected MaxPathLength to be positive, got %d", stats.MaxPathLength)
	}

	if stats.TotalSizeFiles <= 0 {
		t.Errorf("expected TotalSizeFiles to be positive, got %d", stats.TotalSizeFiles)
	}
}
