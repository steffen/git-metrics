package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestFlagFormat(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "test-git-metrics")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-git-metrics")

	// Test help output format
	cmd = exec.Command("./test-git-metrics", "--help")
	output, _ := cmd.CombinedOutput()
	// Help command exits with status 0, so we expect no error

	helpText := string(output)

	// Check that help text shows double-dash format
	expectedLines := []string{
		"  --debug",
		"  --no-progress",
		"  -r, --repository string",
		"  --version",
	}

	for _, expectedLine := range expectedLines {
		if !strings.Contains(helpText, expectedLine) {
			t.Errorf("Help text missing expected line: %s\nActual help text:\n%s", expectedLine, helpText)
		}
	}
}

func TestFlagFunctionality(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "test-git-metrics")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-git-metrics")

	// Test that both single and double dash work for version
	testCases := []string{"--version", "-version"}
	
	for _, flagStyle := range testCases {
		cmd := exec.Command("./test-git-metrics", flagStyle)
		output, err := cmd.Output()
		if err != nil {
			t.Errorf("Flag %s failed: %v", flagStyle, err)
			continue
		}

		if !strings.Contains(string(output), "git-metrics version") {
			t.Errorf("Flag %s didn't produce expected version output. Got: %s", flagStyle, string(output))
		}
	}
}