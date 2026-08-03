package osutils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/osutils"
)

func TestWhich(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectFound bool
	}{
		{"find_go", "go", true},
		{"find_bash", "bash", true},
		{"find_sh", "sh", true},
		{"find_nonexistent", "nonexistentcommand12345", false}, // TODO adjust, this must return an error ErrExecutableNotFound when the binary is not found.
		{"find_empty", "", false},                              // TODO adjust, this must return an error tracedErrors.TracedErrorEmptyString
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := osutils.Which(tt.command)

			if tt.command == "" {
				if err == nil {
					t.Errorf("Expected error for empty command, got nil")
				}
				return
			}

			if tt.expectFound {
				if path == "" {
					t.Errorf("Expected to find '%s', but got empty path", tt.command)
				}
				if !filepath.IsAbs(path) {
					t.Errorf("Expected absolute path for '%s', got '%s'", tt.command, path)
				}
			} else {
				if path != "" {
					t.Errorf("Expected empty path for '%s', but got '%s'", tt.command, path)
				}
			}
		})
	}
}

func TestWhichAll(t *testing.T) {
	tests := []struct {
		name             string
		command          string
		minExpectedCount int
		expectError      bool
	}{
		{"find_all_bash", "bash", 1, false},
		{"find_all_sh", "sh", 1, false},
		{"find_all_nonexistent", "nonexistentcommand12345", 0, false},
		{"find_all_empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := osutils.WhichAll(tt.command)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for empty command, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(results) < tt.minExpectedCount {
				t.Errorf("Expected at least %d results for '%s', got %d", tt.minExpectedCount, tt.command, len(results))
			}

			for _, path := range results {
				if !filepath.IsAbs(path) {
					t.Errorf("Expected absolute path, got '%s'", path)
				}
			}
		})
	}
}

func TestWhichAllReturnsAllOccurrences(t *testing.T) {
	path, err := osutils.Which("go")
	if err != nil {
		t.Skip("Go not found, skipping test")
		return
	}

	if path == "" {
		t.Skip("Go not found in PATH, skipping test")
		return
	}

	allPaths, _ := osutils.WhichAll("go")
	if len(allPaths) == 0 {
		t.Errorf("Expected at least one path for 'go', got 0 paths")
	}

	verified := false
	for _, p := range allPaths {
		if p == path {
			verified = true
			break
		}
	}

	if !verified {
		t.Errorf("Expected WhichAll( 'go') to include Which( 'go') result '%s', but got %v", path, allPaths)
	}
}

func TestWhichAllEmptyPathEntry(t *testing.T) {
	origPATH := os.Getenv("PATH")
	defer os.Setenv("PATH", origPATH)
	os.Setenv("PATH", "/usr/bin: :/usr/local/bin")
	_, err := osutils.WhichAll("bash")
	if err != nil {
		t.Errorf("WhichAll should handle empty PATH entries gracefully, got error: %v", err)
	}
}
