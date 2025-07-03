package models

// Flag constants for command line arguments
var (
	FlagDebug = Flag{
		Name:        "debug",
		Description: "Enable debug output",
	}
	FlagNoProgress = Flag{
		Name:        "no-progress",
		Description: "Disable progress indicators",
	}
	FlagRepository = Flag{
		Name:        "repository",
		ShortName:   "r",
		Description: "Path to git repository",
	}
	FlagVersion = Flag{
		Name:        "version",
		Description: "Display version information and exit",
	}
)
