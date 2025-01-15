package pathsutils

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var ErrPathHasNoParentDirectory = errors.New("path has no parent directory")

// Filter the given path list.
func FilterPaths(pathList []string, pathFilterOptions parameteroptions.PathFilterOptions) (filtered []string, err error) {
	if pathList == nil {
		return nil, tracederrors.TracedErrorNil("pathList")
	}

	if pathFilterOptions == nil {
		return nil, tracederrors.TracedErrorNil("pathFilterOptions")
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
				match, err = stringsutils.MatchesRegex(f, pattern)
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

		filtered = slicesutils.SortStringSliceAndRemoveDuplicates(newFiltered)
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
				match, err = stringsutils.MatchesRegex(filepath.Base(f), pattern)
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

		filtered = slicesutils.SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	if pathFilterOptions.IsMatchBasenamePatternSet() {
		newFiltered := []string{}

		matchBaseNamePatterns, err := pathFilterOptions.GetMatchBasenamePattern()
		if err != nil {
			return nil, err
		}

		for _, pattern := range matchBaseNamePatterns {
			for _, f := range filtered {
				match, err := stringsutils.MatchesRegex(filepath.Base(f), pattern)
				if err != nil {
					return nil, err
				}

				if match {
					newFiltered = append(newFiltered, f)
				}
			}
		}

		filtered = slicesutils.SortStringSliceAndRemoveDuplicates(newFiltered)
	}

	sort.Strings(filtered)

	return filtered, nil
}

// Returns true if path is a relative path.
// An empty string as path will always be false.
func IsRelativePath(path string) (isRelative bool) {
	if path == "" {
		return false
	}

	if IsAbsolutePath(path) {
		return false
	}

	return true
}

// Returns true if path is an absolute path.
// An empty string as path will always be false.
func IsAbsolutePath(path string) (isAbsolute bool) {
	if path == "" {
		return false
	}

	if strings.HasPrefix(path, "/") {
		return true
	}

	re := regexp.MustCompile(`^[a-zA-Z]\:\\`)
	return re.Match([]byte(path))
}

func CheckAbsolutePath(path string) (err error) {
	if IsAbsolutePath(path) {
		return nil
	}

	return tracederrors.TracedErrorf("path '%s' is not absolute", path)
}

func CheckRelativePath(path string) (err error) {
	if IsRelativePath(path) {
		return nil
	}

	return tracederrors.TracedErrorf("path '%s' is not relative", path)
}

func GetAbsolutePath(path string) (absolutePath string, err error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	if IsRelativePath(path) {
		workingDirectoryPath, err := os.Getwd()
		if err != nil {
			return "", err
		}

		path = filepath.Join(workingDirectoryPath, path)
	}

	return path, nil
}

func GetDirPath(inputPath string) (dirPath string, err error) {
	if inputPath == "" {
		return "", tracederrors.TracedErrorEmptyString("inputPath")
	}

	if inputPath == "/" {
		return "", tracederrors.TracedErrorf("%w: '%s'", ErrPathHasNoParentDirectory, err)
	}

	return filepath.Dir(inputPath), nil
}

func GetRelativePathTo(absolutePath string, relativeTo string) (relativePath string, err error) {
	if absolutePath == "" {
		return "", tracederrors.TracedErrorEmptyString("absolutePath")
	}

	if relativeTo == "" {
		return "", tracederrors.TracedErrorEmptyString("relatvieTo")
	}

	err = CheckAbsolutePath(absolutePath)
	if err != nil {
		return "", err
	}

	err = CheckAbsolutePath(relativeTo)
	if err != nil {
		return "", err
	}

	relativeToDirPath := stringsutils.EnsureSuffix(relativeTo, "/")

	if !strings.HasPrefix(absolutePath, relativeToDirPath) {
		return "", tracederrors.TracedErrorf(
			"Only implemented for sub directories but '%s' is not a subdirectory of '%s'",
			absolutePath,
			relativeToDirPath,
		)
	}

	relativePath = strings.TrimPrefix(absolutePath, relativeToDirPath)

	err = CheckRelativePath(relativePath)
	if err != nil {
		return "", err
	}

	if relativePath == "" {
		return "", tracederrors.TracedErrorf("relativePath is empty string after evaluation")
	}

	return relativePath, nil
}

func GetRelativePathsTo(absolutePaths []string, relativeTo string) (relativePaths []string, err error) {
	if absolutePaths == nil {
		return nil, tracederrors.TracedErrorNil("absoultePaths")
	}

	if relativeTo == "" {
		return nil, tracederrors.TracedErrorEmptyString("relativeTo")
	}

	relativePaths = []string{}
	for _, path := range absolutePaths {
		r, err := GetRelativePathTo(path, relativeTo)
		if err != nil {
			return nil, err
		}

		relativePaths = append(relativePaths, r)
	}

	return relativePaths, nil
}

func MustCheckAbsolutePath(path string) {
	err := CheckAbsolutePath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustCheckRelativePath(path string) {
	err := CheckRelativePath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustFilterPaths(pathList []string, pathFilterOptions parameteroptions.PathFilterOptions) (filtered []string) {
	filtered, err := FilterPaths(pathList, pathFilterOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return filtered
}

func MustGetAbsolutePath(path string) (absolutePath string) {
	absolutePath, err := GetAbsolutePath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return absolutePath
}

func MustGetDirPath(inputPath string) (dirPath string) {
	dirPath, err := GetDirPath(inputPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dirPath
}

func MustGetRelativePathTo(absolutePath string, relativeTo string) (relativePath string) {
	relativePath, err := GetRelativePathTo(absolutePath, relativeTo)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return relativePath
}

func MustGetRelativePathsTo(absolutePaths []string, relativeTo string) (relativePaths []string) {
	relativePaths, err := GetRelativePathsTo(absolutePaths, relativeTo)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return relativePaths
}
