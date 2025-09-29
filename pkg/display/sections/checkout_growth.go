package sections

import (
	"fmt"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"sort"
	"strconv"
	"strings"
)

// DisplayCheckoutGrowth displays the checkout growth statistics section
func DisplayCheckoutGrowth(checkoutStatistics map[int]models.CheckoutGrowthStatistics) {
	if len(checkoutStatistics) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("CHECKOUT GROWTH ################################################################################################")
	fmt.Println()
	fmt.Println("Year     Directories    Max depth    Max path length    Files           Total size")
	fmt.Println("----------------------------------------------------------------------------------------------------")

	// Get years and sort them
	var years []int
	for year := range checkoutStatistics {
		years = append(years, year)
	}
	sort.Ints(years)

	// Display each year's statistics
	for _, year := range years {
		stats := checkoutStatistics[year]
		DisplayCheckoutGrowthRow(stats)
	}
}

// DisplayCheckoutGrowthRow displays a single row of checkout growth statistics
func DisplayCheckoutGrowthRow(stats models.CheckoutGrowthStatistics) {
	yearDisplay := strconv.Itoa(stats.Year)
	
	fmt.Printf("%-9s%11s%12d%18d%14s%16s\n",
		yearDisplay,
		utils.FormatNumber(stats.NumberDirectories),
		stats.MaxPathDepth,
		stats.MaxPathLength,
		utils.FormatNumber(stats.NumberFiles),
		strings.TrimSpace(utils.FormatSize(stats.TotalSizeFiles)))
}