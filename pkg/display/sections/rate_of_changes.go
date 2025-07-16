package sections

import (
	"fmt"
	"sort"
	"time"

	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
)

// DisplayRateOfChanges displays commit rate statistics for the default branch
func DisplayRateOfChanges(ratesByYear map[int]models.RateStatistics, defaultBranch string) {
	if len(ratesByYear) == 0 {
		return
	}

	fmt.Println("\nRATE OF CHANGES ################################################################################")
	fmt.Printf("\nCommits to default branch (%s)\n\n", defaultBranch)

	// Table header
	fmt.Println("                                   Average             Daily peak             Daily peak")
	fmt.Println("            Commits                commits                commits                commits")
	fmt.Println("Year       per year                per day         per hour (P95)          per min (P95)")
	fmt.Println("------------------------------------------------------------------------------------------------")

	// Sort years
	var years []int
	for year := range ratesByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// Display statistics for each year
	for _, year := range years {
		stats := ratesByYear[year]

		// Calculate percentage indicators for visual representation
		dailyPercentage := calculateDailyPercentage(stats, ratesByYear)
		hourlyPercentage := calculateHourlyPercentage(stats, ratesByYear)
		minutelyPercentage := calculateMinutelyPercentage(stats, ratesByYear)

		fmt.Printf("%-4d    %8s   %3.0f%%    %8.0f   %3.0f%%    %8d   %3.0f%%    %8.1f   %3.0f%%\n",
			stats.Year,
			utils.FormatNumber(stats.TotalCommits), stats.PercentageOfTotal,
			stats.AverageCommitsPerDay, dailyPercentage,
			stats.HourlyPeakP95, hourlyPercentage,
			stats.MinutelyPeakP95, minutelyPercentage)
	}

	// Add workflow insights section
	fmt.Println()
	displayWorkflowInsights(ratesByYear, years)
}

// calculateDailyPercentage calculates the percentage for daily average commits
func calculateDailyPercentage(current models.RateStatistics, allStats map[int]models.RateStatistics) float64 {
	if current.AverageCommitsPerDay == 0 {
		return 0
	}

	var maxDaily float64
	for _, stats := range allStats {
		if stats.AverageCommitsPerDay > maxDaily {
			maxDaily = stats.AverageCommitsPerDay
		}
	}

	if maxDaily == 0 {
		return 0
	}

	return (current.AverageCommitsPerDay / maxDaily) * 100
}

// calculateHourlyPercentage calculates the percentage for hourly peak commits
func calculateHourlyPercentage(current models.RateStatistics, allStats map[int]models.RateStatistics) float64 {
	if current.HourlyPeakP95 == 0 {
		return 0
	}

	var maxHourly int
	for _, stats := range allStats {
		if stats.HourlyPeakP95 > maxHourly {
			maxHourly = stats.HourlyPeakP95
		}
	}

	if maxHourly == 0 {
		return 0
	}

	return (float64(current.HourlyPeakP95) / float64(maxHourly)) * 100
}

// calculateMinutelyPercentage calculates the percentage for minutely peak commits
func calculateMinutelyPercentage(current models.RateStatistics, allStats map[int]models.RateStatistics) float64 {
	if current.MinutelyPeakP95 == 0 {
		return 0
	}

	var maxMinutely float64
	for _, stats := range allStats {
		if stats.MinutelyPeakP95 > maxMinutely {
			maxMinutely = stats.MinutelyPeakP95
		}
	}

	if maxMinutely == 0 {
		return 0
	}

	return (current.MinutelyPeakP95 / maxMinutely) * 100
}

