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
	currentSizeDelta := currentStats.Compressed - previousStats.Compressed

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
		sizePerDay := float64(currentSizeDelta) / float64(daysPassed)

		predictedCurrentYear = models.GrowthStatistics{
			Year:       currentYear,
			Authors:    previousStats.Authors + int(authorsPerDay*365),
			Commits:    previousStats.Commits + int(commitsPerDay*365),
			Compressed: previousStats.Compressed + int64(sizePerDay*365),
		}
	}

	estimates = append(estimates, predictedCurrentYear)

	// Calculate delta percentage growth rates from current year to apply to future years
	currentYearAuthorsDelta := predictedCurrentYear.Authors - previousStats.Authors
	currentYearCommitsDelta := predictedCurrentYear.Commits - previousStats.Commits
	currentYearSizeDelta := predictedCurrentYear.Compressed - previousStats.Compressed

	// Calculate previous year deltas to determine delta percentage growth rate
	twoYearsAgo, twoYearsExists := yearlyStats[currentYear-2]
	if !twoYearsExists {
		return estimates
	}

	previousYearAuthorsDelta := previousStats.Authors - twoYearsAgo.Authors
	previousYearCommitsDelta := previousStats.Commits - twoYearsAgo.Commits
	previousYearSizeDelta := previousStats.Compressed - twoYearsAgo.Compressed

	// Calculate delta percentage change (Δ% growth rate)
	var authorsDeltaGrowthPercent, commitsDeltaGrowthPercent, sizeDeltaGrowthPercent float64
	if previousYearAuthorsDelta > 0 {
		authorsDeltaGrowthPercent = float64(currentYearAuthorsDelta-previousYearAuthorsDelta) / float64(previousYearAuthorsDelta)
	}
	if previousYearCommitsDelta > 0 {
		commitsDeltaGrowthPercent = float64(currentYearCommitsDelta-previousYearCommitsDelta) / float64(previousYearCommitsDelta)
	}
	if previousYearSizeDelta > 0 {
		sizeDeltaGrowthPercent = float64(currentYearSizeDelta-previousYearSizeDelta) / float64(previousYearSizeDelta)
	}

	// Project future years using delta percentage growth rates
	previousEstimate := predictedCurrentYear
	previousAuthorsDelta := currentYearAuthorsDelta
	previousCommitsDelta := currentYearCommitsDelta
	previousSizeDelta := currentYearSizeDelta

	for futureYear := currentYear + 1; futureYear <= currentYear+5; futureYear++ {
		// Apply delta growth percentages to previous deltas to get new deltas
		nextAuthorsDelta := previousAuthorsDelta + int(float64(previousAuthorsDelta)*authorsDeltaGrowthPercent)
		nextCommitsDelta := previousCommitsDelta + int(float64(previousCommitsDelta)*commitsDeltaGrowthPercent)
		nextSizeDelta := previousSizeDelta + int64(float64(previousSizeDelta)*sizeDeltaGrowthPercent)

		nextEstimate := models.GrowthStatistics{
			Year:       futureYear,
			Authors:    previousEstimate.Authors + nextAuthorsDelta,
			Commits:    previousEstimate.Commits + nextCommitsDelta,
			Compressed: previousEstimate.Compressed + nextSizeDelta,
		}
		estimates = append(estimates, nextEstimate)

		// Update for next iteration
		previousEstimate = nextEstimate
		previousAuthorsDelta = nextAuthorsDelta
		previousCommitsDelta = nextCommitsDelta
		previousSizeDelta = nextSizeDelta
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
	fmt.Println("Year     Authors        Δ     %      Δ%       Commits          Δ     %      Δ%   On-disk size            Δ     %      Δ%")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
}

