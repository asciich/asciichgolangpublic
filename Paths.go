package github.com/asciichgolangpublic

import "strings"

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

	return !strings.HasPrefix(path, "/")
}

// Returns true if path is an absolute path.
// An empty string as path will always be false.
func (p *PathsService) IsAbsolutePath(path string) (isRelative bool) {
	if path == "" {
		return false
	}

	return strings.HasPrefix(path, "/")
}
