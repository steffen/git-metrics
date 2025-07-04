package utils

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// PrintFlagUsage prints a custom usage message for command-line flags.
func PrintFlagUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	// Group flags by usage text to combine flags with same description
	flagGroups := make(map[string][]*flag.Flag)
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		flagGroups[f.Usage] = append(flagGroups[f.Usage], f)
	})

	// Process each group of flags
	for _, flags := range flagGroups {
		// Sort flags to show single-letter options first
		sort.Slice(flags, func(i, j int) bool {
			if len(flags[i].Name) != len(flags[j].Name) {
				return len(flags[i].Name) < len(flags[j].Name)
			}
			return flags[i].Name < flags[j].Name
		})

		// Build flag names display
		var flagNames []string
		var flagType string
		var defaultValue string

		for _, f := range flags {
			// Determine flag type dynamically by checking the value type
			if f.DefValue == "false" || f.DefValue == "true" {
				flagType = "" // bool flag
			} else if f.DefValue != "" {
				// For non-bool flags, determine type based on the value
				if _, err := strconv.Atoi(f.DefValue); err == nil {
					flagType = " int"
				} else if _, err := strconv.ParseFloat(f.DefValue, 64); err == nil {
					flagType = " float64"
				} else {
					flagType = " string"
				}
			} else {
				// Default to string for flags without default values
				flagType = " string"
			}

			// Get default value from the first flag
			if defaultValue == "" && f.DefValue != "" && f.DefValue != "false" {
				defaultValue = f.DefValue
			}

			// Format flag name
			if len(f.Name) == 1 {
				flagNames = append(flagNames, "-"+f.Name)
			} else {
				flagNames = append(flagNames, "--"+f.Name)
			}
		}

		// Display the flags
		fmt.Fprintf(os.Stderr, "  %s%s\n", strings.Join(flagNames, ", "), flagType)
		if defaultValue != "" {
			fmt.Fprintf(os.Stderr, "        %s (default %q)\n", flags[0].Usage, defaultValue)
		} else {
			fmt.Fprintf(os.Stderr, "        %s\n", flags[0].Usage)
		}
	}
}
