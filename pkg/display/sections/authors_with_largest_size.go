package sections

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"git-metrics/pkg/utils"
)

// Format strings for author size contribution table
const (
	// Authors by size section headers and dividers
	formatAuthorsSizeHeader = "\nAUTHORS WITH LARGEST ON-DISK SIZE " +
		"#####################################################################################" // 35 (text+space) + 85 '#' = 120
	formatAuthorsSizeTableHeader = "Year     Author (#1)            On-disk size   Author (#2)            On-disk size   Author (#3)            On-disk size"
	formatAuthorsSizeDivider     = "------------------------------------------------------------------------------------------------------------------------" // 120 '-'

	// Row formats for 3-column layout (size version)
	formatThreeColumnSizeRow = "%-6s │ %-22s%12s %3.0f%% │ %-22s%12s %3.0f%% │ %-22s%12s %3.0f%%\n"
	formatTwoColumnSizeRow   = "%-6s │ %-22s%12s %3.0f%% │ %-22s%12s %3.0f%%\n"
	formatOneColumnSizeRow   = "%-6s │ %-22s%12s %3.0f%%\n"

	// Maximum contributor name length (same as in contributors_with_most_commits.go)
	maxNameLengthSize = 22
)

// authorSizeStats holds author name and total size contribution
type authorSizeStats struct {
	name string
	size int64
}

// truncateAuthorName truncates an author name to maxNameLengthSize and adds ellipsis if needed
func truncateAuthorName(name string) string {
	if utf8.RuneCountInString(name) <= maxNameLengthSize {
		return name
	}

	runes := []rune(name)
	return string(runes[:maxNameLengthSize-3]) + "..."
}

// DisplayAuthorsWithLargestSize displays the top authors by on-disk size contributions per year
func DisplayAuthorsWithLargestSize(authorsByYear map[int][][3]string, totalSizeByYear map[int]int64, allTimeAuthorSizes map[string]int64) {
	fmt.Println(formatAuthorsSizeHeader)
	fmt.Println()
	fmt.Println(formatAuthorsSizeTableHeader)
	fmt.Println(formatAuthorsSizeDivider)

	// Get years and sort them
	var years []int
	for year := range authorsByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// For all-time totals calculation
	var allTimeTotalSize int64

	// Print for each year
	for _, year := range years {
		authors := authorsByYear[year]
		totalSize := totalSizeByYear[year]
		allTimeTotalSize += totalSize

		// Print year row with up to 3 authors
		displayYearRowAuthorsBySize(year, authors, totalSize)
	}

	// Display all-time authors summary
	if len(years) > 0 {
		displayAllTimeAuthorsBySize(allTimeAuthorSizes, allTimeTotalSize)
	}
}

// displayAuthorSizeRow displays a row of authors with their size contributions and percentages
func displayAuthorSizeRow(yearStr string, authors [][3]string, totalSize int64) {
	// Print the row based on how many authors we actually have
	if len(authors) >= 3 {
		// All 3 columns filled
		author1Size, _ := strconv.ParseInt(authors[0][1], 10, 64)
		author1Percentage := float64(author1Size) / float64(totalSize) * 100
		author2Size, _ := strconv.ParseInt(authors[1][1], 10, 64)
		author2Percentage := float64(author2Size) / float64(totalSize) * 100
		author3Size, _ := strconv.ParseInt(authors[2][1], 10, 64)
		author3Percentage := float64(author3Size) / float64(totalSize) * 100

		fmt.Printf(formatThreeColumnSizeRow,
			yearStr,
			truncateAuthorName(authors[0][0]), utils.FormatSize(author1Size), author1Percentage,
			truncateAuthorName(authors[1][0]), utils.FormatSize(author2Size), author2Percentage,
			truncateAuthorName(authors[2][0]), utils.FormatSize(author3Size), author3Percentage)
	} else if len(authors) == 2 {
		// Only 2 columns filled
		author1Size, _ := strconv.ParseInt(authors[0][1], 10, 64)
		author1Percentage := float64(author1Size) / float64(totalSize) * 100
		author2Size, _ := strconv.ParseInt(authors[1][1], 10, 64)
		author2Percentage := float64(author2Size) / float64(totalSize) * 100

		fmt.Printf(formatTwoColumnSizeRow,
			yearStr,
			truncateAuthorName(authors[0][0]), utils.FormatSize(author1Size), author1Percentage,
			truncateAuthorName(authors[1][0]), utils.FormatSize(author2Size), author2Percentage)
	} else if len(authors) == 1 {
		// Only 1 column filled
		author1Size, _ := strconv.ParseInt(authors[0][1], 10, 64)
		author1Percentage := float64(author1Size) / float64(totalSize) * 100

		fmt.Printf(formatOneColumnSizeRow,
			yearStr,
			truncateAuthorName(authors[0][0]), utils.FormatSize(author1Size), author1Percentage)
	}
}

