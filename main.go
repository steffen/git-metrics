package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"git-metrics/pkg/display/sections"
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

	// Define flags with pflag for better help formatting
	repositoryPath := pflag.StringP("repository", "r", ".", "Path to git repository")
	showVersion := pflag.Bool("version", false, "Display version information and exit")
	pflag.BoolVar(&debug, "debug", false, "Enable debug output")
	noProgress := pflag.Bool("no-progress", false, "Disable progress indicators")
	showHelp := pflag.BoolP("help", "h", false, "Display this help message")

	pflag.Parse()

	// Show help and exit if help flag is set
	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	// Show version and exit if version flag is set
	if *showVersion {
		fmt.Printf("git-metrics version %s\n", utils.GetGitMetricsVersion())
		os.Exit(0)
	}

	// Set progress visibility based on --no-progress flag and output destination
	// Automatically disable progress when output is piped to a file or redirected
	progress.ShowProgress = !*noProgress && utils.IsTerminal(os.Stdout)

	if !requirements.CheckRequirements() {
		fmt.Println("\nRequirements not met. Please install listed dependencies above.")
		os.Exit(9)
	}

	// Get Git directory and change to repository directory
	gitDir, err := git.GetGitDirectory(*repositoryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(*repositoryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not change to repository directory: %v\n", err)
		os.Exit(1)
	}

	sections.DisplayRunInformation()

	fmt.Println("\nREPOSITORY #############################################################################################################")
	fmt.Println()

	// Get Git directory last modified time
	lastModified := UnknownValue
	if info, err := os.Stat(gitDir); err == nil {
		lastModified = info.ModTime().Format("Mon, 02 Jan 2006 15:04 MST")
	}

	fmt.Printf("Git directory              %s\n", gitDir)

	// Get fetch time before deciding whether to show last modified time
	recentFetch := git.GetLastFetchTime(gitDir)
	if recentFetch == "" {
		fmt.Printf("Last modified              %s\n", lastModified)
	}

	// Remote URL - only show if there is one
	remoteOutput, err := git.RunGitCommand(debug, "remote", "get-url", "origin")
	remote := ""
	if err == nil && len(strings.TrimSpace(string(remoteOutput))) > 0 {
		if progress.ShowProgress {
			fmt.Printf("Remote                     ... fetching\n")
		}
		remote = strings.TrimSpace(string(remoteOutput))
		if progress.ShowProgress {
			fmt.Printf("\033[1A\033[2KRemote                     %s\n", remote)
		} else {
			fmt.Printf("Remote                     %s\n", remote)
		}
	}

	if recentFetch != "" {
		fmt.Printf("Most recent fetch          %s\n", recentFetch)
	}

	// Most recent commit
	if progress.ShowProgress {
		fmt.Printf("Most recent commit         ... fetching\n")
	}
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
	if progress.ShowProgress {
		fmt.Printf("\033[1A\033[2KMost recent commit         %s\n", lastCommit)
	} else {
		fmt.Printf("Most recent commit         %s\n", lastCommit)
	}

	// First commit and age
	if progress.ShowProgress {
		fmt.Printf("First commit               ... fetching\n")
	}
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
	if progress.ShowProgress {
		fmt.Printf("\033[1A\033[2KFirst commit               %s\n", firstCommit)
	} else {
		fmt.Printf("First commit               %s\n", firstCommit)
	}

	// If there are no commits, exit early
	if firstCommit == UnknownValue {
		fmt.Println("\n\nNo commits found in the repository.")
		os.Exit(2)
	}

	fmt.Printf("Age                        %s\n", ageString)

	// Display the section header before data collection
	fmt.Println()
	fmt.Println("HISTORIC AND ESTIMATED GROWTH ##########################################################################################")
	fmt.Println()

	// Print footnotes above table headers
	fmt.Println("T% columns: each year's delta as share of current totals (^)")
	fmt.Println("Δ% columns: change of this year's delta vs previous year's delta")
	fmt.Println()

	// Print table headers before data collection
	fmt.Println("Year     Authors        Δ    T%      Δ%       Commits          Δ    T%      Δ%   On-disk size            Δ    T%      Δ%")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	// Calculate growth stats and totals
	var previous models.GrowthStatistics
	var totalStatistics models.GrowthStatistics

	yearlyStatistics := make(map[int]models.GrowthStatistics)

	// Start calculation with progress indicator (no newline before progress)
	for year := firstCommitTime.Year(); year <= time.Now().Year(); year++ {
		progress.StartProgress(year, previous, startTime) // Start progress updates
		if cumulativeStatistics, err := git.GetGrowthStats(year, previous, debug); err == nil {
			totalStatistics = cumulativeStatistics
			previous = cumulativeStatistics
			yearlyStatistics[year] = cumulativeStatistics
			progress.CurrentProgress.Statistics = cumulativeStatistics // Update current progress
		}
	}
	progress.StopProgress() // Stop and clear progress line

	// Compute cumulative unique authors per year for historic growth
	cumulativeAuthorsByYear, totalAuthors, authorsErr := git.GetCumulativeUniqueAuthorsByYear()
	if authorsErr == nil {
		// Inject authors into yearly statistics
		for year, stats := range yearlyStatistics {
			if authorsCount, ok := cumulativeAuthorsByYear[year]; ok {
				stats.Authors = authorsCount
				yearlyStatistics[year] = stats
			}
		}
	}

	// Save repository information with totals (including authors)
	repositoryInformation := models.RepositoryInformation{
		Remote:         remote,
		LastCommit:     lastCommit,
		FirstCommit:    firstCommit,
		Age:            ageString,
		FirstDate:      firstCommitTime,
		TotalCommits:   totalStatistics.Commits,
		TotalAuthors:   totalAuthors,
		TotalTrees:     totalStatistics.Trees,
		TotalBlobs:     totalStatistics.Blobs,
		CompressedSize: totalStatistics.Compressed,
	}

	// Calculate and store delta, percentage, and delta percentage values
	currentYear := time.Now().Year()
	var previousCumulative models.GrowthStatistics
	var previousDelta models.GrowthStatistics

	// Process each year to calculate and store all derived values
	for year := repositoryInformation.FirstDate.Year(); year <= currentYear; year++ {
		if cumulative, ok := yearlyStatistics[year]; ok {
			// Calculate delta values (year-over-year changes)
			cumulative.AuthorsDelta = cumulative.Authors - previousCumulative.Authors
			cumulative.CommitsDelta = cumulative.Commits - previousCumulative.Commits
			cumulative.TreesDelta = cumulative.Trees - previousCumulative.Trees
			cumulative.BlobsDelta = cumulative.Blobs - previousCumulative.Blobs
			cumulative.CompressedDelta = cumulative.Compressed - previousCumulative.Compressed

			// Calculate percentage of total
			if repositoryInformation.TotalAuthors > 0 {
				cumulative.AuthorsPercent = float64(cumulative.AuthorsDelta) / float64(repositoryInformation.TotalAuthors) * 100
			}
			if repositoryInformation.TotalCommits > 0 {
				cumulative.CommitsPercent = float64(cumulative.CommitsDelta) / float64(repositoryInformation.TotalCommits) * 100
			}
			if repositoryInformation.TotalTrees > 0 {
				cumulative.TreesPercent = float64(cumulative.TreesDelta) / float64(repositoryInformation.TotalTrees) * 100
			}
			if repositoryInformation.TotalBlobs > 0 {
				cumulative.BlobsPercent = float64(cumulative.BlobsDelta) / float64(repositoryInformation.TotalBlobs) * 100
			}
			if repositoryInformation.CompressedSize > 0 {
				cumulative.CompressedPercent = float64(cumulative.CompressedDelta) / float64(repositoryInformation.CompressedSize) * 100
			}

			// Calculate delta percentage changes (Δ%)
			if previousDelta.Year != 0 { // Skip first year
				if previousDelta.AuthorsDelta > 0 {
					cumulative.AuthorsDeltaPercent = float64(cumulative.AuthorsDelta-previousDelta.AuthorsDelta) / float64(previousDelta.AuthorsDelta) * 100
				}
				if previousDelta.CommitsDelta > 0 {
					cumulative.CommitsDeltaPercent = float64(cumulative.CommitsDelta-previousDelta.CommitsDelta) / float64(previousDelta.CommitsDelta) * 100
				}
				if previousDelta.TreesDelta > 0 {
					cumulative.TreesDeltaPercent = float64(cumulative.TreesDelta-previousDelta.TreesDelta) / float64(previousDelta.TreesDelta) * 100
				}
				if previousDelta.BlobsDelta > 0 {
					cumulative.BlobsDeltaPercent = float64(cumulative.BlobsDelta-previousDelta.BlobsDelta) / float64(previousDelta.BlobsDelta) * 100
				}
				if previousDelta.CompressedDelta > 0 {
					cumulative.CompressedDeltaPercent = float64(cumulative.CompressedDelta-previousDelta.CompressedDelta) / float64(previousDelta.CompressedDelta) * 100
				}
			}

			// Store the updated statistics back in the map
			yearlyStatistics[year] = cumulative

			// Update for next iteration
			previousCumulative = cumulative
			previousDelta = cumulative
		}
	}

	// Display unified historic and estimated growth using the new function
	sections.DisplayUnifiedGrowth(yearlyStatistics, repositoryInformation, firstCommitTime, recentFetch, lastModified)

	// Rate of changes analysis - add after historic growth and before largest directories
	if ratesByYear, err := git.GetRateOfChanges(); err == nil && len(ratesByYear) > 0 {
		if defaultBranch, branchErr := git.GetDefaultBranch(); branchErr == nil {
			sections.DisplayRateOfChanges(ratesByYear, defaultBranch)
		}
	}

	// Use the final statistics for largest files
	largestFiles := totalStatistics.LargestFiles
	// Sort by compressed size descending, then by path ascending, and take top 10.
	sort.Slice(largestFiles, func(i, j int) bool {
		if largestFiles[i].CompressedSize != largestFiles[j].CompressedSize {
			return largestFiles[i].CompressedSize > largestFiles[j].CompressedSize
		}
		return largestFiles[i].Path < largestFiles[j].Path
	})

	// Calculate total compressed size from all files
	var totalFilesCompressedSize int64
	for _, file := range largestFiles {
		totalFilesCompressedSize += file.CompressedSize
	}

	if len(largestFiles) > 10 {
		largestFiles = largestFiles[:10]
	}

	// Print largest directories section before largest files
	sections.PrintLargestDirectories(totalStatistics.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)

	sections.PrintLargestFiles(largestFiles, totalFilesCompressedSize, repositoryInformation.TotalBlobs, len(previous.LargestFiles))

	// New call to display top 10 largest file extensions using accumulated blob data.
	sections.PrintTopFileExtensions(previous.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)

	// Print top 3 commit authors and committers per year
	if topAuthorsByYear, totalAuthorsByYear, totalCommitsByYear, topCommittersByYear, totalCommittersByYear, allTimeAuthors, allTimeCommitters, err := git.GetTopCommitAuthors(3); err == nil && len(topAuthorsByYear) > 0 {
		sections.DisplayContributorsWithMostCommits(topAuthorsByYear, totalAuthorsByYear, totalCommitsByYear, topCommittersByYear, totalCommittersByYear, allTimeAuthors, allTimeCommitters)
	}

	// Get memory statistics for final output
	var memoryStatistics runtime.MemStats
	runtime.ReadMemStats(&memoryStatistics)

	fmt.Printf("\nFinished in %s with a memory footprint of %s.\n",
		utils.FormatDuration(time.Since(startTime)),
		strings.TrimSpace(utils.FormatSize(int64(memoryStatistics.Sys))))
}
