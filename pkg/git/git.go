package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
)

// CountedObjects keeps track of Git objects that have been counted
var CountedObjects = make(map[string]bool)

// RunGitCommand runs a git command with the given arguments and returns its output
func RunGitCommand(debug bool, args ...string) ([]byte, error) {
	utils.DebugPrint(debug, "git %s", strings.Join(args, " "))
	command := exec.Command("git", args...)
	return command.Output()
}

// GetGitVersion returns the installed git version
func GetGitVersion() string {
	if output, err := RunGitCommand(false, "version"); err == nil {
		return strings.TrimPrefix(strings.TrimSpace(string(output)), "git version ")
	}
	return "Unknown"
}

// GetDefaultBranch detects and returns the default branch name (main, master, etc.)
func GetDefaultBranch() (string, error) {
	// First try to get the default branch from remote origin
	cmd := exec.Command("git", "remote", "show", "origin")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "HEAD branch:") {
				return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "HEAD branch:")), nil
			}
		}
	}

	// If that fails, check common default branch names
	commonBranches := []string{"main", "master"}
	for _, branch := range commonBranches {
		cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
		if cmd.Run() == nil {
			return branch, nil
		}
	}

	// If all else fails, try to get current branch
	cmd = exec.Command("git", "branch", "--show-current")
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output)), nil
	}

	return "", errors.New("could not determine default branch")
}

// GetBranchFiles returns a map of all files in the given branch
func GetBranchFiles(branch string) (map[string]bool, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", branch)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := make(map[string]bool)
	for _, file := range strings.Split(string(output), "\n") {
		if file != "" {
			files[file] = true
		}
	}

	return files, nil
}

// GetGitDirectory gets the path to the .git directory for a repository
func GetGitDirectory(path string) (string, error) {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("repository path does not exist: %s", path)
	}

	// Run git rev-parse to get git directory
	gitDir, err := RunGitCommand(false, "-C", path, "rev-parse", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("not a git repository: %s", path)
	}

	// Get absolute paths for both the git dir and the repository path
	absRepoPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute repository path: %s", err)
	}

	gitDirPath := strings.TrimSpace(string(gitDir))
	if !filepath.IsAbs(gitDirPath) {
		// If git dir is relative, join it with the repository path
		gitDirPath = filepath.Join(absRepoPath, gitDirPath)
	}

	// Convert to absolute path to clean it up
	absPath, err := filepath.Abs(gitDirPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %s", err)
	}

	return absPath, nil
}

// GetLastFetchTime returns the time of the last git fetch by checking FETCH_HEAD
func GetLastFetchTime(gitDir string) string {
	fetchHead := filepath.Join(gitDir, "FETCH_HEAD")

	// Check if FETCH_HEAD exists
	if fetchInformation, err := os.Stat(fetchHead); err == nil {
		return fetchInformation.ModTime().Format("Mon, 02 Jan 2006 15:04 MST")
	}

	return ""
}

