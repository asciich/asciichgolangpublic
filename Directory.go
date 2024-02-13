package asciichgolangpublic

type Directory interface {
	Create(verbose bool) (err error)
	Delete(verbose bool) (err error)
	Exists() (exists bool, err error)
	GetLocalPath() (localPath string, err error)
	MustCreate(verbose bool)
	MustDelete(verbose bool)
	MustExists() (exists bool)
	MustGetLocalPath() (localPath string)
}
