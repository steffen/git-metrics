package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"strconv"
)

const historicChangesPerYearBanner = "HISTORIC CHANGES PER YEAR ######################################################################"

// PrintHistoricChangesPerYearHeader prints the header for the historic changes per year table
func PrintHistoricChangesPerYearHeader() {
	fmt.Println()
	fmt.Println(historicChangesPerYearBanner)
	fmt.Println()
	fmt.Println("Year        Commits                  Trees                  Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")
}

// PrintHistoricChangesPerYearRow prints a row of the historic changes per year table.
// statistics here represent the delta (changes) for the year, not cumulative totals.
// previousDelta represents the prior year's delta to compute relative change percentages.
func PrintHistoricChangesPerYearRow(statistics, previousDelta models.GrowthStatistics, currentYear int) {
	// Relative change from previous year's delta. Guard division by zero.
	commitsDifference := 0.0
	treesDifference := 0.0
	blobsDifference := 0.0
	compressedDifference := 0.0

	if previousDelta.Commits != 0 {
		commitsDifference = float64(statistics.Commits-previousDelta.Commits) / float64(previousDelta.Commits) * 100
	}
	if previousDelta.Trees != 0 {
		treesDifference = float64(statistics.Trees-previousDelta.Trees) / float64(previousDelta.Trees) * 100
	}
	if previousDelta.Blobs != 0 {
		blobsDifference = float64(statistics.Blobs-previousDelta.Blobs) / float64(previousDelta.Blobs) * 100
	}
	if previousDelta.Compressed != 0 {
		compressedDifference = float64(statistics.Compressed-previousDelta.Compressed) / float64(previousDelta.Compressed) * 100
	}

	yearDisplay := strconv.Itoa(statistics.Year)
	if statistics.Year == currentYear {
		yearDisplay += "^"
	}

	if previousDelta.Year == 0 { // First year: show raw deltas only without +0% noise
		// Need to preserve column alignment where each percentage block normally takes: ' %+5.0f %%  ' => 8 chars
		blank := "        " // 8 spaces placeholder
		fmt.Printf("%-5s %13s %s %13s %s %13s %s %13s %s\n",
			yearDisplay,
			utils.FormatNumber(statistics.Commits), blank,
			utils.FormatNumber(statistics.Trees), blank,
			utils.FormatNumber(statistics.Blobs), blank,
			utils.FormatSize(statistics.Compressed), blank)
		return
	}

	fmt.Printf("%-5s %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%\n",
		yearDisplay,
		utils.FormatNumber(statistics.Commits), commitsDifference,
		utils.FormatNumber(statistics.Trees), treesDifference,
		utils.FormatNumber(statistics.Blobs), blobsDifference,
		utils.FormatSize(statistics.Compressed), compressedDifference)
}
