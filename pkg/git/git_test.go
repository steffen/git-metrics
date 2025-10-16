package git

import (
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

func TestGetBranchCount(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		cleanup   func(string)
		wantMin   int
		wantErr   bool
	}{
		{
			name: "Repository with at least one branch",
			setupFunc: func() string {
				// Create a temporary directory and initialize a git repo
				tempDir, _ := os.MkdirTemp("", "git-branch-test")
				cmd := exec.Command("git", "init", tempDir)
				cmd.Run()
				// Set a user name and email for the test repo
				exec.Command("git", "-C", tempDir, "config", "user.name", "Test User").Run()
				exec.Command("git", "-C", tempDir, "config", "user.email", "test@example.com").Run()
				// Create an initial commit to ensure there's a branch
				exec.Command("git", "-C", tempDir, "commit", "--allow-empty", "-m", "Initial commit").Run()
				return tempDir
			},
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
			wantMin: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var originalDir string
			if tt.setupFunc != nil {
				originalDir, _ = os.Getwd()
				path := tt.setupFunc()
				os.Chdir(path)
				if tt.cleanup != nil {
					defer func() {
						os.Chdir(originalDir)
						tt.cleanup(path)
					}()
				}
			}

			count, err := GetBranchCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBranchCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && count < tt.wantMin {
				t.Errorf("GetBranchCount() = %v, want at least %v", count, tt.wantMin)
			}
		})
	}
}

func TestGetTagCount(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		cleanup   func(string)
		wantCount int
		wantErr   bool
	}{
		{
			name: "Repository with no tags",
			setupFunc: func() string {
				// Create a temporary directory and initialize a git repo
				tempDir, _ := os.MkdirTemp("", "git-tag-test")
				cmd := exec.Command("git", "init", tempDir)
				cmd.Run()
				// Set a user name and email for the test repo
				exec.Command("git", "-C", tempDir, "config", "user.name", "Test User").Run()
				exec.Command("git", "-C", tempDir, "config", "user.email", "test@example.com").Run()
				// Create an initial commit
				exec.Command("git", "-C", tempDir, "commit", "--allow-empty", "-m", "Initial commit").Run()
				return tempDir
			},
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "Repository with tags",
			setupFunc: func() string {
				// Create a temporary directory and initialize a git repo
				tempDir, _ := os.MkdirTemp("", "git-tag-test-2")
				cmd := exec.Command("git", "init", tempDir)
				cmd.Run()
				// Set a user name and email for the test repo
				exec.Command("git", "-C", tempDir, "config", "user.name", "Test User").Run()
				exec.Command("git", "-C", tempDir, "config", "user.email", "test@example.com").Run()
				// Create an initial commit
				exec.Command("git", "-C", tempDir, "commit", "--allow-empty", "-m", "Initial commit").Run()
				// Create a tag
				exec.Command("git", "-C", tempDir, "tag", "v1.0.0").Run()
				exec.Command("git", "-C", tempDir, "tag", "v2.0.0").Run()
				return tempDir
			},
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var originalDir string
			if tt.setupFunc != nil {
				originalDir, _ = os.Getwd()
				path := tt.setupFunc()
				os.Chdir(path)
				if tt.cleanup != nil {
					defer func() {
						os.Chdir(originalDir)
						tt.cleanup(path)
					}()
				}
			}

			count, err := GetTagCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && count != tt.wantCount {
				t.Errorf("GetTagCount() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}
