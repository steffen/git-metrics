package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"strconv"
)

const estimatedGrowthBanner = "ESTIMATED GROWTH ###############################################################################"

// PrintEstimatedGrowthSectionHeader prints only the section banner (with surrounding spacing)
func PrintEstimatedGrowthSectionHeader() {
	fmt.Println()
	fmt.Println(estimatedGrowthBanner)
	fmt.Println()
}

// PrintEstimatedGrowthTableHeader prints only the table column headers + divider (no banner)
func PrintEstimatedGrowthTableHeader() {
	fmt.Println("Year        Commits                  Trees                  Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")
}

// PrintGrowthEstimateRow prints a row in the estimated growth table
func PrintGrowthEstimateRow(statistics, previous models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	commitsDifference := float64(statistics.Commits-previous.Commits) / float64(information.TotalCommits) * 100
	treesDifference := float64(statistics.Trees-previous.Trees) / float64(information.TotalTrees) * 100
	blobsDifference := float64(statistics.Blobs-previous.Blobs) / float64(information.TotalBlobs) * 100
	compressedDifference := float64(statistics.Compressed-previous.Compressed) / float64(information.CompressedSize) * 100

	yearDisplay := strconv.Itoa(statistics.Year) + "*"
	fmt.Printf("%-5s %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%\n",
		yearDisplay,
		utils.FormatNumber(statistics.Commits), commitsDifference,
		utils.FormatNumber(statistics.Trees), treesDifference,
		utils.FormatNumber(statistics.Blobs), blobsDifference,
		utils.FormatSize(statistics.Compressed), compressedDifference)
}
