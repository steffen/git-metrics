package display

import (
	"testing"
)

func TestCreatePathFootnote(t *testing.T) {
	tests := []struct {
		name                 string
		path                 string
		maxDisplayLength     int
		currentFootnoteCount int
		wantDisplayPath      string
		wantFootnoteIndex    int
	}{
		{
			name:                 "Short path needs no truncation",
			path:                 "README.md",
			maxDisplayLength:     43,
			currentFootnoteCount: 0,
			wantDisplayPath:      "README.md",
			wantFootnoteIndex:    0,
		},
		{
			name:                 "Path exactly at max length",
			path:                 "this-file-name-is-exactly-43-chars-long.jpg",
			maxDisplayLength:     43,
			currentFootnoteCount: 0,
			wantDisplayPath:      "this-file-name-is-exactly-43-chars-long.jpg",
			wantFootnoteIndex:    0,
		},
		{
			name:                 "Path one character over max length",
			path:                 "this-file-name-is-exactly-44-chars-long.jpeg",
			maxDisplayLength:     43,
			currentFootnoteCount: 0,
			wantDisplayPath:      "this-file-name-is-...44-chars-long.jpeg [1]",
			wantFootnoteIndex:    1,
		},
		{
			name:                 "Very long path",
			path:                 "a/very/long/path/that/exceeds/the/limit/for/display/in/the/table/and/should/be/truncated/by/the/tool/very-long-file-name-1.txt",
			maxDisplayLength:     43,
			currentFootnoteCount: 0,
			wantDisplayPath:      "a/very/long/path/t...ng-file-name-1.txt [1]",
			wantFootnoteIndex:    1,
		},
		{
			name:                 "With existing footnotes",
			path:                 "another/long/path/to/truncate.txt",
			maxDisplayLength:     43,
			currentFootnoteCount: 5,
			wantDisplayPath:      "another/long/path/to/truncate.txt",
			wantFootnoteIndex:    0,
		},
		{
			name:                 "Extremely long path with large footnote index",
			path:                 "this/is/an/extremely/long/path/that/would/need/significant/truncation/even/with/a/large/footnote/index.txt",
			maxDisplayLength:     43,
			currentFootnoteCount: 99,
			wantDisplayPath:      "this/is/an/extrem...ootnote/index.txt [100]",
			wantFootnoteIndex:    100,
		},
		{
			name:                 "Small max display length",
			path:                 "README.md",
			maxDisplayLength:     10,
			currentFootnoteCount: 0,
			wantDisplayPath:      "README.md",
			wantFootnoteIndex:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreatePathFootnote(tt.path, tt.maxDisplayLength, tt.currentFootnoteCount)

			if got.DisplayPath != tt.wantDisplayPath {
				t.Errorf("CreatePathFootnote().DisplayPath = %q, want %q", got.DisplayPath, tt.wantDisplayPath)
			}

			if got.Index != tt.wantFootnoteIndex {
				t.Errorf("CreatePathFootnote().Index = %d, want %d", got.Index, tt.wantFootnoteIndex)
			}

			if got.FullPath != tt.path {
				t.Errorf("CreatePathFootnote().FullPath = %q, want %q", got.FullPath, tt.path)
			}

			// Check that the display path doesn't exceed max length
			if len(got.DisplayPath) > tt.maxDisplayLength {
				t.Errorf("Display path length %d exceeds maximum length %d", len(got.DisplayPath), tt.maxDisplayLength)
			}
		})
	}
}
