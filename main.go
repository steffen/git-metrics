package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"git-metrics/pkg/display"
	"git-metrics/pkg/git"
	"git-metrics/pkg/models"
	"git-metrics/pkg/progress"
	"git-metrics/pkg/requirements"
	"git-metrics/pkg/utils"
)

var debug bool

const (
	UnknownValue = "Unknown"
)

func main() {
	startTime := time.Now()
	repositoryPath := flag.String("r", ".", "Path to git repository")
	flag.StringVar(repositoryPath, "repository", ".", "Path to git repository")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	noProgress := flag.Bool("no-progress", false, "Disable progress output")
	flag.Parse()

	// Set progress visibility based on --no-progress flag
	progress.ShowProgress = !*noProgress

	if !requirements.CheckRequirements() {
		fmt.Println("\nRequirements not met. Please install listed dependencies above.")
		os.Exit(9)
	}

	// Validate and change to repository directory
	if err := git.ValidateRepository(*repositoryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(*repositoryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not change to repository directory: %v\n", err)
		os.Exit(1)
	}

	display.PrintMachineInformation()

	fmt.Println("\nREPOSITORY #####################################################################################")
	fmt.Println()
	absolutePath, _ := filepath.Abs(".")
	fmt.Printf("Path                       %s\n", absolutePath)

	// Remote URL
	fmt.Printf("Remote                     ... fetching\n")
	remoteOutput, err := git.RunGitCommand(debug, "remote", "get-url", "origin")
	remote := UnknownValue
	if err == nil {
		remote = strings.TrimSpace(string(remoteOutput))
	}
	// Replace the fetching line with the final value
	fmt.Printf("\033[1A\033[2KRemote                     %s\n", remote)

	// Most recent fetch
	fmt.Printf("Most recent fetch          ... fetching\n")
	recentFetch := git.GetLastFetchTime()
	fmt.Printf("\033[1A\033[2KMost recent fetch          %s\n", recentFetch)

	// Most recent commit
	fmt.Printf("Most recent commit         ... fetching\n")
	lastHashOutput, err := git.RunGitCommand(debug, "rev-parse", "--short", "HEAD")
	lastCommit := UnknownValue
	if err == nil {
		lastHash := strings.TrimSpace(string(lastHashOutput))
		dateCommand := exec.Command("git", "show", "-s", "--format=%cD", lastHash)
		commandOutput, err := dateCommand.Output()
		if err == nil {
			lastDate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(string(commandOutput)))
			lastCommit = fmt.Sprintf("%s (%s)", lastDate.Format("Mon, 02 Jan 2006"), lastHash)
		}
	}
	fmt.Printf("\033[1A\033[2KMost recent commit         %s\n", lastCommit)

	// First commit and age
	fmt.Printf("First commit               ... fetching\n")
	firstOutput, err := git.RunGitCommand(debug, "rev-list", "--max-parents=0", "HEAD", "--format=%cD")
	firstCommit := UnknownValue
	ageString := UnknownValue
	var firstCommitTime time.Time
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(firstOutput)), "\n")
		type commit struct {
			hash string
			date time.Time
		}
		var commits []commit
		for i := 0; i < len(lines); i += 2 {
			if i+1 >= len(lines) {
				break
			}
			hash := strings.TrimPrefix(lines[i], "commit ")[:6]
			if date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(lines[i+1])); err == nil {
				commits = append(commits, commit{hash: hash, date: date})
			}
		}
		if len(commits) > 0 {
			sort.Slice(commits, func(i, j int) bool {
				return commits[i].date.Before(commits[j].date)
			})
			first := commits[0]
			firstCommitTime = first.date
			firstCommit = fmt.Sprintf("%s (%s)", first.date.Format("Mon, 02 Jan 2006"), first.hash)
			now := time.Now()
			years, months, days := utils.CalculateYearsMonthsDays(first.date, now)
			var parts []string
			if years > 0 {
				parts = append(parts, fmt.Sprintf("%d years", years))
			}
			if months > 0 {
				parts = append(parts, fmt.Sprintf("%d months", months))
			}
			if days > 0 {
				parts = append(parts, fmt.Sprintf("%d days", days))
			}
			ageString = strings.Join(parts, " ")
		}
	}
	fmt.Printf("\033[1A\033[2KFirst commit               %s\n", firstCommit)

	// If there are no commits, exit early
	if firstCommit == UnknownValue {
		fmt.Println("\n\nNo commits found in the repository.")
		os.Exit(2)
	}

	fmt.Printf("Age                        %s\n", ageString)

	// Print growth table header first
	display.PrintGrowthTableHeader()

	// Then calculate growth stats and totals
	var previous models.GrowthStatistics
	var totalStatistics models.GrowthStatistics

	// Estimation
	var estimationTotalDeltaStatistics models.GrowthStatistics
	var estimationYearlyAverage models.GrowthStatistics
	var estimationStartYear = firstCommitTime.Year() + 1
	var estimationEndYear = time.Now().Year() - 1
	var estimationYears = estimationEndYear - estimationStartYear + 1
	var minimumRequiredEstimationYears = 1
	var maximumEstimationYears = 5
	var estimationDisplayYears = 6

	if estimationYears > maximumEstimationYears {
		estimationYears = 5
		estimationStartYear = estimationEndYear - maximumEstimationYears + 1
	}

	yearlyStatistics := make(map[int]models.GrowthStatistics)

	// Start calculation with progress indicator (no newline before progress)
	for year := firstCommitTime.Year(); year <= time.Now().Year(); year++ {
		progress.StartProgress(year, previous, startTime) // Start progress updates
		if cumulativeStatistics, err := git.GetGrowthStats(year, previous, debug); err == nil {
			totalStatistics = cumulativeStatistics
			previousForEstimation := previous
			previous = cumulativeStatistics
			yearlyStatistics[year] = cumulativeStatistics
			progress.CurrentProgress.Statistics = cumulativeStatistics // Update current progress

			if estimationYears < minimumRequiredEstimationYears {
				continue
			}

			if year < estimationStartYear || year > estimationEndYear {
				continue
			}

			estimationTotalDeltaStatistics.Commits += totalStatistics.Commits - previousForEstimation.Commits
			estimationTotalDeltaStatistics.Trees += totalStatistics.Trees - previousForEstimation.Trees
			estimationTotalDeltaStatistics.Blobs += totalStatistics.Blobs - previousForEstimation.Blobs
			estimationTotalDeltaStatistics.Compressed += totalStatistics.Compressed - previousForEstimation.Compressed

			if year == estimationEndYear {
				// Calculate average growth per year for the estimation period
				estimationYearlyAverage = models.GrowthStatistics{
					Commits:    estimationTotalDeltaStatistics.Commits / estimationYears,
					Trees:      estimationTotalDeltaStatistics.Trees / estimationYears,
					Blobs:      estimationTotalDeltaStatistics.Blobs / estimationYears,
					Compressed: estimationTotalDeltaStatistics.Compressed / int64(estimationYears),
				}
			}
		}
	}
	progress.StopProgress() // Stop and clear progress line

	// Save repository information with totals
	repositoryInformation := models.RepositoryInformation{
		Remote:         remote,
		LastCommit:     lastCommit,
		FirstCommit:    firstCommit,
		Age:            ageString,
		FirstDate:      firstCommitTime,
		TotalCommits:   totalStatistics.Commits,
		TotalTrees:     totalStatistics.Trees,
		TotalBlobs:     totalStatistics.Blobs,
		CompressedSize: totalStatistics.Compressed,
	}

	// Print growth table using stored statistics
	previous = models.GrowthStatistics{} // Reset for display
	currentYear := time.Now().Year()

	// Print historical data
	for year := repositoryInformation.FirstDate.Year(); year <= currentYear; year++ {
		if statistics, ok := yearlyStatistics[year]; ok {
			display.PrintGrowthTableRow(statistics, previous, repositoryInformation, false, currentYear)
			previous = statistics
		}
	}

	if estimationYears > 0 {
		// Print separator for projections
		fmt.Println("------------------------------------------------------------------------------------------------")

		// Print 6 years of projections including current year
		lastStatistics := yearlyStatistics[currentYear-1]

		for i := 1; i <= estimationDisplayYears; i++ {
			projected := git.CalculateEstimate(lastStatistics, estimationYearlyAverage)
			display.PrintGrowthTableRow(projected, lastStatistics, repositoryInformation, true, currentYear)
			lastStatistics = projected
		}

		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println()
		fmt.Printf("^ Current totals as of the last fetch on %s\n", recentFetch[:11])
		fmt.Println("* Estimated growth based on the last five years")
		fmt.Println("% Percentages show the increase relative to the current total (^)")
	} else {
		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println("No growth estimation possible: Repository is too young")
	}

	// Use the final statistics for largest files
	largestFiles := totalStatistics.LargestFiles
	// Sort by compressed size descending and take top 10.
	sort.Slice(largestFiles, func(i, j int) bool {
		return largestFiles[i].CompressedSize > largestFiles[j].CompressedSize
	})

	// Calculate total compressed size from all files
	var totalFilesCompressedSize int64
	for _, file := range largestFiles {
		totalFilesCompressedSize += file.CompressedSize
	}

	if len(largestFiles) > 10 {
		largestFiles = largestFiles[:10]
	}
	display.PrintLargestFiles(largestFiles, totalFilesCompressedSize, repositoryInformation.TotalBlobs, len(previous.LargestFiles))

	// New call to display top 10 largest file extensions using accumulated blob data.
	display.PrintTopFileExtensions(previous.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)

	// Get memory statistics for final output
	var memoryStatistics runtime.MemStats
	runtime.ReadMemStats(&memoryStatistics)

	fmt.Printf("\nFinished in %s with a memory footprint of %s.\n",
		utils.FormatDuration(time.Since(startTime)),
		strings.TrimSpace(utils.FormatSize(int64(memoryStatistics.Sys))))
}
