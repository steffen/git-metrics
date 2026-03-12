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

	// Start spinner for fetching repository information
	stopRepositorySpinner := progress.StartSimpleSpinner()

	// Remote URL - only show if there is one
	remoteOutput, err := git.RunGitCommand(debug, "remote", "get-url", "origin")
	remote := ""
	if err == nil && len(strings.TrimSpace(string(remoteOutput))) > 0 {
		remote = strings.TrimSpace(string(remoteOutput))
	}

	// Get fetch time
	recentFetch := git.GetLastFetchTime(gitDir)

	// Most recent commit
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

	// First commit and age
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

	stopRepositorySpinner()

	// Print all repository information
	if remote != "" {
		fmt.Printf("Remote                     %s\n", remote)
	}
	if recentFetch == "" {
		fmt.Printf("Last modified              %s\n", lastModified)
	}
	if recentFetch != "" {
		fmt.Printf("Most recent fetch          %s\n", recentFetch)
	}
	fmt.Printf("Most recent commit         %s\n", lastCommit)
	fmt.Printf("First commit               %s\n", firstCommit)

	// If there are no commits, exit early
	if firstCommit == UnknownValue {
		fmt.Println("\n\nNo commits found in the repository.")
		os.Exit(2)
	}

	fmt.Printf("Age                        %s\n", ageString)

	// Display the section header before data collection
	fmt.Println()
	fmt.Println("HISTORIC & ESTIMATED GROWTH ############################################################################################")
	fmt.Println()

	// Calculate growth stats and totals
	var previous models.GrowthStatistics
	var totalStatistics models.GrowthStatistics

	yearlyStatistics := make(map[int]models.GrowthStatistics)
	currentYear := time.Now().Year()

	// Track previous cumulative for computing deltas during progressive output
	var interimPreviousCumulative models.GrowthStatistics

	// Track number of lines printed for preview (to clear later)
	var previewLinesCount int

	// Print table headers before data collection
	if progress.ShowProgress {
		fmt.Println("Year          Commits          Δ     %   ○     Object size            Δ     %   ○    On-disk size            Δ     %   ○")
		fmt.Println("------------------------------------------------------------------------------------------------------------------------")
		previewLinesCount = 2 // header + divider
	}

	// Collect growth statistics year by year
	for year := firstCommitTime.Year(); year <= currentYear; year++ {
		progress.StartSectionProgress(year)
		if cumulativeStatistics, err := git.GetGrowthStats(year, previous, debug); err == nil {
			totalStatistics = cumulativeStatistics
			previous = cumulativeStatistics
			yearlyStatistics[year] = cumulativeStatistics
		}
		progress.StopSectionProgress()

		// When running in a terminal, print each completed year's row immediately (preview with ... for percentages)
		if progress.ShowProgress {
			if cumulative, ok := yearlyStatistics[year]; ok {
				// Add row separator before current year
				if year == currentYear {
					fmt.Println("------------------------------------------------------------------------------------------------------------------------")
					previewLinesCount++
				}
				sections.PrintGrowthHistoryRowPreview(cumulative, interimPreviousCumulative, currentYear)
				previewLinesCount++
				// Add row separator after current year
				if year == currentYear {
					fmt.Println("------------------------------------------------------------------------------------------------------------------------")
					previewLinesCount++
				}

				interimPreviousCumulative = cumulative
			}
		}
	}

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

	// Save repository information with final totals (including authors)
	repositoryInformation := models.RepositoryInformation{
		Remote:           remote,
		LastCommit:       lastCommit,
		FirstCommit:      firstCommit,
		Age:              ageString,
		FirstDate:        firstCommitTime,
		TotalCommits:     totalStatistics.Commits,
		TotalAuthors:     totalAuthors,
		TotalTrees:       totalStatistics.Trees,
		TotalBlobs:       totalStatistics.Blobs,
		CompressedSize:   totalStatistics.Compressed,
		UncompressedSize: totalStatistics.Uncompressed,
	}

	// Recalculate final percentages using definitive totals for estimated growth computation
	var previousCumulative models.GrowthStatistics
	var previousDelta models.GrowthStatistics

	for year := repositoryInformation.FirstDate.Year(); year <= currentYear; year++ {
		if cumulative, ok := yearlyStatistics[year]; ok {
			cumulative.AuthorsDelta = cumulative.Authors - previousCumulative.Authors
			cumulative.CommitsDelta = cumulative.Commits - previousCumulative.Commits
			cumulative.TreesDelta = cumulative.Trees - previousCumulative.Trees
			cumulative.BlobsDelta = cumulative.Blobs - previousCumulative.Blobs
			cumulative.CompressedDelta = cumulative.Compressed - previousCumulative.Compressed
			cumulative.UncompressedDelta = cumulative.Uncompressed - previousCumulative.Uncompressed

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
			if repositoryInformation.UncompressedSize > 0 {
				cumulative.UncompressedPercent = float64(cumulative.UncompressedDelta) / float64(repositoryInformation.UncompressedSize) * 100
			}

			if previousDelta.Year != 0 {
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
				if previousDelta.UncompressedDelta > 0 {
					cumulative.UncompressedDeltaPercent = float64(cumulative.UncompressedDelta-previousDelta.UncompressedDelta) / float64(previousDelta.UncompressedDelta) * 100
				}
			}

			yearlyStatistics[year] = cumulative
			previousCumulative = cumulative
			previousDelta = cumulative
		}
	}

	// Clear the preview table and display complete growth table with correct percentages
	progress.ClearLines(previewLinesCount)
	sections.DisplayUnifiedGrowth(yearlyStatistics, repositoryInformation, firstCommitTime, recentFetch, lastModified)

	// 1. Largest file extensions
	sections.PrintTopFileExtensions(previous.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)

	// 2. Largest file extensions on-disk size growth
	sections.PrintFileExtensionGrowth(yearlyStatistics)

	// Prepare largest files data once for sections 3 & 4
	largestFiles := totalStatistics.LargestFiles
	sort.Slice(largestFiles, func(i, j int) bool {
		if largestFiles[i].CompressedSize != largestFiles[j].CompressedSize {
			return largestFiles[i].CompressedSize > largestFiles[j].CompressedSize
		}
		return largestFiles[i].Path < largestFiles[j].Path
	})

	var totalFilesCompressedSize int64
	for _, file := range largestFiles {
		totalFilesCompressedSize += file.CompressedSize
	}

	if len(largestFiles) > 10 {
		largestFiles = largestFiles[:10]
	}

	// 3. Largest directories
	sections.PrintLargestDirectories(totalStatistics.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)

	// 4. Largest files
	sections.PrintLargestFiles(largestFiles, totalFilesCompressedSize, repositoryInformation.TotalBlobs, len(previous.LargestFiles))

	// 5. Rate of changes analysis
	fmt.Println("\nRATE OF CHANGES ########################################################################################################")
	stopRateSpinner := progress.StartSimpleSpinner()
	ratesByYear, branchName, rateErr := git.GetRateOfChanges()
	stopRateSpinner()
	if rateErr == nil && len(ratesByYear) > 0 {
		sections.DisplayRateOfChanges(ratesByYear, branchName)
	}

	// 6 & 7. Authors with most commits, then Committers with most commits
	stopContributorsSpinner := progress.StartSimpleSpinner()
	topAuthorsByYear, totalAuthorsByYear, totalCommitsByYear, topCommittersByYear, totalCommittersByYear, allTimeAuthors, allTimeCommitters, contributorsErr := git.GetTopCommitAuthors(3)
	stopContributorsSpinner()
	if contributorsErr == nil && len(topAuthorsByYear) > 0 {
		sections.DisplayContributorsWithMostCommits(topAuthorsByYear, totalAuthorsByYear, totalCommitsByYear, topCommittersByYear, totalCommittersByYear, allTimeAuthors, allTimeCommitters)
	}

	// Get memory statistics for final output
	var memoryStatistics runtime.MemStats
	runtime.ReadMemStats(&memoryStatistics)

	fmt.Printf("\nFinished in %s with a memory footprint of %s.\n",
		utils.FormatDuration(time.Since(startTime)),
		strings.TrimSpace(utils.FormatSize(int64(memoryStatistics.Sys))))
}
