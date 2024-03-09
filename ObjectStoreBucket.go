package asciichgolangpublic

type ObjectStoreBucket interface {
	Exists() (exists bool, err error)
	MustExists() (exists bool)
}
