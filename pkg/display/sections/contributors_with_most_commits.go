package sections

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"git-metrics/pkg/utils"
)

// Format strings for contributor table rows and formatting
const (
	// Authors section headers and dividers
	formatAuthorsHeader      = "\nAUTHORS WITH MOST COMMITS ###################################################################"
	formatAuthorsTableHeader = "Year    Author (#1)      Commits           Author (#2)      Commits           Author (#3)      Commits"
	formatAuthorsDivider     = "------------------------------------------------------------------------------------------------"

	// Committers section headers and dividers
	formatCommittersHeader      = "\nCOMMITTERS WITH MOST COMMITS ################################################################"
	formatCommittersTableHeader = "Year    Committer (#1)   Commits           Committer (#2)   Commits           Committer (#3)   Commits"
	formatCommittersDivider     = "------------------------------------------------------------------------------------------------"

	// Row formats for 3-column layout
	formatThreeColumnRow = "%-8s%-17s%8s  %5.1f%%   %-17s%8s  %5.1f%%   %-17s%8s  %5.1f%%\n"
	formatTwoColumnRow   = "%-8s%-17s%8s  %5.1f%%   %-17s%8s  %5.1f%%\n"
	formatOneColumnRow   = "%-8s%-17s%8s  %5.1f%%\n"

	// Maximum contributor name length
	maxNameLength = 15
)

// contributorStats holds contributor name and commit count
type contributorStats struct {
	name    string
	commits int
}

// truncateContributorName truncates a contributor name to maxNameLength and adds ellipsis if needed
func truncateContributorName(name string) string {
	if utf8.RuneCountInString(name) <= maxNameLength {
		return name
	}

	runes := []rune(name)
	if len(runes) <= maxNameLength-3 {
		return name
	}

	return string(runes[:maxNameLength-3]) + "..."
}

// DisplayContributorsWithMostCommits displays the top commit authors and committers by number of commits per year
func DisplayContributorsWithMostCommits(authorsByYear map[int][][3]string, totalAuthorsByYear map[int]int, totalCommitsByYear map[int]int,
	committersByYear map[int][][3]string, totalCommittersByYear map[int]int, allTimeAuthors map[string]int, allTimeCommitters map[string]int) {

	// Display Authors Section
	displayAuthorsSection(authorsByYear, totalAuthorsByYear, totalCommitsByYear, allTimeAuthors)

	// Display Committers Section
	displayCommittersSection(committersByYear, totalCommittersByYear, totalCommitsByYear, allTimeCommitters)
}

func displayAuthorsSection(authorsByYear map[int][][3]string, totalAuthorsByYear map[int]int, totalCommitsByYear map[int]int, allTimeAuthors map[string]int) {
	fmt.Println(formatAuthorsHeader)
	fmt.Println()
	fmt.Println(formatAuthorsTableHeader)
	fmt.Println(formatAuthorsDivider)

	// Get years and sort them
	var years []int
	for year := range authorsByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// For all-time totals calculation
	allTimeAuthorCommits := make(map[string]int)
	var allTimeTotalCommits int

	// Print for each year
	for _, year := range years {
		authors := authorsByYear[year]
		totalCommits := totalCommitsByYear[year]
		allTimeTotalCommits += totalCommits

		// Add commits to all-time author totals
		for _, authorData := range authors {
			authorName := authorData[0]
			authorCommits, _ := strconv.Atoi(authorData[1])
			allTimeAuthorCommits[authorName] += authorCommits
		}

		// Print year row with up to 3 authors
		displayYearRowAuthors(year, authors, totalCommits)
	}

	// Display all-time authors summary
	if len(years) > 0 {
		displayAllTimeAuthors(allTimeAuthorCommits, allTimeTotalCommits, allTimeAuthors)
	}
}

