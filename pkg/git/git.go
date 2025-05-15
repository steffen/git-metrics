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
	commandString := fmt.Sprintf("git rev-list --objects --all --before %d-01-01 --after %d-12-31 | git cat-file --batch-check='%%(objecttype) %%(objectname) %%(objectsize:disk) %%(rest)'", year+1, year-1)
	command := exec.Command(ShellToUse(), "-c", commandString)
	output, err := command.Output()
	if err != nil {
		return currentStatistics, err
	}

	// Prepare a map to collect blob files (keyed by file path).
	blobsMap := make(map[string]models.FileInformation)
	var commitsDelta, treesDelta, blobsDelta int
	var compressedDelta int64
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		objectType := fields[0]
		objectIdentifier := fields[1]
		// Filter out objects already counted
		if CountedObjects[objectIdentifier] {
			continue
		}
		CountedObjects[objectIdentifier] = true

		size, _ := strconv.ParseInt(fields[2], 10, 64)
		compressedDelta += size

		switch objectType {
		case "commit":
			commitsDelta++
		case "tree":
			treesDelta++
		case "blob":
			blobsDelta++
			// Collect blob if file path available (4th field onward)
			if len(fields) >= 4 {
				filePath := strings.Join(fields[3:], " ")
				filePath = strings.TrimSpace(filePath)
				if filePath != "" {
					if existing, ok := blobsMap[filePath]; ok {
						existing.Blobs++
						existing.CompressedSize += size
						blobsMap[filePath] = existing
					} else {
						blobsMap[filePath] = models.FileInformation{
							Path:           filePath,
							Blobs:          1,
							CompressedSize: size,
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

// CalculateEstimate calculates estimated future growth based on current stats and average growth
func CalculateEstimate(current models.GrowthStatistics, average models.GrowthStatistics) models.GrowthStatistics {
	return models.GrowthStatistics{
		Year:         current.Year + 1,
		Commits:      current.Commits + average.Commits,
		Trees:        current.Trees + average.Trees,
		Blobs:        current.Blobs + average.Blobs,
		Compressed:   current.Compressed + average.Compressed,
		LargestFiles: []models.FileInformation{}, // No estimate for largest files
	}
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
