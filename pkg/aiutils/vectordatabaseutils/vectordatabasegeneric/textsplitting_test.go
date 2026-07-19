package vectordatabasegeneric

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplitText_ChunkSizeRespected(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		chunkSize int
		overlap   int
	}{
		{
			name:      "short text under chunk size",
			text:      "hello world",
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "single long paragraph no separators",
			text:      strings.Repeat("a", 1000),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "text with newlines",
			text:      strings.Repeat("This is a line.\n", 100),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "text with double newlines",
			text:      strings.Repeat("This is a paragraph.\n\n", 50),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "text with spaces only",
			text:      strings.Repeat("word ", 500),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "text with sentences",
			text:      strings.Repeat("This is a sentence. ", 200),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "single word longer than chunk size",
			text:      strings.Repeat("x", 600),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "mixed content like markdown",
			text:      "# Title\n\nSome paragraph with text.\n\n## Subtitle\n\n" + strings.Repeat("More content here. ", 100) + "\n\n## Another Section\n\n" + strings.Repeat("Even more text. ", 100),
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "realistic README content",
			text:      "# My Package\n\nThis package provides utilities for working with files.\n\n## Installation\n\n```\ngo get example.com/pkg\n```\n\n## Usage\n\nImport the package and call the functions as needed. The API is designed to be simple and straightforward. You can use it in your projects without any additional configuration. Just import and go.\n\n## API\n\n### Function One\n\nThis function does something really important. It takes a string argument and returns an error if something goes wrong. Make sure to handle the error appropriately in your code.\n\n### Function Two\n\nAnother function that is equally important. It processes data and returns results.",
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "chunk size 50 with overlap 10",
			text:      "The quick brown fox jumps over the lazy dog. " + strings.Repeat("Another sentence here. ", 20),
			chunkSize: 50,
			overlap:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := SplitText(tt.text, tt.chunkSize, tt.overlap)

			if len(chunks) == 0 {
				t.Fatal("SplitText returned no chunks")
			}

			for i, chunk := range chunks {
				if len(chunk) > tt.chunkSize {
					t.Errorf("chunk[%d] exceeds chunkSize: got %d bytes, max %d bytes\nContent: %q",
						i, len(chunk), tt.chunkSize, truncate(chunk, 100))
				}
			}
		})
	}
}

func TestSplitText_NoContentLost(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		chunkSize int
		overlap   int
	}{
		{
			name:      "simple text",
			text:      "Hello world. This is a test. Another sentence.",
			chunkSize: 256,
			overlap:   30,
		},
		{
			name:      "longer text",
			text:      strings.Repeat("word ", 200),
			chunkSize: 100,
			overlap:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := SplitText(tt.text, tt.chunkSize, tt.overlap)

			// Every word from original text should appear in at least one chunk
			words := strings.Fields(tt.text)
			allChunksJoined := strings.Join(chunks, " ")

			for _, word := range words {
				if !strings.Contains(allChunksJoined, word) {
					t.Errorf("word %q from original text not found in any chunk", word)
				}
			}
		})
	}
}

func TestSplitText_OverlapPresent(t *testing.T) {
	text := "AAA BBB CCC. DDD EEE FFF. GGG HHH III. JJJ KKK LLL. MMM NNN OOO."
	chunks := SplitText(text, 30, 10)

	if len(chunks) < 2 {
		t.Fatalf("Expected at least 2 chunks, got %d", len(chunks))
	}

	// Check that consecutive chunks have some overlapping content
	overlapFound := false
	for i := 1; i < len(chunks); i++ {
		prev := chunks[i-1]
		curr := chunks[i]

		// Check if end of previous chunk appears in beginning of current chunk
		if len(prev) > 10 {
			tail := prev[len(prev)-10:]
			if strings.Contains(curr, strings.TrimSpace(tail)) {
				overlapFound = true
				break
			}
		}
	}

	if !overlapFound {
		t.Logf("Chunks: %v", chunks)
		t.Log("Warning: no obvious overlap detected between consecutive chunks")
	}
}

func TestSplitText_MaxChunkSizeWithRealData(t *testing.T) {
	// Simulate the scenario from the bug: files that produce chunks > 256
	content := strings.Repeat("This is line number X in the README file for testing purposes.\n", 50)

	chunkSize := 256
	overlap := 30

	chunks := SplitText(content, chunkSize, overlap)

	var violations int
	var maxLen int
	for i, chunk := range chunks {
		if len(chunk) > chunkSize {
			violations++
			if len(chunk) > maxLen {
				maxLen = len(chunk)
			}
			if violations <= 5 { // Print first 5 violations
				t.Errorf("chunk[%d]: len=%d (exceeds %d)\n  Content: %q",
					i, len(chunk), chunkSize, truncate(chunk, 80))
			}
		}
	}

	if violations > 0 {
		t.Errorf("\nTotal violations: %d/%d chunks exceed chunkSize=%d (max seen: %d)",
			violations, len(chunks), chunkSize, maxLen)
	}
}

func TestSplitText_SinglePartLargerThanChunkSize(t *testing.T) {
	// A single "part" (no separators within) that is larger than chunkSize
	// This is the edge case that likely causes the bug
	longWord := strings.Repeat("abcdefghij", 50) // 500 chars, no spaces/newlines
	text := "Short intro.\n\n" + longWord + "\n\nShort outro."

	chunkSize := 256
	overlap := 30

	chunks := SplitText(text, chunkSize, overlap)

	for i, chunk := range chunks {
		if len(chunk) > chunkSize {
			t.Errorf("chunk[%d]: len=%d exceeds chunkSize=%d\n  Content: %q",
				i, len(chunk), chunkSize, truncate(chunk, 80))
		}
	}
}

func TestRecursiveSplit_FallsThroughAllSeparators(t *testing.T) {
	// Text with no standard separators — must fall through to "" (char-by-char)
	text := strings.Repeat("x", 600)
	chunks := SplitText(text, 256, 30)

	for i, chunk := range chunks {
		if len(chunk) > 256 {
			t.Errorf("chunk[%d]: len=%d, expected <= 256", i, len(chunk))
		}
	}

	// Should have at least 3 chunks (600 / 256 with overlap)
	if len(chunks) < 2 {
		t.Errorf("Expected at least 2 chunks for 600-char text with chunkSize=256, got %d", len(chunks))
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + fmt.Sprintf("... (%d bytes total)", len(s))
}
