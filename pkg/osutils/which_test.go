package osutils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
		{"find_nonexistent", "nonexistentcommand12345", false},
		{"find_empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := osutils.Which(tt.command)

			if tt.command == "" {
				require.Error(t, err, "Expected error for empty command")
				return
			}

			if tt.expectFound {
				require.NotEmpty(t, path, "Expected to find '%s'", tt.command)
				require.True(t, filepath.IsAbs(path), "Expected absolute path for '%s', got '%s'", tt.command, path)
			} else {
				require.Empty(t, path, "Expected empty path for '%s', but got '%s'", tt.command, path)
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
				require.Error(t, err, "Expected error for empty command")
				return
			}

			require.NoError(t, err)
			require.GreaterOrEqual(t, len(results), tt.minExpectedCount, "Expected at least %d results for '%s'", tt.minExpectedCount, tt.command)

			for _, path := range results {
				require.True(t, filepath.IsAbs(path), "Expected absolute path, got '%s'", path)
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

	allPaths, err := osutils.WhichAll("go")
	require.NoError(t, err)
	require.NotEmpty(t, allPaths, "Expected at least one path for 'go'")

	require.Contains(t, allPaths, path, "Expected WhichAll('go') to include Which('go') result '%s'", path)
}

func TestWhichAllEmptyPathEntry(t *testing.T) {
	origPATH := os.Getenv("PATH")
	defer os.Setenv("PATH", origPATH)
	os.Setenv("PATH", "/usr/bin: :/usr/local/bin")

	_, err := osutils.WhichAll("bash")
	require.NoError(t, err, "WhichAll should handle empty PATH entries gracefully")
}
