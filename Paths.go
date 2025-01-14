package asciichgolangpublic

import (
	"errors"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	"github.com/asciich/asciichgolangpublic/logging"

	aerrors "github.com/asciich/asciichgolangpublic/errors"
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
		return nil, aerrors.TracedErrorNil("pathList")
	}

	if pathFilterOptions == nil {
		return nil, aerrors.TracedErrorNil("pathFilterOptions")
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
				match, err = astrings.MatchesRegex(f, pattern)
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

		filtered = aslices.SortStringSliceAndRemoveDuplicates(newFiltered)
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
				match, err = astrings.MatchesRegex(filepath.Base(f), pattern)
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

		filtered = aslices.SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	if pathFilterOptions.IsMatchBasenamePatternSet() {
		newFiltered := []string{}

		matchBaseNamePatterns, err := pathFilterOptions.GetMatchBasenamePattern()
		if err != nil {
			return nil, err
		}

		for _, pattern := range matchBaseNamePatterns {
			for _, f := range filtered {
				match, err := astrings.MatchesRegex(filepath.Base(f), pattern)
				if err != nil {
					return nil, err
				}

				if match {
					newFiltered = append(newFiltered, f)
				}
			}
		}

		filtered = aslices.SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	sort.Strings(filtered)

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

	return aerrors.TracedErrorf("path '%s' is not absolute", path)
}

func (p *PathsService) CheckRelativePath(path string) (err error) {
	if p.IsRelativePath(path) {
		return nil
	}

	return aerrors.TracedErrorf("path '%s' is not relative", path)
}

func (p *PathsService) GetAbsolutePath(path string) (absolutePath string, err error) {
	if path == "" {
		return "", aerrors.TracedErrorEmptyString("path")
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
		return "", aerrors.TracedErrorEmptyString("inputPath")
	}

	if inputPath == "/" {
		return "", aerrors.TracedErrorf("%w: '%s'", ErrPathHasNoParentDirectory, err)
	}

	return filepath.Dir(inputPath), nil
}

func (p *PathsService) GetRelativePathTo(absolutePath string, relativeTo string) (relativePath string, err error) {
	if absolutePath == "" {
		return "", aerrors.TracedErrorEmptyString("absolutePath")
	}

	if relativeTo == "" {
		return "", aerrors.TracedErrorEmptyString("relatvieTo")
	}

	err = Paths().CheckAbsolutePath(absolutePath)
	if err != nil {
		return "", err
	}

	err = Paths().CheckAbsolutePath(relativeTo)
	if err != nil {
		return "", err
	}

	relativeToDirPath := astrings.EnsureSuffix(relativeTo, "/")

	if !strings.HasPrefix(absolutePath, relativeToDirPath) {
		return "", aerrors.TracedErrorf(
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
		return "", aerrors.TracedErrorf("relativePath is empty string after evaluation")
	}

	return relativePath, nil
}

func (p *PathsService) GetRelativePathsTo(absolutePaths []string, relativeTo string) (relativePaths []string, err error) {
	if absolutePaths == nil {
		return nil, aerrors.TracedErrorNil("absoultePaths")
	}

	if relativeTo == "" {
		return nil, aerrors.TracedErrorEmptyString("relativeTo")
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
		logging.LogGoErrorFatal(err)
	}
}

func (p *PathsService) MustCheckRelativePath(path string) {
	err := p.CheckRelativePath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PathsService) MustFilterPaths(pathList []string, pathFilterOptions PathFilterOptions) (filtered []string) {
	filtered, err := p.FilterPaths(pathList, pathFilterOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return filtered
}

func (p *PathsService) MustGetAbsolutePath(path string) (absolutePath string) {
	absolutePath, err := p.GetAbsolutePath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return absolutePath
}

func (p *PathsService) MustGetDirPath(inputPath string) (dirPath string) {
	dirPath, err := p.GetDirPath(inputPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dirPath
}

func (p *PathsService) MustGetRelativePathTo(absolutePath string, relativeTo string) (relativePath string) {
	relativePath, err := p.GetRelativePathTo(absolutePath, relativeTo)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return relativePath
}

func (p *PathsService) MustGetRelativePathsTo(absolutePaths []string, relativeTo string) (relativePaths []string) {
	relativePaths, err := p.GetRelativePathsTo(absolutePaths, relativeTo)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return relativePaths
}
