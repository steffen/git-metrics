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

func TestMailmapSupport(t *testing.T) {
	// Create a temporary directory for the test repository
	tempDir, err := os.MkdirTemp("", "git-repo-mailmap-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repository
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user for commits
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create a commit with one author name/email
	if err := os.WriteFile("file1.txt", []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}
	exec.Command("git", "add", "file1.txt").Run()
	exec.Command("git", "-c", "user.name=John Doe", "-c", "user.email=john@example.com", "commit", "-m", "First commit").Run()

	// Create another commit with a different name/email for the same person
	if err := os.WriteFile("file2.txt", []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create file2.txt: %v", err)
	}
	exec.Command("git", "add", "file2.txt").Run()
	exec.Command("git", "-c", "user.name=J. Doe", "-c", "user.email=jdoe@example.com", "commit", "-m", "Second commit").Run()

	// Test without .mailmap - should see two different authors
	contributors, err := GetContributors()
	if err != nil {
		t.Fatalf("GetContributors() failed: %v", err)
	}

	// Count unique authors from the output
	authorsWithoutMailmap := make(map[string]bool)
	for _, line := range contributors {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 1 {
			authorsWithoutMailmap[parts[0]] = true
		}
	}

	// Create .mailmap file to map both identities to one
	mailmapContent := "John Doe <john@example.com> J. Doe <jdoe@example.com>\n"
	if err := os.WriteFile(".mailmap", []byte(mailmapContent), 0644); err != nil {
		t.Fatalf("Failed to create .mailmap: %v", err)
	}

	// Test with .mailmap - should see only one author
	contributors, err = GetContributors()
	if err != nil {
		t.Fatalf("GetContributors() with mailmap failed: %v", err)
	}

	// Count unique authors with mailmap
	authorsWithMailmap := make(map[string]bool)
	for _, line := range contributors {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 1 {
			authorsWithMailmap[parts[0]] = true
		}
	}

	// With mailmap, we should have only 1 unique author (John Doe)
	// Without checking this would fail if git doesn't respect --use-mailmap
	if len(authorsWithMailmap) != 1 {
		t.Errorf("Expected 1 unique author with mailmap, got %d: %v", len(authorsWithMailmap), authorsWithMailmap)
	}

	// Verify the consolidated name is "John Doe"
	if !authorsWithMailmap["John Doe"] {
		t.Errorf("Expected author 'John Doe' to be present, got: %v", authorsWithMailmap)
	}
}
