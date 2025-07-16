package models

import "time"

// RepositoryInformation holds information about a git repository
type RepositoryInformation struct {
	Remote         string
	Age            string
	FirstCommit    string
	LastCommit     string
	FirstDate      time.Time
	TotalCommits   int
	TotalTrees     int
	TotalBlobs     int
	CompressedSize int64
}

// GrowthStatistics holds statistics about repository growth
type GrowthStatistics struct {
	Year         int
	Commits      int
	Trees        int
	Blobs        int
	Compressed   int64
	RunTime      time.Duration
	LargestFiles []FileInformation
}

// FileInformation holds information about a file in the repository
type FileInformation struct {
	Path           string
	Blobs          int
	CompressedSize int64
	LastChange     time.Time
}

// GitObject represents a git object with its details
type GitObject struct {
	ObjectType       string
	ObjectIdentifier string
	ObjectSize       int64
	Additional       string // typically the file path if available
}

// RateStatistics holds commit rate statistics for a specific year
type RateStatistics struct {
	Year                 int
	TotalCommits         int
	AverageCommitsPerDay float64
	DailyPeakP95         int     // 95th percentile of daily commits
	HourlyPeakP95        int     // 95th percentile of hourly commits (for peak days)
	MinutelyPeakP95      float64 // 95th percentile of commits per minute (for peak hours)
	PercentageOfTotal    float64
	MergeCommits         int     // Commits with >1 parent
	DirectCommits        int     // Regular commits
	MergeRatio           float64 // Percentage of commits that are merges
	BusiestDay           string  // Date with most commits
	BusiestDayCommits    int     // Number of commits on busiest day
	WorkdayCommits       int     // Commits during weekdays
	WeekendCommits       int     // Commits during weekends
	WorkdayWeekendRatio  float64 // Ratio of workday to weekend commits
}
