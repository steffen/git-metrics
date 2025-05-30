package sections

import (
	"fmt"
	"sort"
	"strconv"
	"strings" // Added import
	"unicode/utf8"

	"git-metrics/pkg/utils"
)

// Format strings for contributor table rows and formatting
const (
	// Headers and dividers
	formatSectionHeader = "\nAUTHORS & COMMITTERS WITH MOST COMMITS #########################################################"
	formatTableHeader   = "Year    Author                   Commits                Committer                Commits"
	formatTableDivider  = "------------------------------------------------------------------------------------------------"
	formatRowSeparator  = "        ┌───────────────────────────────────────        ┌───────────────────────────────────────"
	// Row format (unified)
	formatRow = "%-8s%-24s%8s  %5.1f%%        %-24s%8s  %5.1f%%\n" // unified row format for all data

	// Summary row formats
	formatTopRow   = "        ├─ Top %-8s      %11s  %5.1f%%        ├─ Top %-8s      %11s  %5.1f%%\n"
	formatOutOfRow = "        └─ Out of %-8s   %11s  %5.1f%%        └─ Out of %-8s   %11s  %5.1f%%\n"

	// Maximum contributor name length
	maxNameLength = 24
)

// truncateContributorName truncates a contributor name to maxNameLength and adds ellipsis if needed
func truncateContributorName(name string) string {
	// Count runes (characters) instead of bytes
	if utf8.RuneCountInString(name) <= maxNameLength {
		return name
	}

	// Convert to runes to handle multi-byte characters correctly
	runes := []rune(name)
	// Take maxNameLength-1 runes and add ellipsis
	return string(runes[:maxNameLength-1]) + "…"
}