// GetGrowthStats calculates repository growth statistics for a given year
func GetGrowthStats(year int, previousGrowthStatistics models.GrowthStatistics, debug bool) (models.GrowthStatistics, error) {
	utils.DebugPrint(debug, "Calculating stats for year %d", year)
	currentStatistics := models.GrowthStatistics{Year: year}
	startTime := time.Now()

	// Build shell command with before and after dates.
	commandString := fmt.Sprintf("git rev-list --objects --all --before %d-01-01 --after %d-12-31 | git cat-file --batch-check='%%(objecttype) %%(objectname) %%(objectsize) %%(objectsize:disk) %%(rest)'", year+1, year-1)
	command := exec.Command(ShellToUse(), "-c", commandString)
	output, err := command.Output()
	if err != nil {
		return currentStatistics, err
	}

	// Prepare a map to collect blob files (keyed by file path).
	blobsMap := make(map[string]models.FileInformation)
	var commitsDelta, treesDelta, blobsDelta int
	var compressedDelta, uncompressedDelta int64
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		objectType := fields[0]
		objectIdentifier := fields[1]
		// Filter out objects already counted
		if CountedObjects[objectIdentifier] {
			continue
		}
		CountedObjects[objectIdentifier] = true

		uncompressedSize, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			continue // Skip invalid size entries
		}
		compressedSize, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			continue // Skip invalid size entries
		}
		compressedDelta += compressedSize
		uncompressedDelta += uncompressedSize

		switch objectType {
		case "commit":
			commitsDelta++
		case "tree":
			treesDelta++
		case "blob":
			blobsDelta++
			// Collect blob if file path available (5th field onward)
			if len(fields) >= 5 {
				filePath := strings.Join(fields[4:], " ")
				filePath = strings.TrimSpace(filePath)
				if filePath != "" {
					if existing, ok := blobsMap[filePath]; ok {
						existing.Blobs++
						existing.CompressedSize += compressedSize
						existing.UncompressedSize += uncompressedSize
						blobsMap[filePath] = existing
					} else {
						blobsMap[filePath] = models.FileInformation{
							Path:             filePath,
							Blobs:            1,
							CompressedSize:   compressedSize,
							UncompressedSize: uncompressedSize,
							// LastChange remains zero as we do not parse it here
						}
					}
				}
			}
		}
	}

	currentStatistics.Commits = previousGrowthStatistics.Commits + commitsDelta
	currentStatistics.Trees = previousGrowthStatistics.Trees + treesDelta
	currentStatistics.Blobs = previousGrowthStatistics.Blobs + blobsDelta
	currentStatistics.Compressed = previousGrowthStatistics.Compressed + compressedDelta
	currentStatistics.Uncompressed = previousGrowthStatistics.Uncompressed + uncompressedDelta
	currentStatistics.RunTime = time.Since(startTime)

	// Convert blobsMap to slice.
	var currentYearBlobs []models.FileInformation
	for _, fileInfo := range blobsMap {
		currentYearBlobs = append(currentYearBlobs, fileInfo)
	}
	// Merge with previousGrowthStatistics largest blobs.
	mergedBlobsMap := make(map[string]models.FileInformation)
	for _, blob := range previousGrowthStatistics.LargestFiles {
		mergedBlobsMap[blob.Path] = blob
	}
	for _, blob := range currentYearBlobs {
		if existing, ok := mergedBlobsMap[blob.Path]; ok {
			existing.Blobs += blob.Blobs
			existing.CompressedSize += blob.CompressedSize
			existing.UncompressedSize += blob.UncompressedSize
			mergedBlobsMap[blob.Path] = existing
		} else {
			mergedBlobsMap[blob.Path] = blob
		}
	}
	var mergedBlobs []models.FileInformation
	for _, blob := range mergedBlobsMap {
		mergedBlobs = append(mergedBlobs, blob)
	}
	currentStatistics.LargestFiles = mergedBlobs

	utils.DebugPrint(debug, "Finished calculating stats for year %d in %v", year, currentStatistics.RunTime)
	return currentStatistics, nil
}

func ShellToUse() string {
	if runtime.GOOS == "windows" {
		return "bash"
	}
	return "sh"
}

// GetContributors returns all commit authors and committers with dates from git history
func GetContributors() ([]string, error) {
	// Execute the git command to get all contributors with their commit dates
	command := exec.Command("git", "log", "--all", "--format=%an|%cn|%cd", "--date=format:%Y")
	output, err := command.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

// contributorEntry stores the name and count for a contributor.
type contributorEntry struct {
	Name  string
	Count int
}

// commitInfo stores information about a commit for rate calculations
type commitInfo struct {
	timestamp time.Time
	isMerge   bool
	isWorkday bool
	author    string
}

// processContributors takes a map of contributor names to counts,
// sorts them, and returns the top N contributors along with the total unique contributor count.
func processContributors(contributors map[string]int, n int, year int) ([][3]string, int) {
	var contributorList []contributorEntry
	for name, count := range contributors {
		contributorList = append(contributorList, contributorEntry{Name: name, Count: count})
	}

	// Sort by commit count (descending) and then by name (ascending, case-insensitive)
	sort.Slice(contributorList, func(i, j int) bool {
		if contributorList[i].Count != contributorList[j].Count {
			return contributorList[i].Count > contributorList[j].Count
		}
		return strings.ToLower(contributorList[i].Name) < strings.ToLower(contributorList[j].Name)
	})

	var topNContributors [][3]string
	for i, contributor := range contributorList {
		if i >= n {
			break
		}
		topNContributors = append(topNContributors, [3]string{
			contributor.Name,
			strconv.Itoa(contributor.Count),
			strconv.Itoa(year),
		})
	}
	return topNContributors, len(contributors)
}

// GetTopCommitAuthors returns the top N commit authors and committers by number of commits, grouped by year
func GetTopCommitAuthors(n int) (map[int][][3]string, map[int]int, map[int]int, map[int][][3]string, map[int]int, map[string]int, map[string]int, error) {
	// Get all commit authors and committers with dates
	lines, err := GetContributors()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	authorsByYear := make(map[int]map[string]int)
	committersByYear := make(map[int]map[string]int)
	totalCommitsByYear := make(map[int]int)

	// Maps to track all unique authors and committers across all years
	allTimeAuthors := make(map[string]int)
	allTimeCommitters := make(map[string]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}

		author := parts[0]
		committer := parts[1]
		yearStr := parts[2]

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			continue
		}

		if _, exists := authorsByYear[year]; !exists {
			authorsByYear[year] = make(map[string]int)
		}
		if _, exists := committersByYear[year]; !exists {
			committersByYear[year] = make(map[string]int)
		}

		authorsByYear[year][author]++
		committersByYear[year][committer]++
		totalCommitsByYear[year]++

		// Track all unique authors and committers across all years
		allTimeAuthors[author]++
		allTimeCommitters[committer]++
	}

	// Convert to result format: map[year] -> sorted authors/committers
	authorResult := make(map[int][][3]string)
	committerResult := make(map[int][][3]string)
	totalAuthorsByYear := make(map[int]int)
	totalCommittersByYear := make(map[int]int)

	// Process authors
	for year, authors := range authorsByYear {
		topAuthors, total := processContributors(authors, n, year)
		authorResult[year] = topAuthors
		totalAuthorsByYear[year] = total
	}

	// Process committers
	for year, committers := range committersByYear {
		topCommitters, total := processContributors(committers, n, year)
		committerResult[year] = topCommitters
		totalCommittersByYear[year] = total
	}

	return authorResult, totalAuthorsByYear, totalCommitsByYear, committerResult, totalCommittersByYear, allTimeAuthors, allTimeCommitters, nil
}

