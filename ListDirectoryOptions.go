package asciichgolangpublic

type ListDirectoryOptions struct {
	Recursive bool
}

func NewListDirectoryOptions() (l *ListDirectoryOptions) {
	return new(ListDirectoryOptions)
}

func (l *ListDirectoryOptions) GetRecursive() (recursive bool) {

	return l.Recursive
}

func (l *ListDirectoryOptions) SetRecursive(recursive bool) {
	l.Recursive = recursive
}
