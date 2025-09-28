package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// CalculateNewEstimate calculates estimated growth using new prediction logic based on delta changes
func CalculateNewEstimate(yearlyStats map[int]models.GrowthStatistics, currentYear int, fetchTimeStr string) []models.GrowthStatistics {
	var estimates []models.GrowthStatistics

	// Get current year statistics and previous year for comparison
	currentStats, currentExists := yearlyStats[currentYear]
	previousStats, previousExists := yearlyStats[currentYear-1]

	if !currentExists || !previousExists {
		return estimates
	}

	// Calculate current year deltas
	currentAuthorsDelta := currentStats.Authors - previousStats.Authors
	currentCommitsDelta := currentStats.Commits - previousStats.Commits
	currentCompressedSizeDelta := currentStats.Compressed - previousStats.Compressed
	currentUncompressedSizeDelta := currentStats.Uncompressed - previousStats.Uncompressed

	// Determine if we need to predict current year or use existing
	var now time.Time
	if fetchTimeStr != "" {
		// Parse the fetch time string (format: "Mon, 02 Jan 2006 15:04 MST")
		if parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04 MST", fetchTimeStr); err == nil {
			now = parsedTime
		} else {
			now = time.Now() // fallback if parsing fails
		}
	} else {
		now = time.Now() // fallback if no fetch time available
	}
	currentYearTime := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	daysPassed := int(now.Sub(currentYearTime).Hours() / 24)

	var predictedCurrentYear models.GrowthStatistics

	// If less than 2 months (60 days) into the year, use previous year data
	if daysPassed < 60 {
		// Simply use the previous year of the latest fetch date
		predictedCurrentYear = previousStats
		predictedCurrentYear.Year = currentYear
	} else {
		// If 2+ months into the year, predict full year by extrapolating current progress
		commitsPerDay := float64(currentCommitsDelta) / float64(daysPassed)
		authorsPerDay := float64(currentAuthorsDelta) / float64(daysPassed)
		compressedSizePerDay := float64(currentCompressedSizeDelta) / float64(daysPassed)
		uncompressedSizePerDay := float64(currentUncompressedSizeDelta) / float64(daysPassed)

		predictedCurrentYear = models.GrowthStatistics{
			Year:         currentYear,
			Authors:      previousStats.Authors + int(authorsPerDay*365),
			Commits:      previousStats.Commits + int(commitsPerDay*365),
			Compressed:   previousStats.Compressed + int64(compressedSizePerDay*365),
			Uncompressed: previousStats.Uncompressed + int64(uncompressedSizePerDay*365),
		}
	}

	estimates = append(estimates, predictedCurrentYear)

	// Calculate deltas for the predicted current year
	currentYearAuthorsDelta := predictedCurrentYear.Authors - previousStats.Authors
	currentYearCommitsDelta := predictedCurrentYear.Commits - previousStats.Commits
	currentYearCompressedSizeDelta := predictedCurrentYear.Compressed - previousStats.Compressed
	currentYearUncompressedSizeDelta := predictedCurrentYear.Uncompressed - previousStats.Uncompressed

	// Store delta values for the predicted current year
	predictedCurrentYear.AuthorsDelta = currentYearAuthorsDelta
	predictedCurrentYear.CommitsDelta = currentYearCommitsDelta
	predictedCurrentYear.CompressedDelta = currentYearCompressedSizeDelta
	predictedCurrentYear.UncompressedDelta = currentYearUncompressedSizeDelta

	// Update the first estimate in the slice with calculated deltas
	estimates[0] = predictedCurrentYear

	// Project future years using linear growth (same delta each year)
	previousEstimate := predictedCurrentYear

	for futureYear := currentYear + 1; futureYear <= currentYear+5; futureYear++ {
		// Use the same deltas from the current year prediction for all future years (linear growth)
		nextAuthorsDelta := currentYearAuthorsDelta
		nextCommitsDelta := currentYearCommitsDelta
		nextCompressedSizeDelta := currentYearCompressedSizeDelta
		nextUncompressedSizeDelta := currentYearUncompressedSizeDelta

		nextEstimate := models.GrowthStatistics{
			Year:         futureYear,
			Authors:      previousEstimate.Authors + nextAuthorsDelta,
			Commits:      previousEstimate.Commits + nextCommitsDelta,
			Compressed:   previousEstimate.Compressed + nextCompressedSizeDelta,
			Uncompressed: previousEstimate.Uncompressed + nextUncompressedSizeDelta,
			// Store delta values
			AuthorsDelta:      nextAuthorsDelta,
			CommitsDelta:      nextCommitsDelta,
			CompressedDelta:   nextCompressedSizeDelta,
			UncompressedDelta: nextUncompressedSizeDelta,
			// Delta percentage values are not needed for linear growth
			AuthorsDeltaPercent:      0,
			CommitsDeltaPercent:      0,
			CompressedDeltaPercent:   0,
			UncompressedDeltaPercent: 0,
		}
		estimates = append(estimates, nextEstimate)

		// Update for next iteration
		previousEstimate = nextEstimate
	}

	return estimates
}

const estimatedGrowthBanner = "ESTIMATED GROWTH ###############################################################################"

// PrintEstimatedGrowthSectionHeader prints only the section banner (with surrounding spacing)
func PrintEstimatedGrowthSectionHeader() {
	fmt.Println()
	fmt.Println(estimatedGrowthBanner)
	fmt.Println()
}

