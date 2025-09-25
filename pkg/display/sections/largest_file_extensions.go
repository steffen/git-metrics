package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"path/filepath"
	"sort"
	"strconv"
)

// PrintTopFileExtensions prints the top file extensions by size
func PrintTopFileExtensions(blobs []models.FileInformation, totalBlobs int, totalSize int64) {
	extensionStatistics := make(map[string]struct {
		size       int64
		filesCount int
		blobsCount int
	})
	for _, blob := range blobs {
		extension := filepath.Ext(blob.Path)
		if extension == "" {
			extension = "No Extension"
		}
		statistics := extensionStatistics[extension]
		statistics.size += blob.CompressedSize
		statistics.filesCount++
		statistics.blobsCount += blob.Blobs
		extensionStatistics[extension] = statistics
	}

	// Create a slice for sorting.
	type extensionStatistic struct {
		extension  string
		size       int64
		filesCount int
		blobsCount int
	}
	var statistics []extensionStatistic
	for extension, statistic := range extensionStatistics {
		statistics = append(statistics, extensionStatistic{
			extension:  extension,
			size:       statistic.size,
			filesCount: statistic.filesCount,
			blobsCount: statistic.blobsCount,
		})
	}
	sort.Slice(statistics, func(i, j int) bool {
		return statistics[i].size > statistics[j].size
	})

	// Calculate totals from all extensions first
	var totalExtFilesCount, totalExtBlobsCount int
	var totalExtSize int64
	for _, statistic := range extensionStatistics {
		totalExtFilesCount += statistic.filesCount
		totalExtBlobsCount += statistic.blobsCount
		totalExtSize += statistic.size
	}

	// Limit to top 10
	if len(statistics) > 10 {
		statistics = statistics[:10]
	}

	// Track totals for displayed extensions (top 10)
	var selectedFilesCount int
	var selectedBlobsCount int
	var selectedSize int64

	// Display results.
	fmt.Println("\nLARGEST FILE EXTENSIONS ########################################################################")
	fmt.Println()
	fmt.Println("Extension                            Files                  Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, statistic := range statistics {
		percentageFiles := float64(statistic.filesCount) / float64(totalExtFilesCount) * 100
		percentageBlobs := float64(statistic.blobsCount) / float64(totalBlobs) * 100
		percentageSize := float64(statistic.size) / float64(totalSize) * 100
		fmt.Printf("%-28s %13s %5.1f %%  %13s %5.1f %%  %13s %5.1f %%\n",
			statistic.extension, utils.FormatNumber(statistic.filesCount), percentageFiles, utils.FormatNumber(statistic.blobsCount), percentageBlobs, utils.FormatSize(statistic.size), percentageSize)

		selectedFilesCount += statistic.filesCount
		selectedBlobsCount += statistic.blobsCount
		selectedSize += statistic.size
	}

	// Print separator and top 10 totals row
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("%-28s %13s %5.1f %%  %13s %5.1f %%  %13s %5.1f %%\n",
		fmt.Sprintf("├─ Top %s", utils.FormatNumber(len(statistics))),
		utils.FormatNumber(selectedFilesCount),
		float64(selectedFilesCount)/float64(totalExtFilesCount)*100,
		utils.FormatNumber(selectedBlobsCount),
		float64(selectedBlobsCount)/float64(totalExtBlobsCount)*100,
		utils.FormatSize(selectedSize),
		float64(selectedSize)/float64(totalExtSize)*100)

	// Print grand totals row using full totals
	fmt.Printf("%-28s %13s %5.1f %%  %13s %5.1f %%  %13s %5.1f %%\n",
		fmt.Sprintf("└─ Out of %s", utils.FormatNumber(len(extensionStatistics))),
		utils.FormatNumber(totalExtFilesCount),
		100.0, // Always 100% for totals
		utils.FormatNumber(totalExtBlobsCount),
		100.0,
		utils.FormatSize(totalExtSize),
		100.0)
}

// Format strings for extension growth table
const (
	formatExtensionGrowthHeader = "\nLARGEST FILE EXTENSIONS ON-DISK SIZE GROWTH " +
		"#############################################################################" // 44 (text+space) + 76 '#' = 120
	formatExtensionGrowthTableHeader = "Year     Extension (#1)         Growth        Extension (#2)         Growth        Extension (#3)         Growth        Extension (#4)         Growth        Extension (#5)         Growth"
	formatExtensionGrowthDivider     = "----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------" // 184 '-' for 5-column layout
	
	// Row formats for 5-column layout
	formatFiveColumnRow  = "%-6s │ %-18s%10s │ %-18s%10s │ %-18s%10s │ %-18s%10s │ %-18s%10s\n"
	formatFourColumnRow  = "%-6s │ %-18s%10s │ %-18s%10s │ %-18s%10s │ %-18s%10s\n"
	formatThreeColumnRowExt = "%-6s │ %-18s%10s │ %-18s%10s │ %-18s%10s\n"
	formatTwoColumnRowExt   = "%-6s │ %-18s%10s │ %-18s%10s\n"
	formatOneColumnRowExt   = "%-6s │ %-18s%10s\n"
	
	// Maximum extension name length
	maxExtensionNameLength = 18
)