// displayAuthorSizeRowAllTime displays a row of authors with their size contributions for all-time stats
func displayAuthorSizeRowAllTime(yearStr string, authors []authorSizeStats, totalSize int64) {
	// Print the row based on how many authors we actually have
	if len(authors) >= 3 {
		// All 3 columns filled
		author1Percentage := float64(authors[0].size) / float64(totalSize) * 100
		author2Percentage := float64(authors[1].size) / float64(totalSize) * 100
		author3Percentage := float64(authors[2].size) / float64(totalSize) * 100

		fmt.Printf(formatThreeColumnSizeRow,
			yearStr,
			truncateAuthorName(authors[0].name), utils.FormatSize(authors[0].size), author1Percentage,
			truncateAuthorName(authors[1].name), utils.FormatSize(authors[1].size), author2Percentage,
			truncateAuthorName(authors[2].name), utils.FormatSize(authors[2].size), author3Percentage)
	} else if len(authors) == 2 {
		// Only 2 columns filled
		author1Percentage := float64(authors[0].size) / float64(totalSize) * 100
		author2Percentage := float64(authors[1].size) / float64(totalSize) * 100

		fmt.Printf(formatTwoColumnSizeRow,
			yearStr,
			truncateAuthorName(authors[0].name), utils.FormatSize(authors[0].size), author1Percentage,
			truncateAuthorName(authors[1].name), utils.FormatSize(authors[1].size), author2Percentage)
	} else if len(authors) == 1 {
		// Only 1 column filled
		author1Percentage := float64(authors[0].size) / float64(totalSize) * 100

		fmt.Printf(formatOneColumnSizeRow,
			yearStr,
			truncateAuthorName(authors[0].name), utils.FormatSize(authors[0].size), author1Percentage)
	}
}

func displayYearRowAuthorsBySize(year int, authors [][3]string, totalSize int64) {
	yearStr := fmt.Sprintf("%d", year)
	displayAuthorSizeRow(yearStr, authors, totalSize)
}

func displayAllTimeAuthorsBySize(allTimeAuthorSizes map[string]int64, allTimeTotalSize int64) {
	// Create slice for all-time authors
	var allTimeAuthorsList []authorSizeStats

	// Fill authors slice from the map
	for name, size := range allTimeAuthorSizes {
		allTimeAuthorsList = append(allTimeAuthorsList, authorSizeStats{name: name, size: size})
	}

	// Sort by size (descending) and then by name (ascending, case-insensitive) as a secondary criteria
	sort.Slice(allTimeAuthorsList, func(i, j int) bool {
		if allTimeAuthorsList[i].size != allTimeAuthorsList[j].size {
			return allTimeAuthorsList[i].size > allTimeAuthorsList[j].size
		}
		return strings.ToLower(allTimeAuthorsList[i].name) < strings.ToLower(allTimeAuthorsList[j].name)
	})

	// Limit to top 3 authors
	maxDisplayCount := 3
	if len(allTimeAuthorsList) > maxDisplayCount {
		allTimeAuthorsList = allTimeAuthorsList[:maxDisplayCount]
	}

	// Print all-time stats
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	// Print total row for all-time authors
	displayYearRowAuthorsBySizeAllTime("Total", allTimeAuthorsList, allTimeTotalSize)
}

func displayYearRowAuthorsBySizeAllTime(yearStr string, authors []authorSizeStats, totalSize int64) {
	displayAuthorSizeRowAllTime(yearStr, authors, totalSize)
}