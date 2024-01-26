package asciichgolangpublic

// A File represents any kind of file regardless if a local file or a remote file.
type File interface {
	Exists() (exists bool, err error)
	GetLocalPath() (localPath string, err error)
	GetUriAsString() (uri string, err error)
	MustGetLocalPath() (localPath string)
	MustGetUriAsString() (uri string)
}
