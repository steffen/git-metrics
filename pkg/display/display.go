package display

import (
	"fmt"
	"git-metrics/pkg/git"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// PrintLargestDirectories prints directories and files that are >= 1% of total on-disk size, up to 10 levels deep
func PrintLargestDirectories(files []models.FileInformation, totalBlobs int, totalCompressedSize int64) {
	type entry struct {
		Path                  string
		FullPath              string // Full path from root
		Blobs                 int
		CompressedSize        int64
		Level                 int
		IsFile                bool
		ExistsInDefaultBranch bool
	}

	// Calculate the total compressed size of all blobs
	var totalBlobsCompressedSize int64
	for _, file := range files {
		totalBlobsCompressedSize += file.CompressedSize
	}

	// Calculate 1% threshold
	thresholdSize := float64(totalBlobsCompressedSize) * 0.01

	// Get files from default branch for comparison
	defaultBranch, defaultBranchError := git.GetDefaultBranch()
	hasDefaultBranch := defaultBranchError == nil
	defaultBranchFiles, defaultBranchFilesError := git.GetBranchFiles(defaultBranch)

	// Check if path exists in default branch
	pathExistsInDefaultBranch := func(path string) bool {
		if defaultBranchFilesError != nil || defaultBranchFiles == nil {
			return true
		}
		return defaultBranchFiles[path]
	}

	// Check if directory exists in default branch
	directoryExistsInDefaultBranch := func(dirPath string) bool {
		if defaultBranchFilesError != nil || defaultBranchFiles == nil {
			return true
		}

		// For directories, check if any file starts with this directory
		prefix := dirPath + "/"
		for file := range defaultBranchFiles {
			if strings.HasPrefix(file, prefix) || file == dirPath {
				return true
			}
		}
		return false
	}

	// Collect all directories and their stats
	directoryStats := make(map[string]*entry)
	
	// First pass: collect all directories
	for _, file := range files {
		pathParts := strings.Split(file.Path, "/")
		currentPath := ""

		// Create entries for all directories in the path (up to 10 levels)
		for level := 0; level < len(pathParts)-1 && level < 10; level++ {
			if currentPath == "" {
				currentPath = pathParts[level]
			} else {
				currentPath = currentPath + "/" + pathParts[level]
			}

			if _, exists := directoryStats[currentPath]; !exists {
				directoryStats[currentPath] = &entry{
					Path:                  pathParts[level],
					FullPath:              currentPath,
					Level:                 level + 1,
					IsFile:                false,
					ExistsInDefaultBranch: directoryExistsInDefaultBranch(currentPath),
				}
			}

			// Add file stats to this directory
			dir := directoryStats[currentPath]
			dir.Blobs += file.Blobs
			dir.CompressedSize += file.CompressedSize
		}
	}

	// Collect significant entries (directories and files >= 1%)
	var significantEntries []*entry

	// Add significant directories
	for _, dir := range directoryStats {
		if float64(dir.CompressedSize) >= thresholdSize {
			significantEntries = append(significantEntries, dir)
		}
	}

	// Add significant files
	for _, file := range files {
		if float64(file.CompressedSize) >= thresholdSize {
			// Calculate the level based on path depth (limited to 10)
			level := strings.Count(file.Path, "/") + 1
			if level > 10 {
				level = 10
			}

			fileEntry := &entry{
				Path:                  filepath.Base(file.Path),
				FullPath:              file.Path,
				Blobs:                 file.Blobs,
				CompressedSize:        file.CompressedSize,
				Level:                 level,
				IsFile:                true,
				ExistsInDefaultBranch: pathExistsInDefaultBranch(file.Path),
			}
			significantEntries = append(significantEntries, fileEntry)
		}
	}

	// Group entries by their parent directory and level for hierarchical sorting
	levelGroups := make(map[string][]*entry)
	
	for _, entry := range significantEntries {
		var parentPath string
		if entry.Level == 1 {
			parentPath = ""
		} else {
			parentParts := strings.Split(entry.FullPath, "/")
			if entry.IsFile {
				parentPath = strings.Join(parentParts[:len(parentParts)-1], "/")
			} else {
				parentPath = strings.Join(parentParts[:len(parentParts)-1], "/")
			}
		}

		key := fmt.Sprintf("%d:%s", entry.Level, parentPath)
		levelGroups[key] = append(levelGroups[key], entry)
	}

	// Sort entries within each group by size percentage (descending)
	for _, group := range levelGroups {
		sort.Slice(group, func(i, j int) bool {
			percentI := float64(group[i].CompressedSize) / float64(totalBlobsCompressedSize) * 100
			percentJ := float64(group[j].CompressedSize) / float64(totalBlobsCompressedSize) * 100
			if percentI != percentJ {
				return percentI > percentJ
			}
			return group[i].FullPath < group[j].FullPath
		})
	}

	// Build final sorted list following directory structure
	var sortedEntries []*entry
	var processLevel func(level int, parentPath string)
	processLevel = func(level int, parentPath string) {
		key := fmt.Sprintf("%d:%s", level, parentPath)
		if group, exists := levelGroups[key]; exists {
			// Add directories first, then files
			for _, entry := range group {
				if !entry.IsFile {
					sortedEntries = append(sortedEntries, entry)
					// Recursively process children of this directory
					processLevel(level+1, entry.FullPath)
				}
			}
			for _, entry := range group {
				if entry.IsFile {
					sortedEntries = append(sortedEntries, entry)
				}
			}
		}
	}

	// Start processing from level 1 (root level)
	processLevel(1, "")

	// Print header
	fmt.Println("\nLARGEST DIRECTORIES ############################################################################")

	var missingPathsError error = nil
	if defaultBranchError != nil {
		missingPathsError = defaultBranchError
	} else if defaultBranchFilesError != nil {
		missingPathsError = defaultBranchFilesError
	}

	if missingPathsError != nil {
		fmt.Println()
		fmt.Printf("Warning: Could not determine moved, renamed or removed files and directories: %s\n", missingPathsError)
	}

	fmt.Println()
	fmt.Println("Path                                                        Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")

	// Track totals for displayed entries
	var totalSelectedBlobs int
	var totalSelectedSize int64

	// Track if any items aren't in default branch (to show footnote)
	showFootnote := false

	// Track truncated paths for footnotes
	var footnotes []Footnote

	// Print significant entries
	for _, entry := range sortedEntries {
		// Calculate percentages
		percentBlobs := 0.0
		percentSize := 0.0
		if totalBlobs > 0 {
			percentBlobs = float64(entry.Blobs) / float64(totalBlobs) * 100
		}
		if totalBlobsCompressedSize > 0 {
			percentSize = float64(entry.CompressedSize) / float64(totalBlobsCompressedSize) * 100
		}

		// Create indentation based on level
		indent := strings.Repeat("  ", entry.Level-1)
		var prefix string
		if entry.Level == 1 {
			prefix = ""
		} else {
			prefix = indent + "└─ "
		}

		// Add asterisk if not in default branch
		displayName := entry.Path // Use just the name at this level, not full path
		if hasDefaultBranch && !entry.ExistsInDefaultBranch {
			displayName += "*"
			showFootnote = true
		}

		// Calculate available width for path display
		availableWidth := 51 - len(prefix)
		
		// Use CreatePathFootnote for consistent truncation and footnote logic
		result := CreatePathFootnote(displayName, availableWidth, len(footnotes))
		finalDisplayName := result.DisplayPath
		if result.Index > 0 {
			footnotes = append(footnotes, Footnote{
				Index:    result.Index,
				FullPath: result.FullPath,
			})
		}

		// Print entry
		fmt.Printf("%s%-*s %13s%6.1f %%  %13s%6.1f %%\n",
			prefix,
			availableWidth,
			finalDisplayName,
			utils.FormatNumber(entry.Blobs),
			percentBlobs,
			utils.FormatSize(entry.CompressedSize),
			percentSize,
		)

		totalSelectedBlobs += entry.Blobs
		totalSelectedSize += entry.CompressedSize
	}

	// Print separator and summary rows
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("%-51s %13s%6.1f %%  %13s%6.1f %%\n",
		fmt.Sprintf("├─ Shown (%s entries ≥ 1%%)", utils.FormatNumber(len(sortedEntries))),
		utils.FormatNumber(totalSelectedBlobs),
		float64(totalSelectedBlobs)/float64(totalBlobs)*100,
		utils.FormatSize(totalSelectedSize),
		float64(totalSelectedSize)/float64(totalBlobsCompressedSize)*100)
	fmt.Printf("%-51s %13s%6.1f %%  %13s%6.1f %%\n",
		"└─ Total repository",
		utils.FormatNumber(totalBlobs),
		100.0,
		utils.FormatSize(totalBlobsCompressedSize),
		100.0)

	// Add footnote explaining the asterisk meaning
	if hasDefaultBranch && showFootnote {
		fmt.Println()
		fmt.Printf("* File or directory not present in latest commit of %s branch (moved, renamed or removed)\n", defaultBranch)
	}

	// Print footnotes for truncated paths
	if len(footnotes) > 0 {
		fmt.Println()
		for _, footnote := range footnotes {
			fmt.Printf("[%d] %s\n", footnote.Index, footnote.FullPath)
		}
	}
}

// PrintGrowthTableHeader prints the header for the growth table
func PrintGrowthTableHeader() {
	fmt.Println()
	fmt.Println("HISTORIC & ESTIMATED GROWTH ####################################################################")
	fmt.Println()
	fmt.Println("Year        Commits                  Trees                  Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")
}

// PrintGrowthTableRow prints a row of the growth table
func PrintGrowthTableRow(statistics, previous models.GrowthStatistics, information models.RepositoryInformation, isEstimate bool, currentYear int) {
	commitsDifference := float64(statistics.Commits-previous.Commits) / float64(information.TotalCommits) * 100
	treesDifference := float64(statistics.Trees-previous.Trees) / float64(information.TotalTrees) * 100
	blobsDifference := float64(statistics.Blobs-previous.Blobs) / float64(information.TotalBlobs) * 100
	compressedDifference := float64(statistics.Compressed-previous.Compressed) / float64(information.CompressedSize) * 100

	yearDisplay := strconv.Itoa(statistics.Year)
	if isEstimate {
		yearDisplay += "*"
	} else if statistics.Year == currentYear {
		// Only print separator if there are previous years of data
		if previous.Year > 0 {
			fmt.Println("------------------------------------------------------------------------------------------------")
		}
		yearDisplay += "^"
	}

	fmt.Printf("%-5s %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%  %13s %+5.0f %%\n",
		yearDisplay,
		utils.FormatNumber(statistics.Commits), commitsDifference,
		utils.FormatNumber(statistics.Trees), treesDifference,
		utils.FormatNumber(statistics.Blobs), blobsDifference,
		utils.FormatSize(statistics.Compressed), compressedDifference)
}

// PrintLargestFiles prints information about the largest files
func PrintLargestFiles(files []models.FileInformation, totalFilesSize int64, totalBlobs int, totalFiles int) {
	fmt.Println("\nLARGEST FILES ##################################################################################")
	fmt.Println()
	fmt.Println("File path                              Last commit          Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")

	// Track totals for the selected files
	var totalSelectedBlobs int
	var totalSelectedSize int64

	// Track truncated paths for footnotes
	var footnotes []Footnote

	// Calculate total size of all files in repository
	for _, file := range files {
		// Get the last change date for the file
		lastChangeCommand := exec.Command("git", "log", "-1", "--format=%cD", "--", file.Path)
		lastChangeOutput, err := lastChangeCommand.Output()
		if err == nil {
			lastChange, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(string(lastChangeOutput)))
			file.LastChange = lastChange
		}

		percentageSize := float64(file.CompressedSize) / float64(totalFilesSize) * 100
		percentageBlobs := float64(file.Blobs) / float64(totalBlobs) * 100

		// Use CreatePathFootnote for consistent truncation and footnote logic
		result := CreatePathFootnote(file.Path, 43, len(footnotes))
		displayPath := result.DisplayPath
		if result.Index > 0 {
			footnotes = append(footnotes, Footnote{
				Index:    result.Index,
				FullPath: result.FullPath,
			})
		}

		fmt.Printf("%-43s   %s  %13s %5.1f %%  %13s %5.1f %%\n",
			displayPath,
			file.LastChange.Format("2006"),
			utils.FormatNumber(file.Blobs),
			percentageBlobs,
			utils.FormatSize(file.CompressedSize),
			percentageSize)

		totalSelectedBlobs += file.Blobs
		totalSelectedSize += file.CompressedSize
	}

	// Print separator and selected files totals row
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("%-43s   %s  %13s %5.1f %%  %13s %5.1f %%\n",
		fmt.Sprintf("├─ Top %s", utils.FormatNumber(len(files))),
		"    ",
		utils.FormatNumber(totalSelectedBlobs),
		float64(totalSelectedBlobs)/float64(totalBlobs)*100,
		utils.FormatSize(totalSelectedSize),
		float64(totalSelectedSize)/float64(totalFilesSize)*100)

	// Print grand totals row
	fmt.Printf("%-43s   %s  %13s %5.1f %%  %13s %5.1f %%\n",
		fmt.Sprintf("└─ Out of %s", utils.FormatNumber(totalFiles)),
		"    ",
		utils.FormatNumber(totalBlobs),
		100.0,
		utils.FormatSize(totalFilesSize),
		100.0)

	// Print footnotes for truncated paths
	if len(footnotes) > 0 {
		fmt.Println()
		for _, footnote := range footnotes {
			fmt.Printf("[%d] %s\n", footnote.Index, footnote.FullPath)
		}
	}
}