// GetCumulativeUniqueAuthorsByYear returns a map of year -> cumulative unique author count
// along with the final total unique authors across all years.
func GetCumulativeUniqueAuthorsByYear() (map[int]int, int, error) {
	lines, err := GetContributors()
	if err != nil {
		return nil, 0, err
	}
	// authors seen per year (for union later) and global set
	authorsPerYear := make(map[int]map[string]struct{})
	yearsSet := make(map[int]struct{})
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}
		author := parts[0]
		yearStr := parts[2]
		year, convErr := strconv.Atoi(yearStr)
		if convErr != nil {
			continue
		}
		if _, ok := authorsPerYear[year]; !ok {
			authorsPerYear[year] = make(map[string]struct{})
		}
		authorsPerYear[year][author] = struct{}{}
		yearsSet[year] = struct{}{}
	}
	// Sort years
	var years []int
	for y := range yearsSet {
		years = append(years, y)
	}
	sort.Ints(years)
	cumulativeCounts := make(map[int]int)
	cumulativeSet := make(map[string]struct{})
	for _, y := range years {
		for author := range authorsPerYear[y] {
			cumulativeSet[author] = struct{}{}
		}
		cumulativeCounts[y] = len(cumulativeSet)
	}
	return cumulativeCounts, len(cumulativeSet), nil
}

// GetRateOfChanges calculates commit rate statistics for the current branch by year
func GetRateOfChanges() (map[int]models.RateStatistics, string, error) {
	// Get current branch name instead of remote default branch
	cmd := exec.Command("git", "branch", "--show-current")
	branchOutput, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("could not determine current branch: %v", err)
	}
	currentBranch := strings.TrimSpace(string(branchOutput))
	if currentBranch == "" {
		return nil, "", fmt.Errorf("no current branch found")
	}

	// Get all commits from current branch with timestamps, merge info, and authors
	command := exec.Command("git", "log", currentBranch, "--format=%ct|%P|%an", "--reverse")
	output, err := command.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get commit log: %v", err)
	}

	rateStats, err := calculateRateStatistics(string(output))
	return rateStats, currentBranch, err
}

