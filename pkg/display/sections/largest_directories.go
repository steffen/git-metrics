package sections

import (
	"fmt"
	"git-metrics/pkg/git"
	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	// MaxDirectoryDepth is the maximum depth of directories to process
	MaxDirectoryDepth = 10

	// CompressedSizePercentageThreshold is the minimum percentage of total compressed size
	// required for a directory or file to be considered significant (1%)
	CompressedSizePercentageThreshold = 0.01

	// PathColumnWidth is the fixed width for the path column in the output table
	PathColumnWidth = 80

	// MaxTreeLevels is the maximum number of tree levels to process (MaxDirectoryDepth + 1 for root)
	MaxTreeLevels = MaxDirectoryDepth + 1

	// TableRowFormat is the format string for printing table rows
	TableRowFormat = "%11s%6.1f %%   %11s%6.1f %%   %s\n"
)

// Footnote contains the formatted display path and footnote information
type Footnote struct {
	DisplayPath string // The path string to display, possibly truncated with a footnote marker
	Index       int    // Zero if no footnote needed, otherwise the footnote index
	FullPath    string // The full original path (for footnote)
}

// CreatePathFootnote formats a file path for display, truncating it if necessary
// and adding a footnote marker if the path is truncated
// maxDisplayLength is the maximum length of the displayed path
// currentFootnoteCount is the current number of footnotes
// Returns a Footnote containing the formatted path and footnote information
func CreatePathFootnote(path string, maxDisplayLength int, currentFootnoteCount int) Footnote {
	result := Footnote{
		DisplayPath: "",
		Index:       0,
		FullPath:    path,
	}

	// First check if truncation is needed
	truncatedPath := utils.TruncatePath(path, maxDisplayLength)
	if truncatedPath == path {
		// No truncation needed
		result.DisplayPath = path
		return result
	}

	// Truncation needed, add footnote
	footnoteIndex := currentFootnoteCount + 1
	marker := fmt.Sprintf(" [%d]", footnoteIndex)

	// Calculate the maximum truncated length to accommodate the marker
	maxTruncatedLength := maxDisplayLength - len(marker)
	if maxTruncatedLength < 0 {
		maxTruncatedLength = 0
	}

	// Truncate the path to make room for the marker
	truncatedForMarker := utils.TruncatePath(path, maxTruncatedLength)
	displayPath := truncatedForMarker + marker

	// Ensure displayPath is not longer than maxDisplayLength
	// (trim from truncatedForMarker if needed)
	if len(displayPath) > maxDisplayLength {
		// Remove excess from truncatedForMarker part
		excess := len(displayPath) - maxDisplayLength
		if excess < len(truncatedForMarker) {
			truncatedForMarker = truncatedForMarker[:len(truncatedForMarker)-excess]
		} else {
			truncatedForMarker = ""
		}
		displayPath = truncatedForMarker + marker
	}

	result.DisplayPath = displayPath
	result.Index = footnoteIndex

	return result
}

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
		treePrefix            string // Tree formatting prefix
	}

	// Calculate the total compressed size of all blobs
	var totalBlobsCompressedSize int64
	for _, file := range files {
		totalBlobsCompressedSize += file.CompressedSize
	}

	// Calculate 1% threshold
	thresholdSize := float64(totalBlobsCompressedSize) * CompressedSizePercentageThreshold

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

		// Create entries for all directories in the path (up to MaxDirectoryDepth levels)
		for level := 0; level < len(pathParts)-1 && level < MaxDirectoryDepth; level++ {
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
			// Calculate the level based on path depth (limited to MaxDirectoryDepth)
			level := strings.Count(file.Path, "/") + 1
			if level > MaxDirectoryDepth {
				level = MaxDirectoryDepth
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

	// Build final sorted list following directory structure with proper tree formatting
	var sortedEntries []*entry

	// Add root entry first
	rootEntry := &entry{
		Path:                  ".",
		FullPath:              ".",
		Blobs:                 totalBlobs,
		CompressedSize:        totalBlobsCompressedSize,
		Level:                 0,
		IsFile:                false,
		ExistsInDefaultBranch: true,
		treePrefix:            "",
	}
	sortedEntries = append(sortedEntries, rootEntry)

	// Helper function to create tree prefixes
	createTreePrefix := func(level int, isLast []bool) string {
		if level <= 1 {
			return ""
		}

		var prefix strings.Builder
		for i := 1; i < level-1; i++ {
			if i < len(isLast) && isLast[i] {
				prefix.WriteString("   ") // Three spaces for completed branches
			} else {
				prefix.WriteString("│  ") // Pipe and two spaces for continuing branches
			}
		}

		if level > 1 {
			if level-1 < len(isLast) && isLast[level-1] {
				prefix.WriteString("└─ ") // Last item at this level
			} else {
				prefix.WriteString("├─ ") // Not last item at this level
			}
		}

		return prefix.String()
	}

	// Process entries level by level to maintain proper tree structure
	processedPaths := make(map[string]bool)

	var buildTree func(level int, parentPath string, isLastAtLevel []bool)
	buildTree = func(level int, parentPath string, isLastAtLevel []bool) {
		if level > MaxTreeLevels { // Now max MaxTreeLevels levels (0-MaxDirectoryDepth, with 0 being root)
			return
		}

		key := fmt.Sprintf("%d:%s", level, parentPath)
		group, exists := levelGroups[key]
		if !exists {
			return
		}

		// Separate directories and files
		var directories []*entry
		var files []*entry

		for _, entry := range group {
			if processedPaths[entry.FullPath] {
				continue // Skip already processed entries
			}

			if entry.IsFile {
				files = append(files, entry)
			} else {
				directories = append(directories, entry)
			}
		}

		// Combine directories first, then files
		allEntries := append(directories, files...)

		for i, entry := range allEntries {
			if processedPaths[entry.FullPath] {
				continue
			}

			isLast := i == len(allEntries)-1

			// Ensure the slice is large enough
			newIsLastAtLevel := make([]bool, level+1)
			if len(isLastAtLevel) > 0 {
				copy(newIsLastAtLevel, isLastAtLevel)
			}
			newIsLastAtLevel[level] = isLast

			// Create tree prefix for this entry (adjust level for display)
			entry.treePrefix = createTreePrefix(level+1, newIsLastAtLevel)

			sortedEntries = append(sortedEntries, entry)
			processedPaths[entry.FullPath] = true

			// If this is a directory, process its children
			if !entry.IsFile {
				buildTree(level+1, entry.FullPath, newIsLastAtLevel)
			}
		}
	}

	// Start processing from level 1 (root level)
	buildTree(1, "", []bool{})

	// Print header
	fmt.Println("\nLARGEST DIRECTORIES ############################################################################")
	fmt.Println()
	fmt.Println("Showing directories and files that contribute more than 1% of total on-disk size.")

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
	fmt.Println("       Blobs           On-disk size           Path")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

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

		// Create indentation based on tree structure
		prefix := entry.treePrefix

		// Add asterisk if not in default branch
		displayName := entry.Path // Use just the name at this level, not full path

		// Add trailing slash for directories
		if !entry.IsFile {
			displayName += "/"
		}

		if hasDefaultBranch && !entry.ExistsInDefaultBranch {
			displayName += "*"
			showFootnote = true
		} // Calculate available width for path display (fixed width for alignment)
		pathColumnWidth := PathColumnWidth

		// Use CreatePathFootnote for consistent truncation and footnote logic
		result := CreatePathFootnote(displayName, pathColumnWidth-utf8.RuneCountInString(prefix), len(footnotes))
		finalDisplayName := result.DisplayPath
		if result.Index > 0 {
			footnotes = append(footnotes, Footnote{
				Index:    result.Index,
				FullPath: result.FullPath,
			})
		}

		// Create the full path display with prefix
		fullPathDisplay := prefix + finalDisplayName

		// Print entry with fixed column widths
		fmt.Printf(TableRowFormat,
			utils.FormatNumber(entry.Blobs),
			percentBlobs,
			utils.FormatSize(entry.CompressedSize),
			percentSize,
			fullPathDisplay,
		)

		totalSelectedBlobs += entry.Blobs
		totalSelectedSize += entry.CompressedSize
	}

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
