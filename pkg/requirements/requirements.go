package requirements

import (
	"fmt"
	"os/exec"
	"runtime"
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

	// Check bash on Windows
	// We required bash git.GetGrowthStats rely on a pipe that only bash can handle
	// We might remove this in the future
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("bash"); err != nil {
			fmt.Println("bash is not installed or not in PATH")
			fmt.Println("Please install Git Bash or make sure bash is available in your PATH")
			fmt.Println("You can use Git Bash or wsl")
			fmt.Println("")
			allRequirementsMet = false
		}
	}

	return allRequirementsMet
}
