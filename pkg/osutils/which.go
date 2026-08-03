package osutils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Which finds the full path of a command by searching the PATH environment variable.
// Returns the absolute path to the executable if found.
// Returns ErrExecutableNotFound if the command is not found.
// This is a Go-native implementation of the Unix 'which' command.
func Which(command string) (string, error) {
	if command == "" {
		return "", tracederrors.TracedErrorEmptyString("command")
	}

	path, err := exec.LookPath(command)
	if err != nil {
		return "", ErrExecutableNotFound
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	return path, nil
}

// WhichAll finds all occurrences of a command in the PATH environment variable.
// Returns a slice of absolute paths to all matching executables.
// Similar to 'which -a' command.
func WhichAll(command string) ([]string, error) {
	if command == "" {
		return []string{}, tracederrors.TracedErrorEmptyString("command")
	}

	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return []string{}, nil
	}

	pathDirs := filepath.SplitList(pathEnv)
	var results []string

	for _, dir := range pathDirs {
		if dir == "" {
			dir = "."
		}

		fullPath := filepath.Join(dir, command)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if info.IsDir() {
			continue
		}

		if isExecutable(fullPath, info) {
			absPath, err := filepath.Abs(fullPath)
			if err != nil {
				continue
			}
			results = append(results, absPath)
		}
	}

	return results, nil
}

func isExecutable(path string, info os.FileInfo) bool {
	if isWindows() {
		return isWindowsExecutable(path)
	}
	mode := info.Mode()
	return mode&0111 != 0
}

func isWindowsExecutable(path string) bool {
	executableExtensions := []string{".com", ".exe", ".bat", ".cmd", ".ps1", ".vbs", ".msc"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, execExt := range executableExtensions {
		if ext == execExt {
			return true
		}
	}
	return ext == ""
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
