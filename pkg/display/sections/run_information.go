package sections

import (
	"fmt"
	"git-metrics/pkg/git"
	"git-metrics/pkg/utils"
	"runtime"
	"time"
)

// DisplayRunInformation prints information about the system
func DisplayRunInformation() {
	fmt.Println("RUN ####################################################################################################################")
	fmt.Println()
	fmt.Printf("Start time                 %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04 MST"))
	fmt.Printf("Machine                    %d CPU cores with %d GB memory (%s on %s)\n",
		runtime.NumCPU(),
		utils.GetMemoryInGigabytes(),
		utils.GetOperatingSystemInformation(),
		utils.GetChipInformation())
	fmt.Printf("Git metrics version        %s\n", utils.GetGitMetricsVersion())
	fmt.Printf("Git version                %s\n", git.GetGitVersion())
}
