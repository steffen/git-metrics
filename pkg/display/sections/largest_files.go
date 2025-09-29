package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
)

// PrintLargestFiles prints information about the largest files
func PrintLargestFiles(files []models.FileInformation, totalFilesSize int64, totalBlobs int, totalFiles int) {
	fmt.Println("\nLARGEST FILES ##########################################################################################################")
	fmt.Println()
	fmt.Println("       Blobs          On-disk size                                                                                  Path")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	// Track totals for the selected files
	var totalSelectedBlobs int
	var totalSelectedSize int64

	// Track truncated paths for footnotes
	var footnotes []Footnote

	// Calculate total size of all files in repository
	for _, file := range files {
		percentageSize := float64(file.CompressedSize) / float64(totalFilesSize) * 100
		percentageBlobs := float64(file.Blobs) / float64(totalBlobs) * 100

		// Use CreatePathFootnote for consistent truncation and footnote logic
		result := CreatePathFootnote(file.Path, 76, len(footnotes))
		displayPath := result.DisplayPath
		if result.Index > 0 {
			footnotes = append(footnotes, Footnote{
				Index:    result.Index,
				FullPath: result.FullPath,
			})
		}

		fmt.Printf("%11s%6.1f %%   %11s%6.1f %%   %s\n",
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
	fmt.Printf("%11s%6.1f %%   %11s%6.1f %%   %s\n",
		utils.FormatNumber(totalSelectedBlobs),
		float64(totalSelectedBlobs)/float64(totalBlobs)*100,
		utils.FormatSize(totalSelectedSize),
		float64(totalSelectedSize)/float64(totalFilesSize)*100,
		fmt.Sprintf("├─ Top %s", utils.FormatNumber(len(files))))

	// Print grand totals row
	fmt.Printf("%11s%6.1f %%   %11s%6.1f %%   %s\n",
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