func displayCommittersSection(committersByYear map[int][][3]string, totalCommittersByYear map[int]int, totalCommitsByYear map[int]int, allTimeCommitters map[string]int) {
	fmt.Println(formatCommittersHeader)
	fmt.Println()
	fmt.Println(formatCommittersTableHeader)
	fmt.Println(formatCommittersDivider)

	// Get years and sort them
	var years []int
	for year := range committersByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// For all-time totals calculation
	allTimeCommitterCommits := make(map[string]int)
	var allTimeTotalCommits int

	// Print for each year
	for _, year := range years {
		committers := committersByYear[year]
		totalCommits := totalCommitsByYear[year]
		allTimeTotalCommits += totalCommits

		// Add commits to all-time committer totals
		for _, committerData := range committers {
			committerName := committerData[0]
			committerCommits, _ := strconv.Atoi(committerData[1])
			allTimeCommitterCommits[committerName] += committerCommits
		}

		// Print year row with up to 3 committers
		displayYearRowCommitters(year, committers, totalCommits)
	}

	// Display all-time committers summary
	if len(years) > 0 {
		displayAllTimeCommitters(allTimeCommitterCommits, allTimeTotalCommits, allTimeCommitters)
	}
}

func displayYearRowAuthors(year int, authors [][3]string, totalCommits int) {
	yearStr := fmt.Sprintf("%d", year)

	// Print the row based on how many authors we actually have
	if len(authors) >= 3 {
		// All 3 columns filled
		author1Commits, _ := strconv.Atoi(authors[0][1])
		author1Percentage := float64(author1Commits) / float64(totalCommits) * 100
		author2Commits, _ := strconv.Atoi(authors[1][1])
		author2Percentage := float64(author2Commits) / float64(totalCommits) * 100
		author3Commits, _ := strconv.Atoi(authors[2][1])
		author3Percentage := float64(author3Commits) / float64(totalCommits) * 100

		fmt.Printf(formatThreeColumnRow,
			yearStr,
			truncateContributorName(authors[0][0]), utils.FormatNumber(author1Commits), author1Percentage,
			truncateContributorName(authors[1][0]), utils.FormatNumber(author2Commits), author2Percentage,
			truncateContributorName(authors[2][0]), utils.FormatNumber(author3Commits), author3Percentage)
	} else if len(authors) == 2 {
		// Only 2 columns filled
		author1Commits, _ := strconv.Atoi(authors[0][1])
		author1Percentage := float64(author1Commits) / float64(totalCommits) * 100
		author2Commits, _ := strconv.Atoi(authors[1][1])
		author2Percentage := float64(author2Commits) / float64(totalCommits) * 100

		fmt.Printf(formatTwoColumnRow,
			yearStr,
			truncateContributorName(authors[0][0]), utils.FormatNumber(author1Commits), author1Percentage,
			truncateContributorName(authors[1][0]), utils.FormatNumber(author2Commits), author2Percentage)
	} else if len(authors) == 1 {
		// Only 1 column filled
		author1Commits, _ := strconv.Atoi(authors[0][1])
		author1Percentage := float64(author1Commits) / float64(totalCommits) * 100

		fmt.Printf(formatOneColumnRow,
			yearStr,
			truncateContributorName(authors[0][0]), utils.FormatNumber(author1Commits), author1Percentage)
	}
}

