package progress

import (
	"fmt"
	"strings"
	"time"

	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"
)

// Spinner represents a text-based spinner animation
type Spinner struct {
	frames  []string
	current int
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames:  []string{"|", "/", "-", "\\"},
		current: 0,
	}
}

// Next returns the next frame of the spinner
func (s *Spinner) Next() string {
	frame := s.frames[s.current]
	s.current = (s.current + 1) % len(s.frames)
	return frame
}

// ProgressState tracks the current progress state
type ProgressState struct {
	Year               int
	Statistics         models.GrowthStatistics
	PreviousStatistics models.GrowthStatistics
	Active             bool
	ProgramStart       time.Time
}

var (
	// CurrentProgress holds the current progress state
	CurrentProgress ProgressState

	// ProgressSpinner is the spinner used for progress indication
	ProgressSpinner = NewSpinner()

	// ShowProgress determines whether to display progress
	ShowProgress bool

	// spinnerQuitChannel is used to signal the spinner goroutine to stop
	spinnerQuitChannel chan struct{}
)

// formatDelta formats a delta value with a + prefix for display during progress.
func formatDelta(delta int) string {
	if delta == 0 {
		return "..."
	}
	return fmt.Sprintf("+%s", strings.TrimSpace(utils.FormatNumber(delta)))
}

// formatSizeDelta formats a size delta value with a + prefix for display during progress.
func formatSizeDelta(delta int64) string {
	if delta == 0 {
		return "..."
	}
	return fmt.Sprintf("+%s", strings.TrimSpace(utils.FormatSize(delta)))
}

// formatPercent formats a percentage value for display during progress.
func formatPercent(value float64) string {
	if value == 0 {
		return "..."
	}
	return fmt.Sprintf("%d %%", int(value))
}

// UpdateProgress updates the progress display
func UpdateProgress() {
	if !CurrentProgress.Active || !ShowProgress {
		return
	}

	statistics := CurrentProgress.Statistics
	previous := CurrentProgress.PreviousStatistics

	// Calculate deltas from previous year
	commitsDelta := statistics.Commits - previous.Commits
	uncompressedDelta := statistics.Uncompressed - previous.Uncompressed
	compressedDelta := statistics.Compressed - previous.Compressed

	// Calculate percentage of current total
	commitsPercent := 0.0
	if statistics.Commits > 0 {
		commitsPercent = float64(commitsDelta) / float64(statistics.Commits) * 100
	}
	uncompressedPercent := 0.0
	if statistics.Uncompressed > 0 {
		uncompressedPercent = float64(uncompressedDelta) / float64(statistics.Uncompressed) * 100
	}
	compressedPercent := 0.0
	if statistics.Compressed > 0 {
		compressedPercent = float64(compressedDelta) / float64(statistics.Compressed) * 100
	}

	fmt.Printf("\r%-6s %14s %10s %5s %3s │%14s %12s %5s %3s │%14s %12s %5s %3s",
		fmt.Sprintf("%d %s", CurrentProgress.Year, ProgressSpinner.Next()),
		utils.FormatNumber(statistics.Commits),
		formatDelta(commitsDelta),
		formatPercent(commitsPercent),
		".",
		utils.FormatSize(statistics.Uncompressed),
		formatSizeDelta(uncompressedDelta),
		formatPercent(uncompressedPercent),
		".",
		utils.FormatSize(statistics.Compressed),
		formatSizeDelta(compressedDelta),
		formatPercent(compressedPercent),
		".")
}

// StartProgress starts progress tracking
func StartProgress(year int, statistics models.GrowthStatistics, previousStatistics models.GrowthStatistics, programStart time.Time) {
	// Stop any existing spinner goroutine before starting a new one
	StopProgress()

	// Always update the state
	CurrentProgress = ProgressState{
		Year:               year,
		Statistics:         statistics,
		PreviousStatistics: previousStatistics,
		Active:             true,
		ProgramStart:       programStart,
	}

	// Only show visual progress if ShowProgress is true
	if ShowProgress {
		// Create a new quit channel
		spinnerQuitChannel = make(chan struct{})

		// Show initial progress immediately
		UpdateProgress()

		// Start spinner updates with 125ms interval
		go func() {
			ticker := time.NewTicker(125 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if CurrentProgress.Active {
						UpdateProgress()
					}
				case <-spinnerQuitChannel:
					return
				}
			}
		}()
	}
}

// StopProgress stops progress tracking
func StopProgress() {
	// Always update the state
	CurrentProgress.Active = false

	// Signal the spinner goroutine to stop if it's running
	if spinnerQuitChannel != nil {
		close(spinnerQuitChannel)
		spinnerQuitChannel = nil
	}

	// Only clear the progress line if ShowProgress is true
	if ShowProgress {
		fmt.Printf("\r\033[K") // Clear the progress line
	}
}
