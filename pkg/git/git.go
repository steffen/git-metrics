package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// ValidateRepository checks if the given path is a valid git repository
func ValidateRepository(path string) error {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", path)
	}

	// Check if it's a git repository
	gitDirectory := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDirectory); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", path)
	}

	return nil
}

// GetLastFetchTime returns the time of the last git fetch
func GetLastFetchTime() string {
	fetchHead := filepath.Join(".git", "FETCH_HEAD")
	packDirectory := filepath.Join(".git", "objects", "pack")

	// Check if repository has ever been fetched
	fetchInformation, fetchError := os.Stat(fetchHead)
	if fetchError == nil {
		// Has been fetched at least once
		return fetchInformation.ModTime().Format("Mon, 02 Jan 2006 15:04 MST")
	}

	// No fetch found, try to get clone time from pack directory
	if packInformation, err := os.Stat(packDirectory); err == nil {
		return packInformation.ModTime().Format("Mon, 02 Jan 2006 15:04 MST")
	}

	return "Unknown"
}

// GetGrowthStats calculates repository growth statistics for a given year
func GetGrowthStats(year int, previousGrowthStatistics models.GrowthStatistics, debug bool) (models.GrowthStatistics, error) {
	utils.DebugPrint(debug, "Calculating stats for year %d", year)
	currentStatistics := models.GrowthStatistics{Year: year}
	startTime := time.Now()

	// Build shell command with before and after dates.
	commandString := fmt.Sprintf("git rev-list --objects --all --before %d-01-01 --after %d-12-31 | git cat-file --batch-check='%%(objecttype) %%(objectname) %%(objectsize:disk) %%(rest)'", year+1, year-1)
	command := exec.Command("sh", "-c", commandString)
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
