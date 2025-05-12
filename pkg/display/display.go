package display

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"git-metrics/pkg/git"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
)

// PrintLargestDirectories prints the largest root and subdirectories by size and object count
func PrintLargestDirectories(files []models.FileInformation, totalBlobs int, totalCompressedSize int64) {
	type dirStats struct {
		Path                  string
		Blobs                 int
		CompressedSize        int64
		Children              map[string]*dirStats // for subdirectories
		IsRoot                bool
		ExistsInDefaultBranch bool // Whether the directory exists in default branch
	}

	// Helper to get root dir or (root files)
	getRoot := func(path string) string {
		if !strings.Contains(path, "/") {
			return "(root files)"
		}
		return strings.SplitN(path, "/", 2)[0]
	}

	// Get files from default branch for comparison
	defaultBranch, defaultBranchError := git.GetDefaultBranch()
	hasDefaultBranch := defaultBranchError == nil
	defaultBranchFiles, defaultBranchFilesError := git.GetBranchFiles(defaultBranch)

	// Check if directory exists in default branch
	directoryExistsInDefaultBranch := func(dirPath string) bool {
		// If we couldn't get default branch files, assume everything exists
		if defaultBranchFilesError != nil || defaultBranchFiles == nil {
			return true
		}

		// Special case for root files
		if dirPath == "(root files)" {
			// Check if there are any root level files in default branch
			for file := range defaultBranchFiles {
				if !strings.Contains(file, "/") {
					return true
				}
			}
			return false
		}

		// For regular directories, check if any file starts with this directory
		prefix := dirPath + "/"
		for file := range defaultBranchFiles {
			if strings.HasPrefix(file, prefix) || file == dirPath {
				return true
			}
		}
		return false
	}

	// Aggregate stats for root directories and root files
	rootStats := make(map[string]*dirStats)
	// For each file, update root and immediate children (subdirectories or files)
	for _, file := range files {
		root := getRoot(file.Path)
		if _, ok := rootStats[root]; !ok {
			rootStats[root] = &dirStats{
				Path:                  root,
				Children:              make(map[string]*dirStats),
				IsRoot:                true,
				ExistsInDefaultBranch: directoryExistsInDefaultBranch(root),
			}
		}
		stat := rootStats[root]
		stat.Blobs += file.Blobs
		stat.CompressedSize += file.CompressedSize

		// For immediate children: subdirectories or files directly under root
		if root == "(root files)" {
			// Files at root: each file is a child
			if _, ok := stat.Children[file.Path]; !ok {
				existsInDefaultBranch := defaultBranchFilesError == nil && defaultBranchFiles[file.Path]
				stat.Children[file.Path] = &dirStats{
					Path:                  file.Path,
					IsRoot:                false,
					ExistsInDefaultBranch: existsInDefaultBranch,
				}
			}
			child := stat.Children[file.Path]
			child.Blobs += file.Blobs
			child.CompressedSize += file.CompressedSize
		} else {
			// For files under a root directory
			parts := strings.SplitN(file.Path, "/", 3)
			if len(parts) == 2 {
				// File directly under root dir
				name := parts[1]
				fullPath := root + "/" + name
				if _, ok := stat.Children[name]; !ok {
					existsInDefaultBranch := defaultBranchFilesError == nil && defaultBranchFiles[fullPath]
					stat.Children[name] = &dirStats{
						Path:                  name,
						IsRoot:                false,
						ExistsInDefaultBranch: existsInDefaultBranch,
					}
				}
				child := stat.Children[name]
				child.Blobs += file.Blobs
				child.CompressedSize += file.CompressedSize
			} else if len(parts) > 2 {
				// File in a subdirectory: immediate child is the subdir
				sub := parts[1]
				subdirPath := root + "/" + sub
				if _, ok := stat.Children[sub]; !ok {
					// Check if this subdirectory exists in default branch
					existsInDefaultBranch := directoryExistsInDefaultBranch(subdirPath)
					stat.Children[sub] = &dirStats{
						Path:                  sub,
						IsRoot:                false,
						ExistsInDefaultBranch: existsInDefaultBranch,
					}
				}
				child := stat.Children[sub]
				child.Blobs += file.Blobs
				child.CompressedSize += file.CompressedSize
			}
		}
	}

	// Convert to slice and sort by size
	var roots []*dirStats
	for _, stat := range rootStats {
		roots = append(roots, stat)
	}
	sort.Slice(roots, func(i, j int) bool {
		if roots[i].CompressedSize != roots[j].CompressedSize {
			return roots[i].CompressedSize > roots[j].CompressedSize
		}
		return roots[i].Path < roots[j].Path
	})
	if len(roots) > 10 {
		roots = roots[:10]
	}

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

	// Calculate the total compressed size of all blobs (sum of all file.CompressedSize)
	var totalBlobsCompressedSize int64
	for _, file := range files {
		totalBlobsCompressedSize += file.CompressedSize
	}

	// Track totals for displayed roots (top 10)
	var totalSelectedBlobs int
	var totalSelectedSize int64

	// Track if any items aren't in default branch (to show footnote)
	showFootnote := false

	// Print root entries
	for i, stat := range roots {
		// Print separator after each root except the last
		if i > 0 {
			fmt.Println("------------------------------------------------------------------------------------------------")
		}

		// Calculate percentages
		percentBlobs := 0.0
		percentSize := 0.0
		if totalBlobs > 0 {
			percentBlobs = float64(stat.Blobs) / float64(totalBlobs) * 100
		}
		if totalBlobsCompressedSize > 0 {
			percentSize = float64(stat.CompressedSize) / float64(totalBlobsCompressedSize) * 100
		}

		// Add asterisk if not in default branch
		displayPath := stat.Path
		if hasDefaultBranch && !stat.ExistsInDefaultBranch {
			displayPath += "*"
			showFootnote = true
		}

		// Print root
		fmt.Printf("%-51s %13s%6.1f %%  %13s%6.1f %%\n",
			displayPath,
			utils.FormatNumber(stat.Blobs),
			percentBlobs,
			utils.FormatSize(stat.CompressedSize),
			percentSize,
		)

		totalSelectedBlobs += stat.Blobs
		totalSelectedSize += stat.CompressedSize

		// Print up to 10 largest immediate children (subdirs or files) for this root
		var children []*dirStats
		for _, child := range stat.Children {
			children = append(children, child)
		}
		sort.Slice(children, func(i, j int) bool {
			if children[i].CompressedSize != children[j].CompressedSize {
				return children[i].CompressedSize > children[j].CompressedSize
			}
			return children[i].Path < children[j].Path
		})
		if len(children) > 10 {
			children = children[:10]
		}
		for idx, child := range children {
			percentBlobs := 0.0
			percentSize := 0.0
			if totalBlobs > 0 {
				percentBlobs = float64(child.Blobs) / float64(totalBlobs) * 100
			}
			if totalBlobsCompressedSize > 0 {
				percentSize = float64(child.CompressedSize) / float64(totalBlobsCompressedSize) * 100
			}
			prefix := "├─"
			if idx == len(children)-1 {
				prefix = "└─"
			}

			// Add asterisk if not in default branch
			displayPath := child.Path
			if hasDefaultBranch && !child.ExistsInDefaultBranch {
				displayPath += "*"
				showFootnote = true
			}

			fmt.Printf("%s %-48s %13s%6.1f %%  %13s%6.1f %%\n",
				prefix,
				displayPath,
				utils.FormatNumber(child.Blobs),
				percentBlobs,
				utils.FormatSize(child.CompressedSize),
				percentSize,
			)
		}
	}

	// Print separator and summary rows for roots (whole table)
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("%-51s %13s%6.1f %%  %13s%6.1f %%\n",
		fmt.Sprintf("├─ Top %s", utils.FormatNumber(len(roots))),
		utils.FormatNumber(totalSelectedBlobs),
		float64(totalSelectedBlobs)/float64(totalBlobs)*100,
		utils.FormatSize(totalSelectedSize),
		float64(totalSelectedSize)/float64(totalBlobsCompressedSize)*100)
	fmt.Printf("%-51s %13s%6.1f %%  %13s%6.1f %%\n",
		fmt.Sprintf("└─ Out of %s", utils.FormatNumber(len(rootStats))),
		utils.FormatNumber(totalBlobs),
		100.0,
		utils.FormatSize(totalBlobsCompressedSize),
		100.0)

	// Add footnote explaining the asterisk meaning
	if hasDefaultBranch && showFootnote {
		fmt.Println()
		fmt.Printf("* File or directory not present in latest commit of %s branch (moved, renamed or removed)\n", defaultBranch)
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
		fmt.Println("------------------------------------------------------------------------------------------------")
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
	var footnotes []struct {
		index int
		path  string
	}

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

		// Use FormatDisplayPath for consistent truncation and footnote logic
		result := FormatDisplayPath(file.Path, 43, len(footnotes))
		displayPath := result.DisplayPath
		if result.FootnoteIndex > 0 {
			footnotes = append(footnotes, struct {
				index int
				path  string
			}{
				index: result.FootnoteIndex,
				path:  result.FullPath,
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
			fmt.Printf("[%d] %s\n", footnote.index, footnote.path)
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

// PrintMachineInfo prints information about the system
func PrintMachineInformation() {
	fmt.Println()
	fmt.Println("RUN ############################################################################################")
	fmt.Println()
	fmt.Printf("Start time                 %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04 MST"))
	fmt.Printf("Machine                    %d CPU cores with %d GB memory (%s on %s)\n",
		runtime.NumCPU(),
		utils.GetMemoryInGigabytes(),
		utils.GetOperatingSystemInformation(),
		utils.GetChipInformation())
	fmt.Printf("Git version                %s\n", git.GetGitVersion())
}

// PrintTopCommitAuthors prints the top commit authors by number of commits
func PrintTopCommitAuthors(authors [][2]string) {
	fmt.Println("\nAUTHORS WITH MOST COMMITS #######################################################################")
	fmt.Println()
	fmt.Println("Author                                 Commits")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, entry := range authors {
		fmt.Printf("%-38s %8s\n", entry[0], entry[1])
	}
}
