package asciichgolangpublic

import (
	"path/filepath"
	"regexp"
	"strings"
)

type PathsService struct{}

func NewPathsService() (p *PathsService) {
	return new(PathsService)
}

func Paths() (p *PathsService) {
	return NewPathsService()
}

// Returns true if path is a relative path.
// An empty string as path will always be false.
func (p *PathsService) IsRelativePath(path string) (isRelative bool) {
	if path == "" {
		return false
	}

	if p.IsAbsolutePath(path) {
		return false
	}

	return true
}

// Returns true if path is an absolute path.
// An empty string as path will always be false.
func (p *PathsService) IsAbsolutePath(path string) (isAbsolute bool) {
	if path == "" {
		return false
	}

	if strings.HasPrefix(path, "/") {
		return true
	}

	re := regexp.MustCompile(`^[a-zA-Z]\:\\`)
	return re.Match([]byte(path))
}

func (p *PathsService) GetAbsolutePath(path string) (absolutePath string, err error) {
	if path == "" {
		return "", TracedErrorEmptyString("path")
	}

	if Paths().IsRelativePath(path) {
		workingDirectoryPath, err := OS().GetCurrentWorkingDirectoryAsString()
		if err != nil {
			return "", err
		}

		path = filepath.Join(workingDirectoryPath, path)
	}

	return path, nil
}

func (p *PathsService) MustGetAbsolutePath(path string) (absolutePath string) {
	absolutePath, err := p.GetAbsolutePath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return absolutePath
}