func displayYearRowCommitters(year int, committers [][3]string, totalCommits int) {
	yearStr := fmt.Sprintf("%d", year)

	// Print the row based on how many committers we actually have
	if len(committers) >= 3 {
		// All 3 columns filled
		committer1Commits, _ := strconv.Atoi(committers[0][1])
		committer1Percentage := float64(committer1Commits) / float64(totalCommits) * 100
		committer2Commits, _ := strconv.Atoi(committers[1][1])
		committer2Percentage := float64(committer2Commits) / float64(totalCommits) * 100
		committer3Commits, _ := strconv.Atoi(committers[2][1])
		committer3Percentage := float64(committer3Commits) / float64(totalCommits) * 100

		fmt.Printf(formatThreeColumnRow,
			yearStr,
			truncateContributorName(committers[0][0]), utils.FormatNumber(committer1Commits), committer1Percentage,
			truncateContributorName(committers[1][0]), utils.FormatNumber(committer2Commits), committer2Percentage,
			truncateContributorName(committers[2][0]), utils.FormatNumber(committer3Commits), committer3Percentage)
	} else if len(committers) == 2 {
		// Only 2 columns filled
		committer1Commits, _ := strconv.Atoi(committers[0][1])
		committer1Percentage := float64(committer1Commits) / float64(totalCommits) * 100
		committer2Commits, _ := strconv.Atoi(committers[1][1])
		committer2Percentage := float64(committer2Commits) / float64(totalCommits) * 100

		fmt.Printf(formatTwoColumnRow,
			yearStr,
			truncateContributorName(committers[0][0]), utils.FormatNumber(committer1Commits), committer1Percentage,
			truncateContributorName(committers[1][0]), utils.FormatNumber(committer2Commits), committer2Percentage)
	} else if len(committers) == 1 {
		// Only 1 column filled
		committer1Commits, _ := strconv.Atoi(committers[0][1])
		committer1Percentage := float64(committer1Commits) / float64(totalCommits) * 100

		fmt.Printf(formatOneColumnRow,
			yearStr,
			truncateContributorName(committers[0][0]), utils.FormatNumber(committer1Commits), committer1Percentage)
	}
}

func displayAllTimeAuthors(allTimeAuthorCommits map[string]int, allTimeTotalCommits int, allTimeAuthors map[string]int) {
	// Create slice for all-time authors
	var allTimeAuthorsList []contributorStats

	// Fill authors slice from the map
	for name, commits := range allTimeAuthors {
		allTimeAuthorsList = append(allTimeAuthorsList, contributorStats{name: name, commits: commits})
	}

	// Sort by number of commits (descending) and then by name (ascending, case-insensitive) as a secondary criteria
	sort.Slice(allTimeAuthorsList, func(i, j int) bool {
		if allTimeAuthorsList[i].commits != allTimeAuthorsList[j].commits {
			return allTimeAuthorsList[i].commits > allTimeAuthorsList[j].commits
		}
		return strings.ToLower(allTimeAuthorsList[i].name) < strings.ToLower(allTimeAuthorsList[j].name)
	})

	// Limit to top 3 contributors
	maxDisplayCount := 3
	if len(allTimeAuthorsList) > maxDisplayCount {
		allTimeAuthorsList = allTimeAuthorsList[:maxDisplayCount]
	}

	// Print all-time stats
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	// Calculate top authors' total commits
	var topAuthorsTotalCommits int
	for _, author := range allTimeAuthorsList {
		topAuthorsTotalCommits += author.commits
	}

	// Print TOTAL row for all-time authors
	displayYearRowAuthorsAllTime("TOTAL", allTimeAuthorsList, allTimeTotalCommits)
}

func displayAllTimeCommitters(allTimeCommitterCommits map[string]int, allTimeTotalCommits int, allTimeCommitters map[string]int) {
	// Create slice for all-time committers
	var allTimeCommittersList []contributorStats

	// Fill committers slice from the map
	for name, commits := range allTimeCommitters {
		allTimeCommittersList = append(allTimeCommittersList, contributorStats{name: name, commits: commits})
	}

	// Sort by number of commits (descending) and then by name (ascending, case-insensitive) as a secondary criteria
	sort.Slice(allTimeCommittersList, func(i, j int) bool {
		if allTimeCommittersList[i].commits != allTimeCommittersList[j].commits {
			return allTimeCommittersList[i].commits > allTimeCommittersList[j].commits
		}
		return strings.ToLower(allTimeCommittersList[i].name) < strings.ToLower(allTimeCommittersList[j].name)
	})

	// Limit to top 3 contributors
	maxDisplayCount := 3
	if len(allTimeCommittersList) > maxDisplayCount {
		allTimeCommittersList = allTimeCommittersList[:maxDisplayCount]
	}

	// Print all-time stats
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	// Calculate top committers' total commits
	var topCommittersTotalCommits int
	for _, committer := range allTimeCommittersList {
		topCommittersTotalCommits += committer.commits
	}

	// Print TOTAL row for all-time committers
	displayYearRowCommittersAllTime("TOTAL", allTimeCommittersList, allTimeTotalCommits)
}

