package kubernetesutils

// a generic representation of a kubernetes resource like a pod, ingress, role...
type Resource interface {
	CreateByYamlString(roleYaml string, verbose bool) (err error)
	Delete(verbose bool) (err error)
	Exists(verbose bool) (exists bool, err error)
	GetAsYamlString() (yamlString string, err error)
	MustCreateByYamlString(roleYaml string, verbose bool)
	MustDelete(verbose bool)
	MustExists(verbose bool) (exists bool)
	MustGetAsYamlString() (yamlString string)
}
