package asciichgolangpublic

type ListDirectoryOptions struct {
	Recursive bool
	Verbose   bool
}

func NewListDirectoryOptions() (l *ListDirectoryOptions) {
	return new(ListDirectoryOptions)
}

func (l *ListDirectoryOptions) GetRecursive() (recursive bool) {

	return l.Recursive
}

func (l *ListDirectoryOptions) GetVerbose() (verbose bool) {

	return l.Verbose
}

func (l *ListDirectoryOptions) SetRecursive(recursive bool) {
	l.Recursive = recursive
}

func (l *ListDirectoryOptions) SetVerbose(verbose bool) {
	l.Verbose = verbose
}
