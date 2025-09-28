package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"strconv"
	"strings"
)

// PrintGrowthHistoryHeader prints the combined historic growth header.
func PrintGrowthHistoryHeader() {
	fmt.Println()
	fmt.Println("HISTORIC GROWTH ################################################################################")
	fmt.Println()
	// Updated header without Authors columns, with Object size instead, and ○ instead of LoC
	fmt.Println("Year       Commits          Δ     %    ○    Object size            Δ     %    ○   On-disk size            Δ     %    ○")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
}

// PrintGrowthHistoryRow prints a combined cumulative + delta row.
// statistics contains both cumulative totals and pre-calculated delta values.
func PrintGrowthHistoryRow(statistics, _, previousStats models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	// Use pre-calculated values from the statistics struct
	commitsPercentage := statistics.CommitsPercent
	compressedPercentage := statistics.CompressedPercent
	uncompressedPercentage := statistics.UncompressedPercent

	yearDisplay := strconv.Itoa(statistics.Year)
	if statistics.Year == currentYear {
		yearDisplay += "^"
	}

	// Helper to format signed integers with thousand separators
	formatSigned := func(v int) string {
		if v >= 0 {
			return "+" + utils.FormatNumber(v)
		}
		return "-" + utils.FormatNumber(-v)
	}

	commitsDeltaDisplay := formatSigned(statistics.CommitsDelta)

	// Size delta with explicit sign when positive
	var sizeDeltaDisplay string
	if statistics.CompressedDelta >= 0 {
		formatted := utils.FormatSize(statistics.CompressedDelta)
		sizeDeltaDisplay = "+" + strings.TrimLeft(formatted, " ")
	} else {
		formatted := utils.FormatSize(-statistics.CompressedDelta)
		sizeDeltaDisplay = "-" + strings.TrimLeft(formatted, " ")
	}

	// Object size (uncompressed) delta with explicit sign when positive
	var objectSizeDeltaDisplay string
	if statistics.UncompressedDelta >= 0 {
		formatted := utils.FormatSize(statistics.UncompressedDelta)
		objectSizeDeltaDisplay = "+" + strings.TrimLeft(formatted, " ")
	} else {
		formatted := utils.FormatSize(-statistics.UncompressedDelta)
		objectSizeDeltaDisplay = "-" + strings.TrimLeft(formatted, " ")
	}

	// Helper to format integer percentage with trailing %
	formatPercent := func(v float64) string { return fmt.Sprintf("%d %%", int(v+0.5)) }

	commitsPercentDisplay := formatPercent(commitsPercentage)
	compressedPercentDisplay := formatPercent(compressedPercentage)
	uncompressedPercentDisplay := formatPercent(uncompressedPercentage)

	// Get Level of Concern (LoC) symbols
	commitsLoC := utils.GetConcernLevel("commits", int64(statistics.Commits))
	objectSizeLoC := utils.GetConcernLevel("object-size", statistics.Uncompressed)
	diskSizeLoC := utils.GetConcernLevel("disk-size", statistics.Compressed)

	// Print with new formatting: Commits | Object size | On-disk size with LoC columns
	fmt.Printf("%-6s %14s %10s %5s %3s │%14s %12s %5s %3s │%14s %12s %5s %3s\n",
		yearDisplay,
		utils.FormatNumber(statistics.Commits), commitsDeltaDisplay, commitsPercentDisplay, commitsLoC,
		utils.FormatSize(statistics.Uncompressed), objectSizeDeltaDisplay, uncompressedPercentDisplay, objectSizeLoC,
		utils.FormatSize(statistics.Compressed), sizeDeltaDisplay, compressedPercentDisplay, diskSizeLoC)
}
