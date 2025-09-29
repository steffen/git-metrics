package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"os/exec"
	"strings"
	"time"
)

// PrintLargestFiles prints information about the largest files
func PrintLargestFiles(files []models.FileInformation, totalFilesSize int64, totalBlobs int, totalFiles int) {
	fmt.Println("\nLARGEST FILES #########################################################################################################")
	fmt.Println()
	fmt.Println("Last commit          Blobs           On-disk size           File path")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	// Track totals for the selected files
	var totalSelectedBlobs int
	var totalSelectedSize int64

	// Track truncated paths for footnotes
	var footnotes []Footnote

	// Calculate total size of all files in repository
	for _, file := range files {
		// Get the last change date for the file
		lastChangeCommand := exec.Command("git", "log", "-1", "--format=%cD", "--", file.Path)
		lastChangeOutput, err := lastChangeCommand.Output()
		if err == nil {
			lastChange, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(string(lastChangeOutput)))
			file.LastChange = lastChange
		}

		percentageSize := float64(file.CompressedSize) / float64(totalFilesSize) * 100
		percentageBlobs := float64(file.Blobs) / float64(totalBlobs) * 100

		// Use CreatePathFootnote for consistent truncation and footnote logic
		result := CreatePathFootnote(file.Path, 66, len(footnotes))
		displayPath := result.DisplayPath
		if result.Index > 0 {
			footnotes = append(footnotes, Footnote{
				Index:    result.Index,
				FullPath: result.FullPath,
			})
		}

		fmt.Printf("%-10s  %13s %5.1f %%  %13s %5.1f %%  %s\n",
			file.LastChange.Format("2006"),
			utils.FormatNumber(file.Blobs),
			percentageBlobs,
			utils.FormatSize(file.CompressedSize),
			percentageSize,
			displayPath)

		totalSelectedBlobs += file.Blobs
		totalSelectedSize += file.CompressedSize
	}

	// Print separator and selected files totals row
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Printf("%-10s  %13s %5.1f %%  %13s %5.1f %%  %s\n",
		"    ",
		utils.FormatNumber(totalSelectedBlobs),
		float64(totalSelectedBlobs)/float64(totalBlobs)*100,
		utils.FormatSize(totalSelectedSize),
		float64(totalSelectedSize)/float64(totalFilesSize)*100,
		fmt.Sprintf("├─ Top %s", utils.FormatNumber(len(files))))

	// Print grand totals row
	fmt.Printf("%-10s  %13s %5.1f %%  %13s %5.1f %%  %s\n",
		"    ",
		utils.FormatNumber(totalBlobs),
		100.0,
		utils.FormatSize(totalFilesSize),
		100.0,
		fmt.Sprintf("└─ Out of %s", utils.FormatNumber(totalFiles)))

	// Print footnotes for truncated paths
	if len(footnotes) > 0 {
		fmt.Println()
		for _, footnote := range footnotes {
			fmt.Printf("[%d] %s\n", footnote.Index, footnote.FullPath)
		}
	}
}
