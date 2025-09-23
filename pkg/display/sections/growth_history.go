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
	// Adjusted spacing per request (narrower) with separators only between blocks (none after Year)
	fmt.Println("Year     Authors        Δ     %      Δ%       Commits          Δ     %      Δ%   On-disk size            Δ     %      Δ%")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
}

// PrintGrowthHistoryRow prints a combined cumulative + delta row.
// statistics contains both cumulative totals and pre-calculated delta values.
func PrintGrowthHistoryRow(statistics, _, previousStats models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	// Use pre-calculated values from the statistics struct
	authorsPercentage := statistics.AuthorsPercent
	commitsPercentage := statistics.CommitsPercent
	compressedPercentage := statistics.CompressedPercent

	authorsDeltaPercentChange := statistics.AuthorsDeltaPercent
	commitsDeltaPercentChange := statistics.CommitsDeltaPercent
	compressedDeltaPercentChange := statistics.CompressedDeltaPercent

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

	authorsDeltaDisplay := formatSigned(statistics.AuthorsDelta)
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

	// Helper to format integer percentage with trailing %
	formatPercent := func(v float64) string { return fmt.Sprintf("%d %%", int(v+0.5)) }

	authorsPercentDisplay := formatPercent(authorsPercentage)
	commitsPercentDisplay := formatPercent(commitsPercentage)
	compressedPercentDisplay := formatPercent(compressedPercentage)

	authorsDeltaPercentDisplay := "" // blank for first year
	commitsDeltaPercentDisplay := ""
	compressedDeltaPercentDisplay := ""
	if previousStats.Year != 0 {
		formatSignedPercent := func(v float64) string {
			iv := int(v + 0.5)
			if iv > 0 {
				return fmt.Sprintf("+%d %%", iv)
			}
			return fmt.Sprintf("%d %%", iv)
		}
		authorsDeltaPercentDisplay = formatSignedPercent(authorsDeltaPercentChange)
		commitsDeltaPercentDisplay = formatSignedPercent(commitsDeltaPercentChange)
		compressedDeltaPercentDisplay = formatSignedPercent(compressedDeltaPercentChange)
	}

	// Print with adjusted spacing: % column narrower, Δ% wider (extra left padding)
	fmt.Printf("%-5s %10s %8s %5s %7s │%12s %10s %5s %7s │%13s %12s %5s %7s\n",
		yearDisplay,
		utils.FormatNumber(statistics.Authors), authorsDeltaDisplay, authorsPercentDisplay, authorsDeltaPercentDisplay,
		utils.FormatNumber(statistics.Commits), commitsDeltaDisplay, commitsPercentDisplay, commitsDeltaPercentDisplay,
		utils.FormatSize(statistics.Compressed), sizeDeltaDisplay, compressedPercentDisplay, compressedDeltaPercentDisplay)
}
