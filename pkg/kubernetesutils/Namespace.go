package kubernetesutils

type Namespace interface {
	Create(verbose bool) (err error)
	CreateRole(createOptions *CreateRoleOptions) (createdRole Role, err error)
	DeleteRoleByName(name string, verbose bool) (err error)
	GetClusterName() (clusterName string, err error)
	GetKubectlContext(verbose bool) (contextName string, err error)
	GetName() (name string, err error)
	GetResourceByNames(resourceName string, resourceType string) (resource Resource, err error)
	GetRoleByName(name string) (role Role, err error)
	ListRoleNames(verbose bool) (roleNames []string, err error)
	MustCreate(verbose bool)
	MustCreateRole(createOptions *CreateRoleOptions) (createdRole Role)
	MustDeleteRoleByName(name string, verbose bool)
	MustGetClusterName() (clusterName string)
	MustGetKubectlContext(verbose bool) (contextName string)
	MustGetName() (name string)
	MustGetRoleByName(name string) (role Role)
	MustGetResourceByNames(resourceName string, resourceType string) (resource Resource)
	MustListRoleNames(verbose bool) (roleNames []string)
	MustRoleByNameExists(name string, verbose bool) (exists bool)
	RoleByNameExists(name string, verbose bool) (exists bool, err error)
}