// PrintGrowthEstimateRow prints a row in the estimated growth table
func PrintGrowthEstimateRow(statistics, previous models.GrowthStatistics, information models.RepositoryInformation, currentYear int) {
	// Calculate delta values for this estimate row
	currentAuthorsDelta := statistics.Authors - previous.Authors
	currentCommitsDelta := statistics.Commits - previous.Commits
	currentSizeDelta := statistics.Compressed - previous.Compressed

	// Calculate percentages for this estimate
	var authorsPercentage, commitsPercentage, compressedPercentage float64
	if information.TotalAuthors > 0 {
		authorsPercentage = float64(currentAuthorsDelta) / float64(information.TotalAuthors) * 100
	}
	if information.TotalCommits > 0 {
		commitsPercentage = float64(currentCommitsDelta) / float64(information.TotalCommits) * 100
	}
	if information.CompressedSize > 0 {
		compressedPercentage = float64(currentSizeDelta) / float64(information.CompressedSize) * 100
	}

	// Use pre-calculated delta percentage values from the statistics struct
	authorsDeltaPercentChange := statistics.AuthorsDeltaPercent
	commitsDeltaPercentChange := statistics.CommitsDeltaPercent
	compressedDeltaPercentChange := statistics.CompressedDeltaPercent

	yearDisplay := strconv.Itoa(statistics.Year) + "*"

	// Helper to format signed integers with thousand separators
	formatSigned := func(v int) string {
		if v >= 0 {
			return "+" + utils.FormatNumber(v)
		}
		return "-" + utils.FormatNumber(-v)
	}

	authorsDeltaDisplay := formatSigned(statistics.Authors - previous.Authors)
	commitsDeltaDisplay := formatSigned(statistics.Commits - previous.Commits)
	// Size delta with explicit sign when positive
	var sizeDeltaDisplay string
	sizeDelta := statistics.Compressed - previous.Compressed
	if sizeDelta >= 0 {
		formatted := utils.FormatSize(sizeDelta)
		sizeDeltaDisplay = "+" + strings.TrimLeft(formatted, " ")
	} else {
		formatted := utils.FormatSize(-sizeDelta)
		sizeDeltaDisplay = "-" + strings.TrimLeft(formatted, " ")
	}

	// Helper to format integer percentage with trailing %
	formatPercent := func(v float64) string { return fmt.Sprintf("%d %%", int(v+0.5)) }

	authorsPercentDisplay := formatPercent(authorsPercentage)
	commitsPercentDisplay := formatPercent(commitsPercentage)
	compressedPercentDisplay := formatPercent(compressedPercentage)

	// Format Δ% values using the same logic as historic growth
	authorsDeltaPercentDisplay := "" // blank for first estimate
	commitsDeltaPercentDisplay := ""
	compressedDeltaPercentDisplay := ""
	if previous.Year != 0 { // Use previous year check instead of previousPrevious
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

	// Print with same formatting as historic growth: % column narrower, Δ% wider (extra left padding)
	fmt.Printf("%-5s %10s %8s %5s %7s │%12s %10s %5s %7s │%13s %12s %5s %7s\n",
		yearDisplay,
		utils.FormatNumber(statistics.Authors), authorsDeltaDisplay, authorsPercentDisplay, authorsDeltaPercentDisplay,
		utils.FormatNumber(statistics.Commits), commitsDeltaDisplay, commitsPercentDisplay, commitsDeltaPercentDisplay,
		utils.FormatSize(statistics.Compressed), sizeDeltaDisplay, compressedPercentDisplay, compressedDeltaPercentDisplay)
}

// DisplayGrowthEstimates handles the complete growth estimation section including calculation and display
func DisplayGrowthEstimates(yearlyStatistics map[int]models.GrowthStatistics, repositoryInformation models.RepositoryInformation, firstCommitTime time.Time, recentFetch string) {
	currentYear := time.Now().Year()

	// Estimation
	var estimationEndYear = time.Now().Year() - 1
	var estimationYears = estimationEndYear - firstCommitTime.Year()

	if estimationYears > 5 {
		estimationYears = 5
	}

	// Show estimated growth table only when estimation period is sufficient
	PrintEstimatedGrowthSectionHeader()

	if estimationYears > 0 {
		PrintEstimatedGrowthTableHeader()

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

		fmt.Println("------------------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		fmt.Println("* Estimated growth based on delta changes from past years")
		fmt.Println("% % columns: each year's delta as share of current totals (^)")
	} else {
		fmt.Println("Growth estimation unavailable: Requires at least 2 years of commit history")
	}
}
