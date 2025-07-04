package display

import (
	"fmt"
	"git-metrics/pkg/utils"
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
