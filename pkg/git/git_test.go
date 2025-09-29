package git

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGetGitVersion(t *testing.T) {
	version := GetGitVersion()
	if version == "" || version == "Unknown" {
		t.Fatalf("GetGitVersion() returned %q, expected a non-empty git version", version)
	}
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
		{name: "Non-existent path", path: "/path/does/not/exist", wantErr: true},
		{name: "Path exists but not a git repository", setupFunc: func() string { d, _ := os.MkdirTemp("", "not-git-repo"); return d }, cleanupFunc: func(p string) { os.RemoveAll(p) }, wantErr: true},
		{name: "Valid git repository", setupFunc: func() string { d, _ := os.MkdirTemp("", "git-repo"); exec.Command("git", "init", d).Run(); return d }, cleanupFunc: func(p string) { os.RemoveAll(p) }, wantErr: false},
		{name: "Valid bare repository", setupFunc: func() string {
			d, _ := os.MkdirTemp("", "git-repo-bare")
			exec.Command("git", "init", "--bare", d).Run()
			return d
		}, cleanupFunc: func(p string) { os.RemoveAll(p) }, wantErr: false},
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
				if _, err := os.Stat(gitDir); err != nil {
					t.Errorf("GetGitDirectory() returned path %v that does not exist", gitDir)
				}
			}
		})
	}
}

func TestGetCheckoutGrowthStats(t *testing.T) {
	rates, _, err := GetRateOfChanges()
	if err != nil {
		t.Fatalf("GetRateOfChanges() error: %v", err)
	}
	if len(rates) == 0 {
		t.Skip("no commits available for rate statistics; skipping checkout growth test")
	}
	for year, rateStats := range rates {
		stats, cgErr := GetCheckoutGrowthStats(year, rateStats.YearEndCommitHash, false)
		if cgErr != nil {
			t.Fatalf("GetCheckoutGrowthStats() returned error: %v", cgErr)
		}
		if stats.Year != year {
			t.Errorf("expected Year to be %d, got %d", year, stats.Year)
		}
		if stats.NumberFiles < 0 || stats.NumberDirectories < 0 || stats.MaxPathDepth < 0 || stats.MaxPathLength < 0 || stats.TotalSizeFiles < 0 {
			t.Errorf("invalid stats for year %d: %+v", year, stats)
		}
		break // only need one year for basic validation
	}
}
