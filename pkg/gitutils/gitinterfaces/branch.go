package gitinterfaces

type Branch interface {
	GetDeepCopy() Branch
	GetName() (string, error)

	SetName(name string) error
}
