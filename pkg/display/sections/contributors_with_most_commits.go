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
	// Contributor growth section headers and dividers
	formatGrowthHeader      = "\nCONTRIBUTOR GROWTH #########################################################################"
	formatGrowthTableHeader = "Year        Commits   Authors  Commits per Author  Committers  Commits per Committer"
	formatGrowthDivider     = "------------------------------------------------------------------------------------------------"

	// Authors section headers and dividers
	formatAuthorsHeader      = "\nAUTHORS WITH MOST COMMITS ######################################################################"
	formatAuthorsTableHeader = "Year     Author (#1)    Commits        Author (#2)    Commits        Author (#3)    Commits"
	formatAuthorsDivider     = "------------------------------------------------------------------------------------------------"

	// Committers section headers and dividers
	formatCommittersHeader      = "\nCOMMITTERS WITH MOST COMMITS ###################################################################"
	formatCommittersTableHeader = "Year     Committer (#1) Commits        Committer (#2) Commits        Committer (#3) Commits"
	formatCommittersDivider     = "------------------------------------------------------------------------------------------------"

	// Row formats for 3-column layout
	formatThreeColumnRow = "%-6s │ %-14s%8s %3.0f%% │ %-14s%8s %3.0f%% │ %-14s%8s %3.0f%%\n"
	formatTwoColumnRow   = "%-6s │ %-14s%8s %3.0f%% │ %-14s%8s %3.0f%%\n"
	formatOneColumnRow   = "%-6s │ %-14s%8s %3.0f%%\n"

	// Maximum contributor name length
	maxNameLength = 14
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
	return string(runes[:maxNameLength-3]) + "..."
}

// DisplayContributorsWithMostCommits displays the top commit authors and committers by number of commits per year
func DisplayContributorsWithMostCommits(authorsByYear map[int][][3]string, totalAuthorsByYear map[int]int, totalCommitsByYear map[int]int,
	committersByYear map[int][][3]string, totalCommittersByYear map[int]int, allTimeAuthors map[string]int, allTimeCommitters map[string]int) {

	// Display Contributor Growth Section first
	displayContributorGrowthSection(totalAuthorsByYear, totalCommitsByYear, totalCommittersByYear)

	// Display Authors Section
	displayAuthorsSection(authorsByYear, totalAuthorsByYear, totalCommitsByYear, allTimeAuthors)

	// Display Committers Section
	displayCommittersSection(committersByYear, totalCommittersByYear, totalCommitsByYear, allTimeCommitters)
}

