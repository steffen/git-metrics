package progress

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"git-metrics/pkg/models"
	"git-metrics/pkg/utils"

	"golang.org/x/term"
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
	// progressStateMutex protects concurrent access to progress state.
	progressStateMutex sync.RWMutex

	// CurrentProgress holds the current progress state
	CurrentProgress ProgressState

	// ProgressSpinner is the spinner used for progress indication
	ProgressSpinner = NewSpinner()

	// ShowProgress determines whether to display progress
	ShowProgress bool

	// spinnerQuitChannel is used to signal the spinner goroutine to stop
	spinnerQuitChannel chan struct{}

	// previousProgressLength stores the character length of the last progress
	// line that was written. Together with the current terminal width this lets
	// us calculate how many physical rows the previous write occupies after a
	// possible terminal resize, so we can move the cursor up and clear them all.
	previousProgressLength int
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

// getTerminalWidth returns the current terminal width.
// Falls back to 120 columns if the width cannot be determined.
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 120
	}
	return width
}

// clearPreviousProgressLines moves the cursor up and clears every physical
// row that the previous progress write may occupy at the current terminal
// width. This handles the case where the terminal was resized after the
// last write, causing the line to wrap to more rows than originally intended.

func clearPreviousProgressLines(progressLineLength int, terminalWidth int) {
	if progressLineLength == 0 || terminalWidth <= 0 {
		return
	}

	// Calculate how many physical rows the previous write occupies now.
	physicalRows := (progressLineLength + terminalWidth - 1) / terminalWidth

	// Move cursor up for each extra wrapped row.
	if physicalRows > 1 {
		fmt.Printf("\033[%dA\r", physicalRows-1)
	} else {
		fmt.Printf("\r")
	}

	// Clear each row from top to bottom so no stale text remains.
	for i := 0; i < physicalRows; i++ {
		if i > 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("\033[2K")
	}

	// Move back up to the first row where we will write the new progress line.
	if physicalRows > 1 {
		fmt.Printf("\033[%dA", physicalRows-1)
	}
	fmt.Printf("\r")
}

func getCurrentProgressSnapshot() ProgressState {
	progressStateMutex.RLock()
	defer progressStateMutex.RUnlock()
	return CurrentProgress
}

func isCurrentProgressActive() bool {
	progressStateMutex.RLock()
	defer progressStateMutex.RUnlock()
	return CurrentProgress.Active
}

// SetCurrentProgressStatistics updates the current and previous cumulative
// statistics used by the progress renderer.
func SetCurrentProgressStatistics(currentStatistics models.GrowthStatistics, previousStatistics models.GrowthStatistics) {
	progressStateMutex.Lock()
	defer progressStateMutex.Unlock()
	CurrentProgress.Statistics = currentStatistics
	CurrentProgress.PreviousStatistics = previousStatistics
}

// UpdateProgress updates the progress display
func UpdateProgress() {
	snapshot := getCurrentProgressSnapshot()
	if !snapshot.Active || !ShowProgress {
		return
	}

	statistics := snapshot.Statistics
	previous := snapshot.PreviousStatistics

	// Calculate deltas from previous year
	commitsDelta := statistics.Commits - previous.Commits
	uncompressedDelta := statistics.Uncompressed - previous.Uncompressed
	compressedDelta := statistics.Compressed - previous.Compressed

	commitsConcernLevel := utils.GetConcernLevel("commits", int64(statistics.Commits))
	objectSizeConcernLevel := utils.GetConcernLevel("object-size", statistics.Uncompressed)
	onDiskSizeConcernLevel := utils.GetConcernLevel("disk-size", statistics.Compressed)

	progressLine := fmt.Sprintf("%-6s %14s %10s %5s %3s │%14s %12s %5s %3s │%14s %12s %5s %3s",
		fmt.Sprintf("%d %s", snapshot.Year, ProgressSpinner.Next()),
		utils.FormatNumber(statistics.Commits),
		formatDelta(commitsDelta),
		"...",
		commitsConcernLevel,
		utils.FormatSize(statistics.Uncompressed),
		formatSizeDelta(uncompressedDelta),
		"...",
		objectSizeConcernLevel,
		utils.FormatSize(statistics.Compressed),
		formatSizeDelta(compressedDelta),
		"...",
		onDiskSizeConcernLevel)

	terminalWidth := getTerminalWidth()

	progressStateMutex.RLock()
	currentProgressLineLength := previousProgressLength
	progressStateMutex.RUnlock()

	// Clear any wrapped rows from the previous progress write.
	clearPreviousProgressLines(currentProgressLineLength, terminalWidth)

	// Truncate the new line to fit within the current terminal width.
	if len(progressLine) > terminalWidth {
		progressLine = progressLine[:terminalWidth]
	}

	fmt.Printf("%s\033[K", progressLine)

	// Remember how long this write was so the next update can clean it up.
	progressStateMutex.Lock()
	previousProgressLength = len(progressLine)
	progressStateMutex.Unlock()
}

// StartProgress starts progress tracking
func StartProgress(year int, statistics models.GrowthStatistics, previousStatistics models.GrowthStatistics, programStart time.Time) {
	// Stop any existing spinner goroutine before starting a new one
	StopProgress()

	// Always update the state
	progressStateMutex.Lock()
	CurrentProgress = ProgressState{
		Year:               year,
		Statistics:         statistics,
		PreviousStatistics: previousStatistics,
		Active:             true,
		ProgramStart:       programStart,
	}
	progressStateMutex.Unlock()

	// Only show visual progress if ShowProgress is true
	if ShowProgress {
		// Create a new quit channel
		quitChannel := make(chan struct{})
		progressStateMutex.Lock()
		spinnerQuitChannel = quitChannel
		progressStateMutex.Unlock()

		// Show initial progress immediately
		UpdateProgress()

		// Start spinner updates with 125ms interval
		go func() {
			ticker := time.NewTicker(125 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if isCurrentProgressActive() {
						UpdateProgress()
					}
				case <-quitChannel:
					return
				}
			}
		}()
	}
}

// StopProgress stops progress tracking
func StopProgress() {
	var quitChannel chan struct{}
	var currentProgressLineLength int

	// Always update the state
	progressStateMutex.Lock()
	CurrentProgress.Active = false
	quitChannel = spinnerQuitChannel
	spinnerQuitChannel = nil
	currentProgressLineLength = previousProgressLength
	previousProgressLength = 0
	progressStateMutex.Unlock()

	// Signal the spinner goroutine to stop if it's running
	if quitChannel != nil {
		close(quitChannel)
	}

	// Only clear the progress line if ShowProgress is true
	if ShowProgress {
		terminalWidth := getTerminalWidth()
		clearPreviousProgressLines(currentProgressLineLength, terminalWidth)
		fmt.Printf("\033[K")
	}
}
