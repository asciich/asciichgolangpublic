package asciichgolangpublic

type Directory interface {
	GetLocalPath() (localPath string, err error)
	MustGetLocalPath() (localPath string)
}
