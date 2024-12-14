package asciichgolangpublic

// Options to filter a list of paths.
type PathFilterOptions interface {
	// Filters for matching base name:
	GetMatchBasenamePattern() (matchPattern []string, err error)
	IsMatchBasenamePatternSet() (isSet bool)

	// Filters for excluding base name:
	GetExcludeBasenamePattern() (excludePattern []string, err error)
	IsExcludeBasenamePatternSet() (isSet bool)

	// Filters for exluding wholepaths:
	GetExcludePatternWholepath() (excludePattern []string, err error)
	IsExcludePatternWholepathSet() (isSet bool)
}