// DisplayContributorsWithMostCommits displays the top commit authors and committers by number of commits per year
func DisplayContributorsWithMostCommits(authorsByYear map[int][][3]string, totalAuthorsByYear map[int]int, totalCommitsByYear map[int]int,
	committersByYear map[int][][3]string, totalCommittersByYear map[int]int, allTimeAuthors map[string]int, allTimeCommitters map[string]int) {
	fmt.Println(formatSectionHeader)
	fmt.Println()

	fmt.Println(formatTableHeader)
	fmt.Println(formatTableDivider)
	fmt.Println("")

	// Get years and sort them
	var years []int
	for year := range authorsByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// For all-time totals calculation
	allTimeAuthorCommits := make(map[string]int)
	allTimeCommitterCommits := make(map[string]int)
	var allTimeTotalCommits int

	// Print for each year
	for i, year := range years {
		authors := authorsByYear[year]
		totalAuthors := totalAuthorsByYear[year]
		committers := committersByYear[year]
		totalCommitters := totalCommittersByYear[year]
		totalCommits := totalCommitsByYear[year]

		// Update all-time stats while processing each year
		allTimeTotalCommits += totalCommits

		// Add commits to all-time author totals
		for _, authorData := range authors {
			authorName := authorData[0]
			authorCommits, _ := strconv.Atoi(authorData[1])
			allTimeAuthorCommits[authorName] += authorCommits
		}

		// Add commits to all-time committer totals
		for _, committerData := range committers {
			committerName := committerData[0]
			committerCommits, _ := strconv.Atoi(committerData[1])
			allTimeCommitterCommits[committerName] += committerCommits
		}

		// Calculate top authors total commits
		var topAuthorsTotalCommits int
		for _, author := range authors {
			authorCommits, _ := strconv.Atoi(author[1])
			topAuthorsTotalCommits += authorCommits
		}

		// Calculate top committers total commits
		var topCommittersTotalCommits int
		for _, committer := range committers {
			committerCommits, _ := strconv.Atoi(committer[1])
			topCommittersTotalCommits += committerCommits
		}

		// Determine the max number of rows to print (authors or committers)
		maxRows := len(authors)
		if len(committers) > maxRows {
			maxRows = len(committers)
		}

		// Print each row
		for j := 0; j < maxRows; j++ {
			if j == 0 {
				// First row - print year and first author and committer
				if j < len(authors) {
					authorCommits, _ := strconv.Atoi(authors[j][1])
					authorPercentage := float64(authorCommits) / float64(totalCommits) * 100

					if j < len(committers) {
						committerCommits, _ := strconv.Atoi(committers[j][1])
						committerPercentage := float64(committerCommits) / float64(totalCommits) * 100

						fmt.Printf(formatRow,
							fmt.Sprintf("%d", year),
							truncateContributorName(authors[j][0]), utils.FormatNumber(authorCommits), authorPercentage,
							truncateContributorName(committers[j][0]), utils.FormatNumber(committerCommits), committerPercentage)
					} else {
						// No committer for this row
						fmt.Printf(formatRow,
							fmt.Sprintf("%d", year),
							truncateContributorName(authors[j][0]), utils.FormatNumber(authorCommits), authorPercentage,
							"", "", 0.0)
					}
				} else if j < len(committers) {
					// No author for this row but we have a committer
					committerCommits, _ := strconv.Atoi(committers[j][1])
					committerPercentage := float64(committerCommits) / float64(totalCommits) * 100

					fmt.Printf(formatRow,
						fmt.Sprintf("%d", year),
						"", "", 0.0,
						truncateContributorName(committers[j][0]), utils.FormatNumber(committerCommits), committerPercentage)
				}
			} else {
				// Subsequent rows - just author and committer, no year
				if j < len(authors) {
					authorCommits, _ := strconv.Atoi(authors[j][1])
					authorPercentage := float64(authorCommits) / float64(totalCommits) * 100

					if j < len(committers) {
						committerCommits, _ := strconv.Atoi(committers[j][1])
						committerPercentage := float64(committerCommits) / float64(totalCommits) * 100

						fmt.Printf(formatRow,
							"",
							truncateContributorName(authors[j][0]), utils.FormatNumber(authorCommits), authorPercentage,
							truncateContributorName(committers[j][0]), utils.FormatNumber(committerCommits), committerPercentage)
					} else {
						// No committer for this row
						fmt.Printf(formatRow,
							"",
							truncateContributorName(authors[j][0]), utils.FormatNumber(authorCommits), authorPercentage,
							"", "", 0.0)
					}
				} else if j < len(committers) {
					// No author for this row but we have a committer
					committerCommits, _ := strconv.Atoi(committers[j][1])
					committerPercentage := float64(committerCommits) / float64(totalCommits) * 100

					fmt.Printf(formatRow,
						"",
						"", "", 0.0,
						truncateContributorName(committers[j][0]), utils.FormatNumber(committerCommits), committerPercentage)
				}
			}
		}

		// Add separator before summary rows
		fmt.Println(formatRowSeparator)

		// Print summary rows for authors and committers
		topAuthorsPercentage := float64(topAuthorsTotalCommits) / float64(totalCommits) * 100
		topCommittersPercentage := float64(topCommittersTotalCommits) / float64(totalCommits) * 100

		fmt.Printf(formatTopRow,
			utils.FormatNumber(len(authors)), utils.FormatNumber(topAuthorsTotalCommits), topAuthorsPercentage,
			utils.FormatNumber(len(committers)), utils.FormatNumber(topCommittersTotalCommits), topCommittersPercentage)

		fmt.Printf(formatOutOfRow,
			utils.FormatNumber(totalAuthors), utils.FormatNumber(totalCommits), 100.0,
			utils.FormatNumber(totalCommitters), utils.FormatNumber(totalCommits), 100.0)

		// Add separator after each year except the last one
		if i < len(years)-1 {
			fmt.Println("")
		}
	}

	// After printing all years, add a section for all-time stats
	if len(years) > 0 {
		// Convert author and committer maps to sortable slices
		type contributorStats struct {
			name    string
			commits int
		}

		// Create slices for all-time authors and committers
		var allTimeAuthorsList []contributorStats
		var allTimeCommittersList []contributorStats
		var totalAuthorCount = len(allTimeAuthors)
		var totalCommitterCount = len(allTimeCommitters)

		// Fill authors slice from the map
		for name, commits := range allTimeAuthors {
			allTimeAuthorsList = append(allTimeAuthorsList, contributorStats{name: name, commits: commits})
		}

		// Fill committers slice from the map
		for name, commits := range allTimeCommitters {
			allTimeCommittersList = append(allTimeCommittersList, contributorStats{name: name, commits: commits})
		}

		// Sort by number of commits (descending) and then by name (ascending, case-insensitive) as a secondary criteria
		sort.Slice(allTimeAuthorsList, func(i, j int) bool {
			if allTimeAuthorsList[i].commits != allTimeAuthorsList[j].commits {
				return allTimeAuthorsList[i].commits > allTimeAuthorsList[j].commits
			}
			return strings.ToLower(allTimeAuthorsList[i].name) < strings.ToLower(allTimeAuthorsList[j].name)
		})
		sort.Slice(allTimeCommittersList, func(i, j int) bool {
			if allTimeCommittersList[i].commits != allTimeCommittersList[j].commits {
				return allTimeCommittersList[i].commits > allTimeCommittersList[j].commits
			}
			return strings.ToLower(allTimeCommittersList[i].name) < strings.ToLower(allTimeCommittersList[j].name)
		})

		// Limit to top 3 contributors
		maxDisplayCount := 3
		if len(allTimeAuthorsList) > maxDisplayCount {
			allTimeAuthorsList = allTimeAuthorsList[:maxDisplayCount]
		}
		if len(allTimeCommittersList) > maxDisplayCount {
			allTimeCommittersList = allTimeCommittersList[:maxDisplayCount]
		}

		// Print all-time stats
		fmt.Println("\n------------------------------------------------------------------------------------------------")
		fmt.Println("")

		// Determine the max number of rows for all-time display
		maxRows := len(allTimeAuthorsList)
		if len(allTimeCommittersList) > maxRows {
			maxRows = len(allTimeCommittersList)
		}

		// Calculate top contributors' total commits
		var topAuthorsTotalCommits int
		for _, author := range allTimeAuthorsList {
			topAuthorsTotalCommits += author.commits
		}

		var topCommittersTotalCommits int
		for _, committer := range allTimeCommittersList {
			topCommittersTotalCommits += committer.commits
		}

		// Print each row for all-time stats
		for j := 0; j < maxRows; j++ {
			if j == 0 {
				// First row - print TOTAL and first author and committer
				if j < len(allTimeAuthorsList) {
					authorPercentage := float64(allTimeAuthorsList[j].commits) / float64(allTimeTotalCommits) * 100

					if j < len(allTimeCommittersList) {
						committerPercentage := float64(allTimeCommittersList[j].commits) / float64(allTimeTotalCommits) * 100

						fmt.Printf(formatRow,
							"TOTAL",
							truncateContributorName(allTimeAuthorsList[j].name), utils.FormatNumber(allTimeAuthorsList[j].commits), authorPercentage,
							truncateContributorName(allTimeCommittersList[j].name), utils.FormatNumber(allTimeCommittersList[j].commits), committerPercentage)
					} else {
						// No committer for this row
						fmt.Printf(formatRow,
							"TOTAL",
							truncateContributorName(allTimeAuthorsList[j].name), utils.FormatNumber(allTimeAuthorsList[j].commits), authorPercentage,
							"", "", 0.0)
					}
				} else if j < len(allTimeCommittersList) {
					// No author for this row but we have a committer
					committerPercentage := float64(allTimeCommittersList[j].commits) / float64(allTimeTotalCommits) * 100

					fmt.Printf(formatRow,
						"TOTAL",
						"", "", 0.0,
						truncateContributorName(allTimeCommittersList[j].name), utils.FormatNumber(allTimeCommittersList[j].commits), committerPercentage)
				}
			} else {
				// Subsequent rows - just author and committer, no TOTAL
				if j < len(allTimeAuthorsList) {
					authorPercentage := float64(allTimeAuthorsList[j].commits) / float64(allTimeTotalCommits) * 100

					if j < len(allTimeCommittersList) {
						committerPercentage := float64(allTimeCommittersList[j].commits) / float64(allTimeTotalCommits) * 100

						fmt.Printf(formatRow,
							"",
							truncateContributorName(allTimeAuthorsList[j].name), utils.FormatNumber(allTimeAuthorsList[j].commits), authorPercentage,
							truncateContributorName(allTimeCommittersList[j].name), utils.FormatNumber(allTimeCommittersList[j].commits), committerPercentage)
					} else {
						// No committer for this row
						fmt.Printf(formatRow,
							"",
							truncateContributorName(allTimeAuthorsList[j].name), utils.FormatNumber(allTimeAuthorsList[j].commits), authorPercentage,
							"", "", 0.0)
					}
				} else if j < len(allTimeCommittersList) {
					// No author for this row but we have a committer
					committerPercentage := float64(allTimeCommittersList[j].commits) / float64(allTimeTotalCommits) * 100

					fmt.Printf(formatRow,
						"",
						"", "", 0.0,
						truncateContributorName(allTimeCommittersList[j].name), utils.FormatNumber(allTimeCommittersList[j].commits), committerPercentage)
				}
			}
		}

		// Add separator before summary rows
		fmt.Println(formatRowSeparator)

		// Print summary rows for all-time authors and committers
		topAuthorsPercentage := float64(topAuthorsTotalCommits) / float64(allTimeTotalCommits) * 100
		topCommittersPercentage := float64(topCommittersTotalCommits) / float64(allTimeTotalCommits) * 100

		fmt.Printf(formatTopRow,
			utils.FormatNumber(len(allTimeAuthorsList)), utils.FormatNumber(topAuthorsTotalCommits), topAuthorsPercentage,
			utils.FormatNumber(len(allTimeCommittersList)), utils.FormatNumber(topCommittersTotalCommits), topCommittersPercentage)

		fmt.Printf(formatOutOfRow,
			utils.FormatNumber(totalAuthorCount), utils.FormatNumber(allTimeTotalCommits), 100.0,
			utils.FormatNumber(totalCommitterCount), utils.FormatNumber(allTimeTotalCommits), 100.0)
	}
}