// PrintEstimatedGrowthTableHeader prints only the table column headers + divider (no banner)
func PrintEstimatedGrowthTableHeader() {
	fmt.Println("Year       Commits          Δ     %    ○    Object size            Δ     %    ○   On-disk size            Δ     %    ○")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
}

// PrintGrowthEstimateRow prints a row in the estimated growth table
func PrintGrowthEstimateRow(statistics, previous models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	// Calculate delta values for this estimate row
	currentCommitsDelta := statistics.Commits - previous.Commits
	currentCompressedSizeDelta := statistics.Compressed - previous.Compressed
	currentUncompressedSizeDelta := statistics.Uncompressed - previous.Uncompressed

	// Calculate percentages for this estimate
	var commitsPercentage, compressedPercentage, uncompressedPercentage float64
	if information.TotalCommits > 0 {
		commitsPercentage = float64(currentCommitsDelta) / float64(information.TotalCommits) * 100
	}
	if information.CompressedSize > 0 {
		compressedPercentage = float64(currentCompressedSizeDelta) / float64(information.CompressedSize) * 100
	}
	if information.UncompressedSize > 0 {
		uncompressedPercentage = float64(currentUncompressedSizeDelta) / float64(information.UncompressedSize) * 100
	}

	yearDisplay := strconv.Itoa(statistics.Year) + "*"
	if statistics.Year == currentYear {
		// Current year estimate uses ~ to distinguish from future year estimates (*)
		yearDisplay = strconv.Itoa(statistics.Year) + "~"
	}

	// Helper to format signed integers with thousand separators
	formatSigned := func(v int) string {
		if v >= 0 {
			return "+" + utils.FormatNumber(v)
		}
		return "-" + utils.FormatNumber(-v)
	}

	commitsDeltaDisplay := formatSigned(statistics.Commits - previous.Commits)

	// On-disk size delta with explicit sign when positive
	var sizeDeltaDisplay string
	sizeDelta := statistics.Compressed - previous.Compressed
	if sizeDelta >= 0 {
		formatted := utils.FormatSize(sizeDelta)
		sizeDeltaDisplay = "+" + strings.TrimLeft(formatted, " ")
	} else {
		formatted := utils.FormatSize(-sizeDelta)
		sizeDeltaDisplay = "-" + strings.TrimLeft(formatted, " ")
	}

	// Object size (uncompressed) delta with explicit sign when positive
	var objectSizeDeltaDisplay string
	objectSizeDelta := statistics.Uncompressed - previous.Uncompressed
	if objectSizeDelta >= 0 {
		formatted := utils.FormatSize(objectSizeDelta)
		objectSizeDeltaDisplay = "+" + strings.TrimLeft(formatted, " ")
	} else {
		formatted := utils.FormatSize(-objectSizeDelta)
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

// DisplayUnifiedGrowth handles the complete unified historic and estimated growth section
func DisplayUnifiedGrowth(yearlyStatistics map[int]models.GrowthStatistics, repositoryInformation models.RepositoryInformation, firstCommitTime time.Time, recentFetch string, lastModified string) {
	currentYear := time.Now().Year()

	// Table headers and footnotes are now printed before data collection in main.go
	// Display historic growth data
	var previousDelta models.GrowthStatistics
	for year := repositoryInformation.FirstDate.Year(); year <= currentYear; year++ {
		if cumulative, ok := yearlyStatistics[year]; ok {
			// Add row separator before current year
			if year == currentYear {
				fmt.Println("------------------------------------------------------------------------------------------------------------------------")
			}
			PrintGrowthHistoryRow(cumulative, cumulative, previousDelta, repositoryInformation, currentYear)
			// Add row separator after current year
			if year == currentYear {
				fmt.Println("------------------------------------------------------------------------------------------------------------------------")
			}
			previousDelta = cumulative
		}
	}

	// Display estimated growth data if sufficient history exists
	var estimationEndYear = time.Now().Year() - 1
	var estimationYears = estimationEndYear - firstCommitTime.Year()

	if estimationYears > 5 {
		estimationYears = 5
	}

	if estimationYears > 0 {
		// Use new prediction logic based on delta changes
		estimates := CalculateNewEstimate(yearlyStatistics, currentYear, recentFetch)
		for i, estimate := range estimates {
			var previous models.GrowthStatistics
			if i == 0 {
				// For first estimate (current year), use previous year as comparison
				previous = yearlyStatistics[currentYear-1]
			} else {
				// For subsequent estimates, use previous estimate
				previous = estimates[i-1]
			}
			PrintGrowthEstimateRow(estimate, previous, repositoryInformation, currentYear)
		}
	}

	// Separator and footnotes
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("% columns: each year's delta as share of current totals (^)")
	fmt.Println("○ columns: ○ = Unconcerning, ◑ = On-road to concerning, ● = Concerning")
	if recentFetch != "" {
		fmt.Printf("^ Current totals as of the most recent fetch on %s\n", recentFetch[:16])
	} else {
		fmt.Printf("^ Current totals as of Git directory's last modified: %s\n", lastModified[:16])
	}
	if estimationYears > 0 {
		fmt.Println("~ Estimated growth for current year based on year to date deltas (Δ) extrapolated to full year")
		fmt.Println("* Estimated growth based on current year's estimated delta percentages (Δ%)")
	} else {
		fmt.Println("Growth estimation unavailable: Requires at least 2 years of commit history")
	}
}
