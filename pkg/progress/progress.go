package progress

import (
	"fmt"
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

	// sectionSpinnerQuitChannel is used to signal the section spinner goroutine to stop
	sectionSpinnerQuitChannel chan struct{}
)

// stopSectionSpinnerIfActive stops the section spinner goroutine and clears
// its line so the per-year progress can take over without a race condition.
func stopSectionSpinnerIfActive() {
	if sectionSpinnerQuitChannel != nil {
		close(sectionSpinnerQuitChannel)
		sectionSpinnerQuitChannel = nil
		fmt.Printf("\r\033[K")
	}
}

// UpdateProgress updates the progress display
func UpdateProgress() {
	if !CurrentProgress.Active || !ShowProgress {
		return
	}
	fmt.Printf("\r%-6s %14s %10s %5s %3s │%14s %12s %5s %3s │%14s %12s %5s %3s",
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
}

// StartProgress starts progress tracking
func StartProgress(year int, statistics models.GrowthStatistics, programStart time.Time) {
	// Stop any existing spinner goroutine before starting a new one
	StopProgress()

	// Stop any active section spinner so it does not race with per-year progress
	stopSectionSpinnerIfActive()

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
		fmt.Printf("\r\033[K") // Clear the progress line
	}
}

// StartSectionSpinner starts a spinner with the given label on stdout.
// Returns a function that stops the spinner and clears the line.
// If StartProgress is called while a section spinner is active, the section
// spinner is automatically stopped to avoid a race condition.
func StartSectionSpinner(label string) func() {
	if !ShowProgress {
		return func() {}
	}

	sectionSpinner := NewSpinner()
	sectionSpinnerQuitChannel = make(chan struct{})

	fmt.Printf("\r%s %s", label, sectionSpinner.Next())

	go func() {
		ticker := time.NewTicker(125 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r%s %s", label, sectionSpinner.Next())
			case <-sectionSpinnerQuitChannel:
				return
			}
		}
	}()

	return func() {
		if sectionSpinnerQuitChannel != nil {
			close(sectionSpinnerQuitChannel)
			sectionSpinnerQuitChannel = nil
		}
		fmt.Printf("\r\033[K")
	}
}

// StartSimpleSpinner starts a spinner with just the animation character (no text).
// Returns a function that stops the spinner and clears the line.
// Used for sections where the spinner appears below the header.
func StartSimpleSpinner() func() {
	if !ShowProgress {
		return func() {}
	}

	simpleSpinner := NewSpinner()
	sectionSpinnerQuitChannel = make(chan struct{})

	fmt.Printf("%s", simpleSpinner.Next())

	go func() {
		ticker := time.NewTicker(125 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r%s", simpleSpinner.Next())
			case <-sectionSpinnerQuitChannel:
				return
			}
		}
	}()

	return func() {
		if sectionSpinnerQuitChannel != nil {
			close(sectionSpinnerQuitChannel)
			sectionSpinnerQuitChannel = nil
		}
		fmt.Printf("\r\033[K")
	}
}

// SectionSpinner is a dedicated spinner for the section-level progress indicator
var SectionSpinner = NewSpinner()

// SectionProgressYear tracks which year the section spinner is currently displaying
var SectionProgressYear int

// SectionSpinnerActive tracks whether the section spinner is currently running
var SectionSpinnerActive bool

// StartSectionProgress starts a single continuous spinner at the bottom of the section
// showing which year is currently being processed
func StartSectionProgress(year int) {
	StopSectionProgress()

	SectionProgressYear = year
	SectionSpinnerActive = true

	if ShowProgress {
		sectionSpinnerQuitChannel = make(chan struct{})

		// Show initial spinner line
		fmt.Printf("\r%d %s", year, SectionSpinner.Next())

		go func() {
			ticker := time.NewTicker(125 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if SectionSpinnerActive {
						fmt.Printf("\r%d %s", SectionProgressYear, SectionSpinner.Next())
					}
				case <-sectionSpinnerQuitChannel:
					return
				}
			}
		}()
	}
}

// StopSectionProgress stops the section-level spinner and clears its line
func StopSectionProgress() {
	SectionSpinnerActive = false

	if sectionSpinnerQuitChannel != nil {
		close(sectionSpinnerQuitChannel)
		sectionSpinnerQuitChannel = nil
	}

	if ShowProgress {
		fmt.Printf("\r\033[K") // Clear the spinner line
	}
}

// ClearLines moves the cursor up by the specified number of lines and clears each line.
// Used to replace preview output with final output.
func ClearLines(count int) {
	if !ShowProgress || count <= 0 {
		return
	}
	for i := 0; i < count; i++ {
		fmt.Print("\033[A") // Move cursor up one line
		fmt.Print("\033[K") // Clear the line
	}
}
