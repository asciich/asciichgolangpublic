package github.com/asciich/asciichgolangpublic

type Directory interface {
	GetLocalPath() (localPath string, err error)
	MustGetLocalPath() (localPath string)
}
