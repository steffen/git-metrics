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

	// Calculate delta percentage growth rates from current year to apply to future years
	currentYearAuthorsDelta := predictedCurrentYear.Authors - previousStats.Authors
	currentYearCommitsDelta := predictedCurrentYear.Commits - previousStats.Commits
	currentYearCompressedSizeDelta := predictedCurrentYear.Compressed - previousStats.Compressed
	currentYearUncompressedSizeDelta := predictedCurrentYear.Uncompressed - previousStats.Uncompressed

	// Calculate previous year deltas to determine delta percentage growth rate
	twoYearsAgo, twoYearsExists := yearlyStats[currentYear-2]
	if !twoYearsExists {
		return estimates
	}

	previousYearAuthorsDelta := previousStats.Authors - twoYearsAgo.Authors
	previousYearCommitsDelta := previousStats.Commits - twoYearsAgo.Commits
	previousYearCompressedSizeDelta := previousStats.Compressed - twoYearsAgo.Compressed
	previousYearUncompressedSizeDelta := previousStats.Uncompressed - twoYearsAgo.Uncompressed

	// Calculate delta percentage change (Δ% growth rate)
	var authorsDeltaGrowthPercent, commitsDeltaGrowthPercent, compressedSizeDeltaGrowthPercent, uncompressedSizeDeltaGrowthPercent float64
	if previousYearAuthorsDelta > 0 {
		authorsDeltaGrowthPercent = float64(currentYearAuthorsDelta-previousYearAuthorsDelta) / float64(previousYearAuthorsDelta)
	}
	if previousYearCommitsDelta > 0 {
		commitsDeltaGrowthPercent = float64(currentYearCommitsDelta-previousYearCommitsDelta) / float64(previousYearCommitsDelta)
	}
	if previousYearCompressedSizeDelta > 0 {
		compressedSizeDeltaGrowthPercent = float64(currentYearCompressedSizeDelta-previousYearCompressedSizeDelta) / float64(previousYearCompressedSizeDelta)
	}
	if previousYearUncompressedSizeDelta > 0 {
		uncompressedSizeDeltaGrowthPercent = float64(currentYearUncompressedSizeDelta-previousYearUncompressedSizeDelta) / float64(previousYearUncompressedSizeDelta)
	}

	// Store delta percentage values for the predicted current year
	// Calculate delta percentage for current year based on previous year's delta
	if previousYearAuthorsDelta > 0 {
		predictedCurrentYear.AuthorsDeltaPercent = float64(currentYearAuthorsDelta-previousYearAuthorsDelta) / float64(previousYearAuthorsDelta) * 100
	}
	if previousYearCommitsDelta > 0 {
		predictedCurrentYear.CommitsDeltaPercent = float64(currentYearCommitsDelta-previousYearCommitsDelta) / float64(previousYearCommitsDelta) * 100
	}
	if previousYearCompressedSizeDelta > 0 {
		predictedCurrentYear.CompressedDeltaPercent = float64(currentYearCompressedSizeDelta-previousYearCompressedSizeDelta) / float64(previousYearCompressedSizeDelta) * 100
	}
	if previousYearUncompressedSizeDelta > 0 {
		predictedCurrentYear.UncompressedDeltaPercent = float64(currentYearUncompressedSizeDelta-previousYearUncompressedSizeDelta) / float64(previousYearUncompressedSizeDelta) * 100
	}

	// Store delta values for the predicted current year
	predictedCurrentYear.AuthorsDelta = currentYearAuthorsDelta
	predictedCurrentYear.CommitsDelta = currentYearCommitsDelta
	predictedCurrentYear.CompressedDelta = currentYearCompressedSizeDelta
	predictedCurrentYear.UncompressedDelta = currentYearUncompressedSizeDelta

	// Update the first estimate in the slice with calculated delta percentages
	estimates[0] = predictedCurrentYear

	// Project future years using delta percentage growth rates
	previousEstimate := predictedCurrentYear
	previousAuthorsDelta := currentYearAuthorsDelta
	previousCommitsDelta := currentYearCommitsDelta
	previousCompressedSizeDelta := currentYearCompressedSizeDelta
	previousUncompressedSizeDelta := currentYearUncompressedSizeDelta

	for futureYear := currentYear + 1; futureYear <= currentYear+5; futureYear++ {
		// Apply delta growth percentages to previous deltas to get new deltas
		nextAuthorsDelta := previousAuthorsDelta + int(float64(previousAuthorsDelta)*authorsDeltaGrowthPercent)
		nextCommitsDelta := previousCommitsDelta + int(float64(previousCommitsDelta)*commitsDeltaGrowthPercent)
		nextCompressedSizeDelta := previousCompressedSizeDelta + int64(float64(previousCompressedSizeDelta)*compressedSizeDeltaGrowthPercent)
		nextUncompressedSizeDelta := previousUncompressedSizeDelta + int64(float64(previousUncompressedSizeDelta)*uncompressedSizeDeltaGrowthPercent)

		// Calculate delta percentages for future years (consistent with current year calculation)
		var nextAuthorsDeltaPercent, nextCommitsDeltaPercent, nextCompressedSizeDeltaPercent, nextUncompressedSizeDeltaPercent float64
		if previousAuthorsDelta > 0 {
			nextAuthorsDeltaPercent = float64(nextAuthorsDelta-previousAuthorsDelta) / float64(previousAuthorsDelta) * 100
		}
		if previousCommitsDelta > 0 {
			nextCommitsDeltaPercent = float64(nextCommitsDelta-previousCommitsDelta) / float64(previousCommitsDelta) * 100
		}
		if previousCompressedSizeDelta > 0 {
			nextCompressedSizeDeltaPercent = float64(nextCompressedSizeDelta-previousCompressedSizeDelta) / float64(previousCompressedSizeDelta) * 100
		}
		if previousUncompressedSizeDelta > 0 {
			nextUncompressedSizeDeltaPercent = float64(nextUncompressedSizeDelta-previousUncompressedSizeDelta) / float64(previousUncompressedSizeDelta) * 100
		}

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
			// Store delta percentage values
			AuthorsDeltaPercent:      nextAuthorsDeltaPercent,
			CommitsDeltaPercent:      nextCommitsDeltaPercent,
			CompressedDeltaPercent:   nextCompressedSizeDeltaPercent,
			UncompressedDeltaPercent: nextUncompressedSizeDeltaPercent,
		}
		estimates = append(estimates, nextEstimate)

		// Update for next iteration
		previousEstimate = nextEstimate
		previousAuthorsDelta = nextAuthorsDelta
		previousCommitsDelta = nextCommitsDelta
		previousCompressedSizeDelta = nextCompressedSizeDelta
		previousUncompressedSizeDelta = nextUncompressedSizeDelta
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
	fmt.Println("Year       Commits          Δ    T%   LoC    Object size            Δ    T%   LoC   On-disk size            Δ    T%   LoC")
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
		// Current year's estimate references both current totals (^) and estimated (*): show ^*
		yearDisplay = strconv.Itoa(statistics.Year) + "^*"
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
	fmt.Printf("%-6s %12s %10s %5s %3s │%13s %12s %5s %3s │%13s %12s %5s %3s\n",
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
	if recentFetch != "" {
		fmt.Printf("^  Current totals as of the most recent fetch on %s\n", recentFetch[:16])
	} else {
		fmt.Printf("^  Current totals as of Git directory's last modified: %s\n", lastModified[:16])
	}
	if estimationYears > 0 {
		fmt.Println("^* Estimated growth for current year based on year to date deltas (Δ) extrapolated to full year")
		fmt.Println("*  Estimated growth based on current year's estimated delta percentages (Δ%)")
	} else {
		fmt.Println("Growth estimation unavailable: Requires at least 2 years of commit history")
	}
	fmt.Println()
	fmt.Println("Level of Concern (LoC):")
	fmt.Println("○ Unconcerning  ◑ On-road to concerning  ● Concerning")
	fmt.Println("Commits: < 1.5M = ○, >= 1.5M && < 22.5M = ◑, >= 22.5M = ●")
	fmt.Println("Object size: < 10 GB = ○, >= 10 GB && < 160 GB = ◑, >= 160 GB = ●") 
	fmt.Println("On-disk size: < 1 GB = ○, >= 1 GB && < 10 GB = ◑, >= 10 GB = ●")
}
