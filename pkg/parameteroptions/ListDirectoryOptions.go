package parameteroptions

type ListDirectoryOptions struct {
	// Enable recursive file and/or directory listing:
	Recursive bool

	// Return paths relative to the directory to list:
	ReturnRelativePaths bool
}

func NewListDirectoryOptions() (l *ListDirectoryOptions) {
	return new(ListDirectoryOptions)
}

func (l *ListDirectoryOptions) GetRecursive() (recursive bool) {

	return l.Recursive
}

func (l *ListDirectoryOptions) GetReturnRelativePaths() (returnRelativePaths bool) {

	return l.ReturnRelativePaths
}

func (l *ListDirectoryOptions) SetRecursive(recursive bool) {
	l.Recursive = recursive
}

func (l *ListDirectoryOptions) SetReturnRelativePaths(returnRelativePaths bool) {
	l.ReturnRelativePaths = returnRelativePaths
}