// calculateRateStatistics processes git log output and calculates rate statistics
func calculateRateStatistics(gitLogOutput string) (map[int]models.RateStatistics, error) {
	lines := strings.Split(strings.TrimSpace(gitLogOutput), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	// Parse commits by year
	commitsByYear := make(map[int][]commitInfo)
	var totalCommits int

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}

		// Parse timestamp
		timestampStr := parts[0]
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			continue
		}
		commitTime := time.Unix(timestamp, 0)

		// Check if it's a merge commit (has multiple parents)
		parents := strings.TrimSpace(parts[1])
		isMerge := strings.Contains(parents, " ")

		// Get author name
		author := strings.TrimSpace(parts[2])

		// Check if it's a workday (Monday-Friday)
		weekday := commitTime.Weekday()
		isWorkday := weekday >= time.Monday && weekday <= time.Friday

		year := commitTime.Year()
		commitsByYear[year] = append(commitsByYear[year], commitInfo{
			timestamp: commitTime,
			isMerge:   isMerge,
			isWorkday: isWorkday,
			author:    author,
		})
		totalCommits++
	}

	// Calculate statistics for each year
	ratesByYear := make(map[int]models.RateStatistics)

	for year, commits := range commitsByYear {
		stats := models.RateStatistics{
			Year:              year,
			TotalCommits:      len(commits),
			PercentageOfTotal: float64(len(commits)) / float64(totalCommits) * 100,
		}

		// Calculate merge statistics and count unique authors
		uniqueAuthors := make(map[string]bool)
		for _, commit := range commits {
			// Count unique authors
			uniqueAuthors[commit.author] = true
			
			if commit.isMerge {
				stats.MergeCommits++
			} else {
				stats.DirectCommits++
			}

			if commit.isWorkday {
				stats.WorkdayCommits++
			} else {
				stats.WeekendCommits++
			}
		}

		// Set the number of active authors
		stats.ActiveAuthors = len(uniqueAuthors)

		if stats.TotalCommits > 0 {
			stats.MergeRatio = float64(stats.MergeCommits) / float64(stats.TotalCommits) * 100
			if stats.WeekendCommits > 0 {
				stats.WorkdayWeekendRatio = float64(stats.WorkdayCommits) / float64(stats.WeekendCommits)
			} else {
				stats.WorkdayWeekendRatio = float64(stats.WorkdayCommits)
			}
		}

		// Calculate daily statistics
		dailyCommits := make(map[string]int)
		hourlyCommits := make(map[string]int)
		minutelyCommits := make(map[string]int)

		for _, commit := range commits {
			day := commit.timestamp.Format("2006-01-02")
			hour := commit.timestamp.Format("2006-01-02-15")
			minute := commit.timestamp.Format("2006-01-02-15:04")

			dailyCommits[day]++
			hourlyCommits[hour]++
			minutelyCommits[minute]++
		}

		// Calculate average commits per day
		daysInYear := 365
		if isLeapYear(year) {
			daysInYear = 366
		}
		stats.AverageCommitsPerDay = float64(stats.TotalCommits) / float64(daysInYear)

		// Find busiest day and calculate percentiles
		var dailyCounts []int
		var busiestDay string
		maxDailyCommits := 0

		for day, count := range dailyCommits {
			dailyCounts = append(dailyCounts, count)
			if count > maxDailyCommits {
				maxDailyCommits = count
				busiestDay = day
			}
		}

		if len(dailyCounts) > 0 {
			sort.Ints(dailyCounts)
			stats.DailyPeakP95 = calculatePercentile(dailyCounts, 95)
			stats.DailyPeakP99 = calculatePercentile(dailyCounts, 99)
			stats.DailyPeakP100 = dailyCounts[len(dailyCounts)-1] // Maximum value
			stats.BusiestDay = busiestDay
			stats.BusiestDayCommits = maxDailyCommits
		}

		// Calculate hourly percentiles
		var hourlyCounts []int
		for _, count := range hourlyCommits {
			hourlyCounts = append(hourlyCounts, count)
		}

		if len(hourlyCounts) > 0 {
			sort.Ints(hourlyCounts)
			stats.HourlyPeakP95 = calculatePercentile(hourlyCounts, 95)
			stats.HourlyPeakP99 = calculatePercentile(hourlyCounts, 99)
			stats.HourlyPeakP100 = hourlyCounts[len(hourlyCounts)-1] // Maximum value
		}

		// Calculate minutely percentiles
		var minutelyCounts []int
		for _, count := range minutelyCommits {
			minutelyCounts = append(minutelyCounts, count)
		}

		if len(minutelyCounts) > 0 {
			sort.Ints(minutelyCounts)
			stats.MinutelyPeakP95 = calculatePercentile(minutelyCounts, 95)
			stats.MinutelyPeakP99 = calculatePercentile(minutelyCounts, 99)
			stats.MinutelyPeakP100 = minutelyCounts[len(minutelyCounts)-1] // Maximum value
		}

		ratesByYear[year] = stats
	}

	return ratesByYear, nil
}

// calculatePercentile calculates the nth percentile of a sorted slice
func calculatePercentile(sortedData []int, percentile int) int {
	if len(sortedData) == 0 {
		return 0
	}

	index := float64(percentile) / 100.0 * float64(len(sortedData)-1)
	if index == float64(int(index)) {
		return sortedData[int(index)]
	}

	lower := int(index)
	upper := lower + 1
	if upper >= len(sortedData) {
		return sortedData[len(sortedData)-1]
	}

	weight := index - float64(lower)
	return int(float64(sortedData[lower])*(1-weight) + float64(sortedData[upper])*weight)
}

// isLeapYear checks if a given year is a leap year
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
