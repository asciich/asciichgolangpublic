package kubernetes

type Namespace interface {
	GetName() (name string, err error)
	MustGetName() (name string)
}
