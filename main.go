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

	// Flags
	repositoryPath := pflag.StringP("repository", "r", ".", "Path to git repository")
	showVersion := pflag.Bool("version", false, "Display version information and exit")
	pflag.BoolVar(&debug, "debug", false, "Enable debug output")
	noProgress := pflag.Bool("no-progress", false, "Disable progress indicators")
	showHelp := pflag.BoolP("help", "h", false, "Display this help message")
	pflag.Parse()

	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}
	if *showVersion {
		fmt.Printf("git-metrics version %s\n", utils.GetGitMetricsVersion())
		os.Exit(0)
	}

	// Progress visibility (disabled if redirected)
	progress.ShowProgress = !*noProgress && utils.IsTerminal(os.Stdout)

	if !requirements.CheckRequirements() {
		fmt.Println("\nRequirements not met. Please install listed dependencies above.")
		os.Exit(9)
	}

	// Validate repository
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

	// Git directory last modified
	lastModified := UnknownValue
	if info, err := os.Stat(gitDir); err == nil {
		lastModified = info.ModTime().Format("Mon, 02 Jan 2006 15:04 MST")
	}

	fmt.Printf("Git directory              %s\n", gitDir)

	// Remote
	remoteOutput, err := git.RunGitCommand(debug, "remote", "get-url", "origin")
	remote := ""
	if err == nil {
		trimmed := strings.TrimSpace(string(remoteOutput))
		if trimmed != "" {
			if progress.ShowProgress { fmt.Printf("Remote                     ... fetching\n") }
			remote = trimmed
			if progress.ShowProgress {
				fmt.Printf("\033[1A\033[2KRemote                     %s\n", remote)
			} else {
				fmt.Printf("Remote                     %s\n", remote)
			}
		}
	}

	// Fetch time
	recentFetch := git.GetLastFetchTime(gitDir)
	if recentFetch == "" {
		fmt.Printf("Last modified              %s\n", lastModified)
	} else {
		fmt.Printf("Most recent fetch          %s\n", recentFetch)
	}

	// Most recent commit
	if progress.ShowProgress { fmt.Printf("Most recent commit         ... fetching\n") }
	lastCommit := UnknownValue
	if out, err := git.RunGitCommand(debug, "rev-parse", "--short", "HEAD"); err == nil {
		hash := strings.TrimSpace(string(out))
		dateCmd := exec.Command("git", "show", "-s", "--format=%cD", hash)
		if dcOut, derr := dateCmd.Output(); derr == nil {
			if d, perr := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(string(dcOut))); perr == nil {
				lastCommit = fmt.Sprintf("%s (%s)", d.Format("Mon, 02 Jan 2006"), hash)
			}
		}
	}
	if progress.ShowProgress { fmt.Printf("\033[1A\033[2KMost recent commit         %s\n", lastCommit) } else { fmt.Printf("Most recent commit         %s\n", lastCommit) }

	// First commit & age
	if progress.ShowProgress { fmt.Printf("First commit               ... fetching\n") }
	firstCommit := UnknownValue
	ageString := UnknownValue
	var firstCommitTime time.Time
	if out, err := git.RunGitCommand(debug, "rev-list", "--max-parents=0", "HEAD", "--format=%cD"); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		type cinfo struct { hash string; date time.Time }
		var commits []cinfo
		for i := 0; i+1 < len(lines); i += 2 {
			hash := strings.TrimPrefix(lines[i], "commit ")
			if len(hash) >= 6 { hash = hash[:6] }
			if d, perr := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(lines[i+1])); perr == nil {
				commits = append(commits, cinfo{hash: hash, date: d})
			}
		}
		if len(commits) > 0 {
			sort.Slice(commits, func(i, j int) bool { return commits[i].date.Before(commits[j].date) })
			first := commits[0]
			firstCommitTime = first.date
			firstCommit = fmt.Sprintf("%s (%s)", first.date.Format("Mon, 02 Jan 2006"), first.hash)
			now := time.Now()
			years, months, days := utils.CalculateYearsMonthsDays(first.date, now)
			var parts []string
			if years > 0 { parts = append(parts, fmt.Sprintf("%d years", years)) }
			if months > 0 { parts = append(parts, fmt.Sprintf("%d months", months)) }
			if days > 0 { parts = append(parts, fmt.Sprintf("%d days", days)) }
			ageString = strings.Join(parts, " ")
		}
	}
	if progress.ShowProgress { fmt.Printf("\033[1A\033[2KFirst commit               %s\n", firstCommit) } else { fmt.Printf("First commit               %s\n", firstCommit) }
	if firstCommit == UnknownValue { fmt.Println("\n\nNo commits found in the repository."); os.Exit(2) }
	fmt.Printf("Age                        %s\n", ageString)

	// Historic & estimated growth header
	fmt.Println() 
	fmt.Println("HISTORIC & ESTIMATED GROWTH ############################################################################################")
	fmt.Println()
	fmt.Println("Year          Commits          Δ     %   ○     Object size            Δ     %   ○    On-disk size            Δ     %   ○")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	var previous models.GrowthStatistics
	var totalStatistics models.GrowthStatistics
	yearlyStatistics := make(map[int]models.GrowthStatistics)

	for year := firstCommitTime.Year(); year <= time.Now().Year(); year++ {
		progress.StartProgress(year, previous, startTime)
		if stats, err := git.GetGrowthStats(year, previous, debug); err == nil {
			totalStatistics = stats
			previous = stats
			yearlyStatistics[year] = stats
			progress.CurrentProgress.Statistics = stats
		}
	}
	progress.StopProgress()

	// Cumulative authors
	if cumulativeAuthorsByYear, totalAuthors, err := git.GetCumulativeUniqueAuthorsByYear(); err == nil {
		for year, stats := range yearlyStatistics {
			if authors, ok := cumulativeAuthorsByYear[year]; ok {
				stats.Authors = authors
				yearlyStatistics[year] = stats
			}
		}
		// Add repository info later after percentages computed
		_ = totalAuthors
	}

	// Repository info (authors total set later via cumulative data loop)
	repositoryInformation := models.RepositoryInformation{
		Remote:           remote,
		LastCommit:       lastCommit,
		FirstCommit:      firstCommit,
		Age:              ageString,
		FirstDate:        firstCommitTime,
		TotalCommits:     totalStatistics.Commits,
		TotalAuthors:     0, // will adjust below if we can derive
		TotalTrees:       totalStatistics.Trees,
		TotalBlobs:       totalStatistics.Blobs,
		CompressedSize:   totalStatistics.Compressed,
		UncompressedSize: totalStatistics.Uncompressed,
	}

	// Determine total unique authors from last year with data
	if len(yearlyStatistics) > 0 {
		maxYear := 0
		for y := range yearlyStatistics { if y > maxYear { maxYear = y } }
		repositoryInformation.TotalAuthors = yearlyStatistics[maxYear].Authors
	}

	currentYear := time.Now().Year()
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

			if repositoryInformation.TotalAuthors > 0 { cumulative.AuthorsPercent = float64(cumulative.AuthorsDelta) / float64(repositoryInformation.TotalAuthors) * 100 }
			if repositoryInformation.TotalCommits > 0 { cumulative.CommitsPercent = float64(cumulative.CommitsDelta) / float64(repositoryInformation.TotalCommits) * 100 }
			if repositoryInformation.TotalTrees > 0 { cumulative.TreesPercent = float64(cumulative.TreesDelta) / float64(repositoryInformation.TotalTrees) * 100 }
			if repositoryInformation.TotalBlobs > 0 { cumulative.BlobsPercent = float64(cumulative.BlobsDelta) / float64(repositoryInformation.TotalBlobs) * 100 }
			if repositoryInformation.CompressedSize > 0 { cumulative.CompressedPercent = float64(cumulative.CompressedDelta) / float64(repositoryInformation.CompressedSize) * 100 }
			if repositoryInformation.UncompressedSize > 0 { cumulative.UncompressedPercent = float64(cumulative.UncompressedDelta) / float64(repositoryInformation.UncompressedSize) * 100 }

			if previousDelta.Year != 0 {
				if previousDelta.AuthorsDelta > 0 { cumulative.AuthorsDeltaPercent = diffPercent(cumulative.AuthorsDelta, previousDelta.AuthorsDelta) }
				if previousDelta.CommitsDelta > 0 { cumulative.CommitsDeltaPercent = diffPercent(cumulative.CommitsDelta, previousDelta.CommitsDelta) }
				if previousDelta.TreesDelta > 0 { cumulative.TreesDeltaPercent = diffPercent(cumulative.TreesDelta, previousDelta.TreesDelta) }
				if previousDelta.BlobsDelta > 0 { cumulative.BlobsDeltaPercent = diffPercent(cumulative.BlobsDelta, previousDelta.BlobsDelta) }
				if previousDelta.CompressedDelta > 0 { cumulative.CompressedDeltaPercent = diffPercent64(cumulative.CompressedDelta, previousDelta.CompressedDelta) }
				if previousDelta.UncompressedDelta > 0 { cumulative.UncompressedDeltaPercent = diffPercent64(cumulative.UncompressedDelta, previousDelta.UncompressedDelta) }
			}
			yearlyStatistics[year] = cumulative
			previousCumulative = cumulative
			previousDelta = cumulative
		}
	}

	sections.DisplayUnifiedGrowth(yearlyStatistics, repositoryInformation, firstCommitTime, recentFetch, lastModified)

	sections.PrintTopFileExtensions(previous.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)
	sections.PrintFileExtensionGrowth(yearlyStatistics)

	largestFiles := totalStatistics.LargestFiles
	sort.Slice(largestFiles, func(i, j int) bool {
		if largestFiles[i].CompressedSize != largestFiles[j].CompressedSize {
			return largestFiles[i].CompressedSize > largestFiles[j].CompressedSize
		}
		return largestFiles[i].Path < largestFiles[j].Path
	})
	var totalFilesCompressedSize int64
	for _, f := range largestFiles { totalFilesCompressedSize += f.CompressedSize }
	if len(largestFiles) > 10 { largestFiles = largestFiles[:10] }
	sections.PrintLargestDirectories(totalStatistics.LargestFiles, repositoryInformation.TotalBlobs, repositoryInformation.CompressedSize)
	sections.PrintLargestFiles(largestFiles, totalFilesCompressedSize, repositoryInformation.TotalBlobs, len(previous.LargestFiles))

	// Rate of changes (provides commit hashes for checkout growth)
	ratesByYear, branchName, ratesErr := git.GetRateOfChanges()
	if ratesErr == nil && len(ratesByYear) > 0 {
		sections.DisplayRateOfChanges(ratesByYear, branchName)
	}

	// Contributors (authors & committers)
	if topAuthorsByYear, totalAuthorsByYear, totalCommitsByYear, topCommittersByYear, totalCommittersByYear, allTimeAuthors, allTimeCommitters, err := git.GetTopCommitAuthors(3); err == nil && len(topAuthorsByYear) > 0 {
		sections.DisplayContributorsWithMostCommits(topAuthorsByYear, totalAuthorsByYear, totalCommitsByYear, topCommittersByYear, totalCommittersByYear, allTimeAuthors, allTimeCommitters)
	}

	// Checkout growth (reuse ratesByYear commit hashes)
	checkoutStatistics := make(map[int]models.CheckoutGrowthStatistics)
	if len(ratesByYear) > 0 { // only if we have rate data
		for year := firstCommitTime.Year(); year <= time.Now().Year(); year++ {
			commitHash := ""
			if rs, ok := ratesByYear[year]; ok { commitHash = rs.YearEndCommitHash }
			if stats, err := git.GetCheckoutGrowthStats(year, commitHash, debug); err == nil {
				checkoutStatistics[year] = stats
			}
		}
	}
	sections.DisplayCheckoutGrowth(checkoutStatistics)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("\nFinished in %s with a memory footprint of %s.\n", utils.FormatDuration(time.Since(startTime)), strings.TrimSpace(utils.FormatSize(int64(mem.Sys))))
}

func diffPercent(newVal, oldVal int) float64 { return float64(newVal-oldVal) / float64(oldVal) * 100 }
func diffPercent64(newVal, oldVal int64) float64 { return float64(newVal-oldVal) / float64(oldVal) * 100 }
