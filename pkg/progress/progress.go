package progress

import (
	"fmt"
	"os"
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
	Year         int
	Statistics   models.GrowthStatistics
	Active       bool
	ProgramStart time.Time
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

	// previousProgressLength stores the character length of the last progress
	// line that was written. Together with the current terminal width this lets
	// us calculate how many physical rows the previous write occupies after a
	// possible terminal resize, so we can move the cursor up and clear them all.
	previousProgressLength int
)

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
func clearPreviousProgressLines(terminalWidth int) {
	if previousProgressLength == 0 || terminalWidth <= 0 {
		return
	}

	// Calculate how many physical rows the previous write occupies now.
	physicalRows := (previousProgressLength + terminalWidth - 1) / terminalWidth

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

// UpdateProgress updates the progress display
func UpdateProgress() {
	if !CurrentProgress.Active || !ShowProgress {
		return
	}

	progressLine := fmt.Sprintf("%-6s %14s %10s %5s %3s │%14s %12s %5s %3s │%14s %12s %5s %3s",
		fmt.Sprintf("%d %s", CurrentProgress.Year, ProgressSpinner.Next()),
		utils.FormatNumber(CurrentProgress.Statistics.Authors),
		"...",
		"...",
		".",
		utils.FormatNumber(CurrentProgress.Statistics.Commits),
		"...",
		"...",
		".",
		utils.FormatSize(CurrentProgress.Statistics.Compressed),
		"...",
		"...",
		".")

	terminalWidth := getTerminalWidth()

	// Clear any wrapped rows from the previous progress write.
	clearPreviousProgressLines(terminalWidth)

	// Truncate the new line to fit within the current terminal width.
	if len(progressLine) > terminalWidth {
		progressLine = progressLine[:terminalWidth]
	}

	fmt.Printf("%s\033[K", progressLine)

	// Remember how long this write was so the next update can clean it up.
	previousProgressLength = len(progressLine)
}

// StartProgress starts progress tracking
func StartProgress(year int, statistics models.GrowthStatistics, programStart time.Time) {
	// Stop any existing spinner goroutine before starting a new one
	StopProgress()

	// Always update the state
	CurrentProgress = ProgressState{
		Year:         year,
		Statistics:   statistics,
		Active:       true,
		ProgramStart: programStart,
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
		terminalWidth := getTerminalWidth()
		clearPreviousProgressLines(terminalWidth)
		fmt.Printf("\033[K")
		previousProgressLength = 0
	}
}
