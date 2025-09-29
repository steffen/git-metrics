package models

import "time"

// RepositoryInformation holds information about a git repository
type RepositoryInformation struct {
	Remote           string
	Age              string
	FirstCommit      string
	LastCommit       string
	FirstDate        time.Time
	TotalCommits     int
	TotalAuthors     int
	TotalTrees       int
	TotalBlobs       int
	CompressedSize   int64
	UncompressedSize int64
}

// GrowthStatistics holds statistics about repository growth
type GrowthStatistics struct {
	Year         int
	Authors      int
	Commits      int
	Trees        int
	Blobs        int
	Compressed   int64
	Uncompressed int64
	RunTime      time.Duration
	LargestFiles []FileInformation

	// Delta values (year-over-year changes)
	AuthorsDelta      int
	CommitsDelta      int
	TreesDelta        int
	BlobsDelta        int
	CompressedDelta   int64
	UncompressedDelta int64

	// Percentage of total
	AuthorsPercent      float64
	CommitsPercent      float64
	TreesPercent        float64
	BlobsPercent        float64
	CompressedPercent   float64
	UncompressedPercent float64

	// Delta percentage changes (Δ%)
	AuthorsDeltaPercent      float64
	CommitsDeltaPercent      float64
	TreesDeltaPercent        float64
	BlobsDeltaPercent        float64
	CompressedDeltaPercent   float64
	UncompressedDeltaPercent float64
}

// FileInformation holds information about a file in the repository
type FileInformation struct {
	Path             string
	Blobs            int
	CompressedSize   int64
	UncompressedSize int64
	LastChange       time.Time
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
	DailyPeakP95         int // 95th percentile of daily commits
	DailyPeakP99         int // 99th percentile of daily commits
	DailyPeakP100        int // Maximum daily commits
	HourlyPeakP95        int // 95th percentile of hourly commits
	HourlyPeakP99        int // 99th percentile of hourly commits
	HourlyPeakP100       int // Maximum hourly commits
	MinutelyPeakP95      int // 95th percentile of commits per minute
	MinutelyPeakP99      int // 99th percentile of commits per minute
	MinutelyPeakP100     int // Maximum commits per minute
	PercentageOfTotal    float64
	MergeCommits         int     // Commits with >1 parent
	DirectCommits        int     // Regular commits
	MergeRatio           float64 // Percentage of commits that are merges
	BusiestDay           string  // Date with most commits
	BusiestDayCommits    int     // Number of commits on busiest day
	WorkdayCommits       int     // Commits during weekdays
	WeekendCommits       int     // Commits during weekends
	WorkdayWeekendRatio  float64 // Ratio of workday to weekend commits
	YearEndCommitHash    string  // Commit hash representing the final state of the year (last commit in that year)
}

// CheckoutGrowthStatistics holds checkout growth statistics for a specific year
type CheckoutGrowthStatistics struct {
	Year              int
	NumberDirectories int
	MaxPathDepth      int
	MaxPathLength     int
	NumberFiles       int
	TotalSizeFiles    int64
}
