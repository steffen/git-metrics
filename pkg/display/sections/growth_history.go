package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"strconv"
)

// PrintGrowthHistoryHeader prints the header for the historic growth table
func PrintGrowthHistoryHeader() {
	fmt.Println()
	fmt.Println("HISTORIC GROWTH ################################################################################")
	fmt.Println()
	fmt.Println("Year        Authors                 Commits           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")
}

// PrintGrowthHistoryRow prints a row of the historic growth table
func PrintGrowthHistoryRow(statistics, previous models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	var authorsDifference, commitsDifference, compressedDifference float64
	if information.TotalAuthors > 0 {
		authorsDifference = float64(statistics.Authors-previous.Authors) / float64(information.TotalAuthors) * 100
	}
	if information.TotalCommits > 0 {
		commitsDifference = float64(statistics.Commits-previous.Commits) / float64(information.TotalCommits) * 100
	}
	if information.CompressedSize > 0 {
		compressedDifference = float64(statistics.Compressed-previous.Compressed) / float64(information.CompressedSize) * 100
	}

	yearDisplay := strconv.Itoa(statistics.Year)
	if statistics.Year == currentYear {
		// Only print separator if there are previous years of data
		if previous.Year > 0 {
			fmt.Println("------------------------------------------------------------------------------------------------")
		}
		yearDisplay += "^"
	}

	fmt.Printf("%-5s %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%\n",
		yearDisplay,
		utils.FormatNumber(statistics.Authors), authorsDifference,
		utils.FormatNumber(statistics.Commits), commitsDifference,
		utils.FormatSize(statistics.Compressed), compressedDifference)
}