// displayWorkflowInsights shows additional insights about development patterns
func displayWorkflowInsights(ratesByYear map[int]models.RateStatistics, sortedYears []int) {
	if len(sortedYears) == 0 {
		return
	}

	// Find overall patterns
	var busiestYear models.RateStatistics
	var maxCommitsPerYear int
	var totalMergeCommits, totalDirectCommits int
	var allBusiestDays []string
	var maxBusiestDayCommits int

	for _, year := range sortedYears {
		stats := ratesByYear[year]

		if stats.TotalCommits > maxCommitsPerYear {
			maxCommitsPerYear = stats.TotalCommits
			busiestYear = stats
		}

		totalMergeCommits += stats.MergeCommits
		totalDirectCommits += stats.DirectCommits

		if stats.BusiestDayCommits > maxBusiestDayCommits {
			maxBusiestDayCommits = stats.BusiestDayCommits
			allBusiestDays = []string{stats.BusiestDay}
		} else if stats.BusiestDayCommits == maxBusiestDayCommits && stats.BusiestDay != "" {
			allBusiestDays = append(allBusiestDays, stats.BusiestDay)
		}
	}

	fmt.Println("Workflow insights:")

	// Display busiest day
	if len(allBusiestDays) > 0 && maxBusiestDayCommits > 0 {
		busiestDay := allBusiestDays[0]
		if parsedDate, err := time.Parse("2006-01-02", busiestDay); err == nil {
			busiestDay = parsedDate.Format("2006-01-02")
		}
		fmt.Printf("• Busiest day: %s (%s commits)\n", busiestDay, utils.FormatNumber(maxBusiestDayCommits))
	}

	// Display busiest year
	if busiestYear.Year != 0 {
		fmt.Printf("• Most active year: %d (%s commits, %.1f avg/day)\n",
			busiestYear.Year,
			utils.FormatNumber(busiestYear.TotalCommits),
			busiestYear.AverageCommitsPerDay)
	}

	// Display merge ratio
	totalCommits := totalMergeCommits + totalDirectCommits
	if totalCommits > 0 {
		mergeRatio := float64(totalMergeCommits) / float64(totalCommits) * 100
		var workflowType string
		switch {
		case mergeRatio < 10:
			workflowType = "suggests direct commits workflow"
		case mergeRatio < 30:
			workflowType = "suggests mixed workflow"
		default:
			workflowType = "suggests feature branch workflow"
		}
		fmt.Printf("• Merge ratio: %.0f%% (%s)\n", mergeRatio, workflowType)
	}

	// Display weekend activity pattern
	var totalWorkdayCommits, totalWeekendCommits int
	for _, year := range sortedYears {
		stats := ratesByYear[year]
		totalWorkdayCommits += stats.WorkdayCommits
		totalWeekendCommits += stats.WeekendCommits
	}

	if totalWorkdayCommits+totalWeekendCommits > 0 {
		weekendPercentage := float64(totalWeekendCommits) / float64(totalWorkdayCommits+totalWeekendCommits) * 100
		var activityLevel string
		switch {
		case weekendPercentage < 10:
			activityLevel = "low after-hours development"
		case weekendPercentage < 25:
			activityLevel = "moderate after-hours development"
		default:
			activityLevel = "high after-hours development"
		}
		fmt.Printf("• Weekend activity: %.0f%% (%s)\n", weekendPercentage, activityLevel)
	}

	// Display development trend
	if len(sortedYears) >= 2 {
		recentYears := sortedYears[len(sortedYears)-2:]
		oldStats := ratesByYear[recentYears[0]]
		newStats := ratesByYear[recentYears[1]]

		if oldStats.TotalCommits > 0 {
			trendPercentage := float64(newStats.TotalCommits-oldStats.TotalCommits) / float64(oldStats.TotalCommits) * 100
			var trendIndicator string
			var trendDescription string

			switch {
			case trendPercentage > 20:
				trendIndicator = "↗"
				trendDescription = "increasing"
			case trendPercentage < -20:
				trendIndicator = "↘"
				trendDescription = "decreasing"
			default:
				trendIndicator = "→"
				trendDescription = "stable"
			}

			fmt.Printf("• Development trend: %s %s (%.0f%% change)\n",
				trendIndicator, trendDescription, trendPercentage)
		}
	}
}
