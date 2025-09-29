package sections

import (
	"fmt"
	"sort"

	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
)

// DisplayRateOfChanges displays commit rate statistics for the default branch
func DisplayRateOfChanges(ratesByYear map[int]models.RateStatistics, defaultBranch string) {
	if len(ratesByYear) == 0 {
		return
	}

	fmt.Println("\nRATE OF CHANGES #######################################################################################################")
	fmt.Printf("\nCommits to default branch (%s)\n\n", defaultBranch)

	// Table header with subcolumns
	fmt.Println("              Commits     Active                   Peak per day            Peak per hour          Peak per minute      ")
	fmt.Println("Year         per year    Authors                 P95    P99   P100        P95    P99   P100        P95    P99   P100    ")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	// Sort years
	var years []int
	for year := range ratesByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// Display statistics for each year
	for _, year := range years {
		stats := ratesByYear[year]

		fmt.Printf("%-4d      %11s    %7d           │ %6d %6d %6d   │ %6d %6d %6d   │ %6d %6d %6d\n",
			stats.Year,
			utils.FormatNumber(stats.TotalCommits),
			stats.ActiveAuthors,
			stats.DailyPeakP95, stats.DailyPeakP99, stats.DailyPeakP100,
			stats.HourlyPeakP95, stats.HourlyPeakP99, stats.HourlyPeakP100,
			stats.MinutelyPeakP95, stats.MinutelyPeakP99, stats.MinutelyPeakP100)
	}
}
