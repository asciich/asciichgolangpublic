package parameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ListFileOptions struct {
	MatchBasenamePattern          []string
	ExcludeBasenamePattern        []string
	ExcludePatternWholepath       []string
	ReturnRelativePaths           bool
	OnlyFiles                     bool
	NonRecursive                  bool
	AllowEmptyListIfNoFileIsFound bool
}

func NewListFileOptions() (l *ListFileOptions) {
	return new(ListFileOptions)
}

func (l *ListFileOptions) GetAllowEmptyListIfNoFileIsFound() (allowEmptyListIfNoFileIsFound bool, err error) {

	return l.AllowEmptyListIfNoFileIsFound, nil
}

func (l *ListFileOptions) GetExcludeBasenamePattern() (excludePattern []string, err error) {
	if l.ExcludeBasenamePattern == nil {
		return nil, tracederrors.TracedErrorf("ExcludePattern not set")
	}

	if len(l.ExcludeBasenamePattern) <= 0 {
		return nil, tracederrors.TracedErrorf("ExcludePattern has no elements")
	}

	return l.ExcludeBasenamePattern, nil
}

func (l *ListFileOptions) GetExcludeBasenamePatternOrEmptySliceIfUnset() (excludePattern []string) {
	if len(l.ExcludeBasenamePattern) > 0 {
		return l.ExcludeBasenamePattern
	} else {
		return []string{}
	}
}

func (l *ListFileOptions) GetExcludePatternWholepath() (excludePatternWholepath []string, err error) {
	if l.ExcludePatternWholepath == nil {
		return nil, tracederrors.TracedErrorf("ExcludePatternWholepath not set")
	}

	if len(l.ExcludePatternWholepath) <= 0 {
		return nil, tracederrors.TracedErrorf("ExcludePatternWholepath has no elements")
	}

	return l.ExcludePatternWholepath, nil
}

func (l *ListFileOptions) GetMatchBasenamePattern() (matchPattern []string, err error) {
	if l.MatchBasenamePattern == nil {
		return nil, tracederrors.TracedErrorf("MatchPattern not set")
	}

	if len(l.MatchBasenamePattern) <= 0 {
		return nil, tracederrors.TracedErrorf("MatchPattern has no elements")
	}

	return l.MatchBasenamePattern, nil
}

func (l *ListFileOptions) GetMatchBasenamePatternOrEmptySliceIfUnset() (excludePattern []string) {
	if len(l.MatchBasenamePattern) > 0 {
		return l.MatchBasenamePattern
	} else {
		return []string{}
	}
}

func (l *ListFileOptions) GetNonRecursive() (nonRecursive bool, err error) {

	return l.NonRecursive, nil
}

func (l *ListFileOptions) GetOnlyFiles() (onlyFiles bool, err error) {

	return l.OnlyFiles, nil
}

func (l *ListFileOptions) GetReturnRelativePaths() (returnRelativePaths bool, err error) {

	return l.ReturnRelativePaths, nil
}

func (l *ListFileOptions) IsExcludeBasenamePatternSet() (isSet bool) {
	return len(l.ExcludeBasenamePattern) > 0
}

func (l *ListFileOptions) IsExcludePatternWholepathSet() (isSet bool) {
	return len(l.ExcludePatternWholepath) > 0
}

func (l *ListFileOptions) IsMatchBasenamePatternSet() (isSet bool) {
	return len(l.MatchBasenamePattern) > 0
}
func (l *ListFileOptions) SetAllowEmptyListIfNoFileIsFound(allowEmptyListIfNoFileIsFound bool) (err error) {
	l.AllowEmptyListIfNoFileIsFound = allowEmptyListIfNoFileIsFound

	return nil
}

func (l *ListFileOptions) SetExcludeBasenamePattern(excludeBasenamePattern []string) (err error) {
	if excludeBasenamePattern == nil {
		return tracederrors.TracedErrorf("excludeBasenamePattern is nil")
	}

	if len(excludeBasenamePattern) <= 0 {
		return tracederrors.TracedErrorf("excludeBasenamePattern has no elements")
	}

	l.ExcludeBasenamePattern = excludeBasenamePattern

	return nil
}

func (l *ListFileOptions) SetExcludePattern(excludePattern []string) (err error) {
	if excludePattern == nil {
		return tracederrors.TracedErrorf("excludePattern is nil")
	}

	if len(excludePattern) <= 0 {
		return tracederrors.TracedErrorf("excludePattern has no elements")
	}

	l.ExcludeBasenamePattern = excludePattern

	return nil
}

func (l *ListFileOptions) SetExcludePatternWholepath(excludePatternWholepath []string) (err error) {
	if excludePatternWholepath == nil {
		return tracederrors.TracedErrorf("excludePatternWholepath is nil")
	}

	if len(excludePatternWholepath) <= 0 {
		return tracederrors.TracedErrorf("excludePatternWholepath has no elements")
	}

	l.ExcludePatternWholepath = excludePatternWholepath

	return nil
}

func (l *ListFileOptions) SetMatchBasenamePattern(matchBasenamePattern []string) (err error) {
	if matchBasenamePattern == nil {
		return tracederrors.TracedErrorf("matchBasenamePattern is nil")
	}

	if len(matchBasenamePattern) <= 0 {
		return tracederrors.TracedErrorf("matchBasenamePattern has no elements")
	}

	l.MatchBasenamePattern = matchBasenamePattern

	return nil
}

func (l *ListFileOptions) SetMatchPattern(matchPattern []string) (err error) {
	if matchPattern == nil {
		return tracederrors.TracedErrorf("matchPattern is nil")
	}

	if len(matchPattern) <= 0 {
		return tracederrors.TracedErrorf("matchPattern has no elements")
	}

	l.MatchBasenamePattern = matchPattern

	return nil
}

func (l *ListFileOptions) SetNonRecursive(nonRecursive bool) (err error) {
	l.NonRecursive = nonRecursive

	return nil
}

func (l *ListFileOptions) SetOnlyFiles(onlyFiles bool) (err error) {
	l.OnlyFiles = onlyFiles

	return nil
}

func (l *ListFileOptions) SetReturnRelativePaths(returnRelativePaths bool) (err error) {
	l.ReturnRelativePaths = returnRelativePaths

	return nil
}

func (o *ListFileOptions) GetDeepCopy() (deepCopy *ListFileOptions) {
	deepCopy = new(ListFileOptions)

	*deepCopy = *o

	return deepCopy
}