func displayYearRowAuthorsAllTime(yearStr string, authors []contributorStats, totalCommits int) {
	// Print the row based on how many authors we actually have
	if len(authors) >= 3 {
		// All 3 columns filled
		author1Percentage := float64(authors[0].commits) / float64(totalCommits) * 100
		author2Percentage := float64(authors[1].commits) / float64(totalCommits) * 100
		author3Percentage := float64(authors[2].commits) / float64(totalCommits) * 100

		fmt.Printf(formatThreeColumnRow,
			yearStr,
			truncateContributorName(authors[0].name), utils.FormatNumber(authors[0].commits), author1Percentage,
			truncateContributorName(authors[1].name), utils.FormatNumber(authors[1].commits), author2Percentage,
			truncateContributorName(authors[2].name), utils.FormatNumber(authors[2].commits), author3Percentage)
	} else if len(authors) == 2 {
		// Only 2 columns filled
		author1Percentage := float64(authors[0].commits) / float64(totalCommits) * 100
		author2Percentage := float64(authors[1].commits) / float64(totalCommits) * 100

		fmt.Printf(formatTwoColumnRow,
			yearStr,
			truncateContributorName(authors[0].name), utils.FormatNumber(authors[0].commits), author1Percentage,
			truncateContributorName(authors[1].name), utils.FormatNumber(authors[1].commits), author2Percentage)
	} else if len(authors) == 1 {
		// Only 1 column filled
		author1Percentage := float64(authors[0].commits) / float64(totalCommits) * 100

		fmt.Printf(formatOneColumnRow,
			yearStr,
			truncateContributorName(authors[0].name), utils.FormatNumber(authors[0].commits), author1Percentage)
	}
}

func displayYearRowCommittersAllTime(yearStr string, committers []contributorStats, totalCommits int) {
	// Print the row based on how many committers we actually have
	if len(committers) >= 3 {
		// All 3 columns filled
		committer1Percentage := float64(committers[0].commits) / float64(totalCommits) * 100
		committer2Percentage := float64(committers[1].commits) / float64(totalCommits) * 100
		committer3Percentage := float64(committers[2].commits) / float64(totalCommits) * 100

		fmt.Printf(formatThreeColumnRow,
			yearStr,
			truncateContributorName(committers[0].name), utils.FormatNumber(committers[0].commits), committer1Percentage,
			truncateContributorName(committers[1].name), utils.FormatNumber(committers[1].commits), committer2Percentage,
			truncateContributorName(committers[2].name), utils.FormatNumber(committers[2].commits), committer3Percentage)
	} else if len(committers) == 2 {
		// Only 2 columns filled
		committer1Percentage := float64(committers[0].commits) / float64(totalCommits) * 100
		committer2Percentage := float64(committers[1].commits) / float64(totalCommits) * 100

		fmt.Printf(formatTwoColumnRow,
			yearStr,
			truncateContributorName(committers[0].name), utils.FormatNumber(committers[0].commits), committer1Percentage,
			truncateContributorName(committers[1].name), utils.FormatNumber(committers[1].commits), committer2Percentage)
	} else if len(committers) == 1 {
		// Only 1 column filled
		committer1Percentage := float64(committers[0].commits) / float64(totalCommits) * 100

		fmt.Printf(formatOneColumnRow,
			yearStr,
			truncateContributorName(committers[0].name), utils.FormatNumber(committers[0].commits), committer1Percentage)
	}
}
