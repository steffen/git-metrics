package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"strconv"
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
// cumulative = totals up to this year, delta = per-year additions.
func PrintGrowthHistoryRow(cumulative, delta, previousDelta models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	var authorsPercentage, commitsPercentage, compressedPercentage float64
	if information.TotalAuthors > 0 {
		authorsPercentage = float64(delta.Authors) / float64(information.TotalAuthors) * 100
	}
	if information.TotalCommits > 0 {
		commitsPercentage = float64(delta.Commits) / float64(information.TotalCommits) * 100
	}
	if information.CompressedSize > 0 {
		compressedPercentage = float64(delta.Compressed) / float64(information.CompressedSize) * 100
	}

	var authorsDeltaPercentChange, commitsDeltaPercentChange, compressedDeltaPercentChange float64
	if previousDelta.Authors > 0 {
		authorsDeltaPercentChange = float64(delta.Authors-previousDelta.Authors) / float64(previousDelta.Authors) * 100
	}
	if previousDelta.Commits > 0 {
		commitsDeltaPercentChange = float64(delta.Commits-previousDelta.Commits) / float64(previousDelta.Commits) * 100
	}
	if previousDelta.Compressed > 0 {
		compressedDeltaPercentChange = float64(delta.Compressed-previousDelta.Compressed) / float64(previousDelta.Compressed) * 100
	}

	yearDisplay := strconv.Itoa(cumulative.Year)
	if cumulative.Year == currentYear {
		yearDisplay += "^"
	}

	// Helper to format signed integers with thousand separators
	formatSigned := func(v int) string {
		if v >= 0 {
			return "+" + utils.FormatNumber(v)
		}
		return "-" + utils.FormatNumber(-v)
	}

	authorsDeltaDisplay := formatSigned(delta.Authors)
	commitsDeltaDisplay := formatSigned(delta.Commits)
	sizeDeltaDisplay := utils.FormatSize(delta.Compressed)

	// Helper to format integer percentage with trailing %
	formatPercent := func(v float64) string { return fmt.Sprintf("%d %%", int(v+0.5)) }

	authorsPercentDisplay := formatPercent(authorsPercentage)
	commitsPercentDisplay := formatPercent(commitsPercentage)
	compressedPercentDisplay := formatPercent(compressedPercentage)

	authorsDeltaPercentDisplay := "" // blank for first year
	commitsDeltaPercentDisplay := ""
	compressedDeltaPercentDisplay := ""
	if previousDelta.Year != 0 {
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
		utils.FormatNumber(cumulative.Authors), authorsDeltaDisplay, authorsPercentDisplay, authorsDeltaPercentDisplay,
		utils.FormatNumber(cumulative.Commits), commitsDeltaDisplay, commitsPercentDisplay, commitsDeltaPercentDisplay,
		utils.FormatSize(cumulative.Compressed), sizeDeltaDisplay, compressedPercentDisplay, compressedDeltaPercentDisplay)
}