func displayContributorGrowthSection(totalAuthorsByYear map[int]int, totalCommitsByYear map[int]int, totalCommittersByYear map[int]int) {
	fmt.Println(formatGrowthHeader)
	fmt.Println()
	fmt.Println(formatGrowthTableHeader)
	fmt.Println(formatGrowthDivider)

	// Get years and sort them
	var years []int
	for year := range totalCommitsByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// Print growth data for each year
	for _, year := range years {
		commits := totalCommitsByYear[year]
		authors := totalAuthorsByYear[year]
		committers := totalCommittersByYear[year]
		
		var commitsPerAuthor float64
		if authors > 0 {
			commitsPerAuthor = float64(commits) / float64(authors)
		}
		
		var commitsPerCommitter float64
		if committers > 0 {
			commitsPerCommitter = float64(commits) / float64(committers)
		}

		fmt.Printf("%-11d%9s%9d%18.1f%12d%19.1f\n",
			year,
			utils.FormatNumber(commits),
			authors,
			commitsPerAuthor,
			committers,
			commitsPerCommitter)
	}
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

// displayContributorRow displays a row of contributors with their commit counts and percentages
func displayContributorRow(yearStr string, contributors [][3]string, totalCommits int) {
	// Print the row based on how many contributors we actually have
	if len(contributors) >= 3 {
		// All 3 columns filled
		contributor1Commits, _ := strconv.Atoi(contributors[0][1])
		contributor1Percentage := float64(contributor1Commits) / float64(totalCommits) * 100
		contributor2Commits, _ := strconv.Atoi(contributors[1][1])
		contributor2Percentage := float64(contributor2Commits) / float64(totalCommits) * 100
		contributor3Commits, _ := strconv.Atoi(contributors[2][1])
		contributor3Percentage := float64(contributor3Commits) / float64(totalCommits) * 100

		fmt.Printf(formatThreeColumnRow,
			yearStr,
			truncateContributorName(contributors[0][0]), utils.FormatNumber(contributor1Commits), contributor1Percentage,
			truncateContributorName(contributors[1][0]), utils.FormatNumber(contributor2Commits), contributor2Percentage,
			truncateContributorName(contributors[2][0]), utils.FormatNumber(contributor3Commits), contributor3Percentage)
	} else if len(contributors) == 2 {
		// Only 2 columns filled
		contributor1Commits, _ := strconv.Atoi(contributors[0][1])
		contributor1Percentage := float64(contributor1Commits) / float64(totalCommits) * 100
		contributor2Commits, _ := strconv.Atoi(contributors[1][1])
		contributor2Percentage := float64(contributor2Commits) / float64(totalCommits) * 100

		fmt.Printf(formatTwoColumnRow,
			yearStr,
			truncateContributorName(contributors[0][0]), utils.FormatNumber(contributor1Commits), contributor1Percentage,
			truncateContributorName(contributors[1][0]), utils.FormatNumber(contributor2Commits), contributor2Percentage)
	} else if len(contributors) == 1 {
		// Only 1 column filled
		contributor1Commits, _ := strconv.Atoi(contributors[0][1])
		contributor1Percentage := float64(contributor1Commits) / float64(totalCommits) * 100

		fmt.Printf(formatOneColumnRow,
			yearStr,
			truncateContributorName(contributors[0][0]), utils.FormatNumber(contributor1Commits), contributor1Percentage)
	}
}

// displayContributorRowAllTime displays a row of contributors with their commit counts and percentages for all-time stats
func displayContributorRowAllTime(yearStr string, contributors []contributorStats, totalCommits int) {
	// Print the row based on how many contributors we actually have
	if len(contributors) >= 3 {
		// All 3 columns filled
		contributor1Percentage := float64(contributors[0].commits) / float64(totalCommits) * 100
		contributor2Percentage := float64(contributors[1].commits) / float64(totalCommits) * 100
		contributor3Percentage := float64(contributors[2].commits) / float64(totalCommits) * 100

		fmt.Printf(formatThreeColumnRow,
			yearStr,
			truncateContributorName(contributors[0].name), utils.FormatNumber(contributors[0].commits), contributor1Percentage,
			truncateContributorName(contributors[1].name), utils.FormatNumber(contributors[1].commits), contributor2Percentage,
			truncateContributorName(contributors[2].name), utils.FormatNumber(contributors[2].commits), contributor3Percentage)
	} else if len(contributors) == 2 {
		// Only 2 columns filled
		contributor1Percentage := float64(contributors[0].commits) / float64(totalCommits) * 100
		contributor2Percentage := float64(contributors[1].commits) / float64(totalCommits) * 100

		fmt.Printf(formatTwoColumnRow,
			yearStr,
			truncateContributorName(contributors[0].name), utils.FormatNumber(contributors[0].commits), contributor1Percentage,
			truncateContributorName(contributors[1].name), utils.FormatNumber(contributors[1].commits), contributor2Percentage)
	} else if len(contributors) == 1 {
		// Only 1 column filled
		contributor1Percentage := float64(contributors[0].commits) / float64(totalCommits) * 100

		fmt.Printf(formatOneColumnRow,
			yearStr,
			truncateContributorName(contributors[0].name), utils.FormatNumber(contributors[0].commits), contributor1Percentage)
	}
}

func displayYearRowAuthors(year int, authors [][3]string, totalCommits int) {
	yearStr := fmt.Sprintf("%d", year)
	displayContributorRow(yearStr, authors, totalCommits)
}

func displayYearRowCommitters(year int, committers [][3]string, totalCommits int) {
	yearStr := fmt.Sprintf("%d", year)
	displayContributorRow(yearStr, committers, totalCommits)
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
	fmt.Println("------------------------------------------------------------------------------------------------")

	// Print total row for all-time authors
	displayYearRowAuthorsAllTime("Total", allTimeAuthorsList, allTimeTotalCommits)
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
	fmt.Println("------------------------------------------------------------------------------------------------")

	// Print total row for all-time committers
	displayYearRowCommittersAllTime("Total", allTimeCommittersList, allTimeTotalCommits)
}

func displayYearRowAuthorsAllTime(yearStr string, authors []contributorStats, totalCommits int) {
	displayContributorRowAllTime(yearStr, authors, totalCommits)
}

func displayYearRowCommittersAllTime(yearStr string, committers []contributorStats, totalCommits int) {
	displayContributorRowAllTime(yearStr, committers, totalCommits)
}
