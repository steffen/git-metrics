package requirements

import (
	"fmt"
	"os/exec"
)

// CheckRequirements validates if required dependencies are available
func CheckRequirements() bool {
	allRequirementsMet := true

	// Check git availability
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("git is not installed or not in PATH")
		fmt.Println("Please install git and make sure it's available in your PATH")
		fmt.Println("")
		allRequirementsMet = false
	}

	return allRequirementsMet
}