// extensionGrowthStats holds extension growth information
type extensionGrowthStats struct {
	extension  string
	growth     int64
	percentage float64
}

// PrintFileExtensionGrowth displays the top 5 extensions with largest size growth per year
func PrintFileExtensionGrowth(yearlyStatistics map[int]models.GrowthStatistics) {
	if len(yearlyStatistics) < 2 {
		return // Need at least 2 years to calculate growth
	}
	
	fmt.Println(formatExtensionGrowthHeader)
	fmt.Println()
	fmt.Println(formatExtensionGrowthTableHeader)
	fmt.Println(formatExtensionGrowthDivider)
	
	// Get years and sort them
	var years []int
	for year := range yearlyStatistics {
		years = append(years, year)
	}
	sort.Ints(years)
	
	// Calculate extension statistics for each year
	yearlyExtensionStats := make(map[int]map[string]int64)
	
	for _, year := range years {
		stats := yearlyStatistics[year]
		extensionSizes := make(map[string]int64)
		
		for _, blob := range stats.LargestFiles {
			extension := filepath.Ext(blob.Path)
			if extension == "" {
				extension = "No Extension"
			}
			extensionSizes[extension] += blob.CompressedSize
		}
		
		yearlyExtensionStats[year] = extensionSizes
	}
	
	// Display growth for each year (starting from second year)
	for i := 1; i < len(years); i++ {
		currentYear := years[i]
		previousYear := years[i-1]
		
		currentStats := yearlyExtensionStats[currentYear]
		previousStats := yearlyExtensionStats[previousYear]
		
		// Calculate growth for each extension
		var growthStats []extensionGrowthStats
		totalGrowth := int64(0)
		
		for extension, currentSize := range currentStats {
			previousSize := previousStats[extension] // will be 0 if extension didn't exist previously
			growth := currentSize - previousSize
			if growth > 0 {
				growthStats = append(growthStats, extensionGrowthStats{
					extension: extension,
					growth:    growth,
				})
				totalGrowth += growth
			}
		}
		
		// Sort by growth (descending)
		sort.Slice(growthStats, func(i, j int) bool {
			return growthStats[i].growth > growthStats[j].growth
		})
		
		// Calculate percentages
		for j := range growthStats {
			if totalGrowth > 0 {
				growthStats[j].percentage = float64(growthStats[j].growth) / float64(totalGrowth) * 100
			}
		}
		
		// Limit to top 5
		if len(growthStats) > 5 {
			growthStats = growthStats[:5]
		}
		
		// Display the row
		displayExtensionGrowthRow(strconv.Itoa(currentYear), growthStats)
	}
}

// truncateExtensionName truncates an extension name to maxExtensionNameLength and adds ellipsis if needed
func truncateExtensionName(name string) string {
	if len(name) <= maxExtensionNameLength {
		return name
	}
	return name[:maxExtensionNameLength-3] + "..."
}

// displayExtensionGrowthRow displays a row of extensions with their growth
func displayExtensionGrowthRow(yearStr string, growthStats []extensionGrowthStats) {
	// Prepare data arrays for each column
	var extensions [5]string
	var growths [5]string
	
	for i := 0; i < 5; i++ {
		if i < len(growthStats) {
			extensions[i] = truncateExtensionName(growthStats[i].extension)
			growths[i] = utils.FormatSize(growthStats[i].growth)
		} else {
			extensions[i] = ""
			growths[i] = ""
		}
	}
	
	// Print the row based on how many extensions we actually have
	switch len(growthStats) {
	case 5:
		fmt.Printf(formatFiveColumnRow,
			yearStr,
			extensions[0], growths[0],
			extensions[1], growths[1],
			extensions[2], growths[2],
			extensions[3], growths[3],
			extensions[4], growths[4])
	case 4:
		fmt.Printf(formatFourColumnRow,
			yearStr,
			extensions[0], growths[0],
			extensions[1], growths[1],
			extensions[2], growths[2],
			extensions[3], growths[3])
	case 3:
		fmt.Printf(formatThreeColumnRowExt,
			yearStr,
			extensions[0], growths[0],
			extensions[1], growths[1],
			extensions[2], growths[2])
	case 2:
		fmt.Printf(formatTwoColumnRowExt,
			yearStr,
			extensions[0], growths[0],
			extensions[1], growths[1])
	case 1:
		fmt.Printf(formatOneColumnRowExt,
			yearStr,
			extensions[0], growths[0])
	}
}
