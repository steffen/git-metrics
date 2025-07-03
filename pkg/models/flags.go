package models

// Flag constants for command line arguments
var (
	FlagDebug = Flag{
		Name:        "debug",
		Description: "Enable debug output",
		Type:        "bool",
	}
	FlagNoProgress = Flag{
		Name:        "no-progress",
		Description: "Disable progress indicators",
		Type:        "bool",
	}
	FlagRepository = Flag{
		Name:        "repository",
		ShortName:   "r",
		Description: "Path to git repository",
		Type:        "string",
		Default:     ".",
	}
	FlagVersion = Flag{
		Name:        "version",
		Description: "Display version information and exit",
		Type:        "bool",
	}
)

// AllFlags returns a slice of all defined flags for iteration
func AllFlags() []Flag {
	return []Flag{
		FlagRepository,
		FlagNoProgress, 
		FlagVersion,
		FlagDebug,
	}
}
