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
