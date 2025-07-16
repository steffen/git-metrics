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

	// Table header with subcolumns
	fmt.Println("                                    Daily peak                  Hourly peak                 Minutely peak")
	fmt.Println("            Commits")
	fmt.Println("Year       per year           P95   P99  P100           P95   P99  P100           P95   P99  P100")
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

		fmt.Printf("%-4d    %8s       %5d %5d %5d       %5d %5d %5d       %5.1f %5.1f %5.1f\n",
			stats.Year,
			utils.FormatNumber(stats.TotalCommits),
			stats.DailyPeakP95, stats.DailyPeakP99, stats.DailyPeakP100,
			stats.HourlyPeakP95, stats.HourlyPeakP99, stats.HourlyPeakP100,
			stats.MinutelyPeakP95, stats.MinutelyPeakP99, stats.MinutelyPeakP100)
	}

	// Add workflow insights section
	fmt.Println()
	displayWorkflowInsights(ratesByYear, years)
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