// PrintTopFileExtensions prints the top file extensions by size
func PrintTopFileExtensions(blobs []models.FileInformation, totalBlobs int, totalSize int64) {
	extensionStatistics := make(map[string]struct {
		size       int64
		filesCount int
		blobsCount int
	})
	for _, blob := range blobs {
		extension := filepath.Ext(blob.Path)
		if extension == "" {
			extension = "No Extension"
		}
		statistics := extensionStatistics[extension]
		statistics.size += blob.CompressedSize
		statistics.filesCount++
		statistics.blobsCount += blob.Blobs
		extensionStatistics[extension] = statistics
	}

	// Create a slice for sorting.
	type extensionStatistic struct {
		extension  string
		size       int64
		filesCount int
		blobsCount int
	}
	var statistics []extensionStatistic
	for extension, statistic := range extensionStatistics {
		statistics = append(statistics, extensionStatistic{
			extension:  extension,
			size:       statistic.size,
			filesCount: statistic.filesCount,
			blobsCount: statistic.blobsCount,
		})
	}
	sort.Slice(statistics, func(i, j int) bool {
		return statistics[i].size > statistics[j].size
	})

	// Calculate totals from all extensions first
	var totalExtFilesCount, totalExtBlobsCount int
	var totalExtSize int64
	for _, statistic := range extensionStatistics {
		totalExtFilesCount += statistic.filesCount
		totalExtBlobsCount += statistic.blobsCount
		totalExtSize += statistic.size
	}

	// Limit to top 10
	if len(statistics) > 10 {
		statistics = statistics[:10]
	}

	// Track totals for displayed extensions (top 10)
	var selectedFilesCount int
	var selectedBlobsCount int
	var selectedSize int64

	// Display results.
	fmt.Println("\nLARGEST FILE EXTENSIONS ########################################################################")
	fmt.Println()
	fmt.Println("Extension                            Files                  Blobs           On-disk size")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, statistic := range statistics {
		percentageFiles := float64(statistic.filesCount) / float64(totalExtFilesCount) * 100
		percentageBlobs := float64(statistic.blobsCount) / float64(totalBlobs) * 100
		percentageSize := float64(statistic.size) / float64(totalSize) * 100
		fmt.Printf("%-28s %13s %5.1f %%  %13s %5.1f %%  %13s %5.1f %%\n",
			statistic.extension, utils.FormatNumber(statistic.filesCount), percentageFiles, utils.FormatNumber(statistic.blobsCount), percentageBlobs, utils.FormatSize(statistic.size), percentageSize)

		selectedFilesCount += statistic.filesCount
		selectedBlobsCount += statistic.blobsCount
		selectedSize += statistic.size
	}

	// Print separator and top 10 totals row
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("%-28s %13s %5.1f %%  %13s %5.1f %%  %13s %5.1f %%\n",
		fmt.Sprintf("├─ Top %s", utils.FormatNumber(len(statistics))),
		utils.FormatNumber(selectedFilesCount),
		float64(selectedFilesCount)/float64(totalExtFilesCount)*100,
		utils.FormatNumber(selectedBlobsCount),
		float64(selectedBlobsCount)/float64(totalExtBlobsCount)*100,
		utils.FormatSize(selectedSize),
		float64(selectedSize)/float64(totalExtSize)*100)

	// Print grand totals row using full totals
	fmt.Printf("%-28s %13s %5.1f %%  %13s %5.1f %%  %13s %5.1f %%\n",
		fmt.Sprintf("└─ Out of %s", utils.FormatNumber(len(extensionStatistics))),
		utils.FormatNumber(totalExtFilesCount),
		100.0, // Always 100% for totals
		utils.FormatNumber(totalExtBlobsCount),
		100.0,
		utils.FormatSize(totalExtSize),
		100.0)
}
