package asciichgolangpublic

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrPathHasNoParentDirectory = errors.New("path has no parent directory")

type PathsService struct{}

func NewPathsService() (p *PathsService) {
	return new(PathsService)
}

func Paths() (p *PathsService) {
	return NewPathsService()
}

// Filter the given path list.
func (p *PathsService) FilterPaths(pathList []string, pathFilterOptions PathFilterOptions) (filtered []string, err error) {
	if pathList == nil {
		return nil, TracedErrorNil("pathList")
	}

	if pathFilterOptions == nil {
		return nil, TracedErrorNil("pathFilterOptions")
	}

	filtered = pathList

	if pathFilterOptions.IsExcludePatternWholepathSet() {
		newFiltered := []string{}

		excludePatterns, err := pathFilterOptions.GetExcludePatternWholepath()
		if err != nil {
			return nil, err
		}

		for _, f := range filtered {
			match := false
			for _, pattern := range excludePatterns {
				match, err = Strings().MatchesRegex(f, pattern)
				if err != nil {
					return nil, err
				}

				if match {
					break
				}
			}

			if !match {
				newFiltered = append(newFiltered, f)
			}
		}

		filtered = Slices().SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	if pathFilterOptions.IsExcludeBasenamePatternSet() {
		newFiltered := []string{}

		excludePatterns, err := pathFilterOptions.GetExcludeBasenamePattern()
		if err != nil {
			return nil, err
		}

		for _, f := range filtered {
			match := false
			for _, pattern := range excludePatterns {
				match, err = Strings().MatchesRegex(filepath.Base(f), pattern)
				if err != nil {
					return nil, err
				}

				if match {
					break
				}
			}

			if !match {
				newFiltered = append(newFiltered, f)
			}
		}

		filtered = Slices().SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	if pathFilterOptions.IsMatchBasenamePatternSet() {
		newFiltered := []string{}

		matchBaseNamePatterns, err := pathFilterOptions.GetMatchBasenamePattern()
		if err != nil {
			return nil, err
		}

		for _, pattern := range matchBaseNamePatterns {
			for _, f := range filtered {
				match, err := Strings().MatchesRegex(filepath.Base(f), pattern)
				if err != nil {
					return nil, err
				}

				if match {
					newFiltered = append(newFiltered, f)
				}
			}
		}

		filtered = Slices().SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	filtered = Slices().SortStringSlice(filtered)

	return filtered, nil
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

func (p *PathsService) CheckAbsolutePath(path string) (err error) {
	if p.IsAbsolutePath(path) {
		return nil
	}

	return TracedErrorf("path '%s' is not absolute", path)
}

func (p *PathsService) CheckRelativePath(path string) (err error) {
	if p.IsRelativePath(path) {
		return nil
	}

	return TracedErrorf("path '%s' is not relative", path)
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

func (p *PathsService) GetDirPath(inputPath string) (dirPath string, err error) {
	if inputPath == "" {
		return "", TracedErrorEmptyString("inputPath")
	}

	if inputPath == "/" {
		return "", TracedErrorf("%w: '%s'", ErrPathHasNoParentDirectory, err)
	}

	return filepath.Dir(inputPath), nil
}

func (p *PathsService) GetRelativePathTo(absolutePath string, relativeTo string) (relativePath string, err error) {
	if absolutePath == "" {
		return "", TracedErrorEmptyString("absolutePath")
	}

	if relativeTo == "" {
		return "", TracedErrorEmptyString("relatvieTo")
	}

	err = Paths().CheckAbsolutePath(absolutePath)
	if err != nil {
		return "", err
	}

	err = Paths().CheckAbsolutePath(relativeTo)
	if err != nil {
		return "", err
	}

	relativeToDirPath := Strings().EnsureSuffix(relativeTo, "/")

	if !strings.HasPrefix(absolutePath, relativeToDirPath) {
		return "", TracedErrorf(
			"Only implemented for sub directories but '%s' is not a subdirectory of '%s'",
			absolutePath,
			relativeToDirPath,
		)
	}

	relativePath = strings.TrimPrefix(absolutePath, relativeToDirPath)

	err = Paths().CheckRelativePath(relativePath)
	if err != nil {
		return "", err
	}

	if relativePath == "" {
		return "", TracedErrorf("relativePath is empty string after evaluation")
	}

	return relativePath, nil
}

func (p *PathsService) GetRelativePathsTo(absolutePaths []string, relativeTo string) (relativePaths []string, err error) {
	if absolutePaths == nil {
		return nil, TracedErrorNil("absoultePaths")
	}

	if relativeTo == "" {
		return nil, TracedErrorEmptyString("relativeTo")
	}

	relativePaths = []string{}
	for _, path := range absolutePaths {
		r, err := p.GetRelativePathTo(path, relativeTo)
		if err != nil {
			return nil, err
		}

		relativePaths = append(relativePaths, r)
	}

	return relativePaths, nil
}

func (p *PathsService) MustCheckAbsolutePath(path string) {
	err := p.CheckAbsolutePath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *PathsService) MustCheckRelativePath(path string) {
	err := p.CheckRelativePath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *PathsService) MustFilterPaths(pathList []string, pathFilterOptions PathFilterOptions) (filtered []string) {
	filtered, err := p.FilterPaths(pathList, pathFilterOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filtered
}

func (p *PathsService) MustGetAbsolutePath(path string) (absolutePath string) {
	absolutePath, err := p.GetAbsolutePath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return absolutePath
}

func (p *PathsService) MustGetDirPath(inputPath string) (dirPath string) {
	dirPath, err := p.GetDirPath(inputPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dirPath
}

func (p *PathsService) MustGetRelativePathTo(absolutePath string, relativeTo string) (relativePath string) {
	relativePath, err := p.GetRelativePathTo(absolutePath, relativeTo)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return relativePath
}

func (p *PathsService) MustGetRelativePathsTo(absolutePaths []string, relativeTo string) (relativePaths []string) {
	relativePaths, err := p.GetRelativePathsTo(absolutePaths, relativeTo)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return relativePaths
}
