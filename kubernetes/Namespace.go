package kubernetes

type Namespace interface {
	CreateRole(createOptions *CreateRoleOptions) (createdRole Role, err error)
	DeleteRoleByName(name string, verbose bool) (err error)
	GetClusterName() (clusterName string, err error)
	GetName() (name string, err error)
	GetRoleByName(name string) (role Role, err error)
	ListRoleNames(verbose bool) (roleNames []string, err error)
	MustCreateRole(createOptions *CreateRoleOptions) (createdRole Role)
	MustDeleteRoleByName(name string, verbose bool)
	MustGetClusterName() (clusterName string)
	MustGetName() (name string)
	MustGetRoleByName(name string) (role Role)
	MustListRoleNames(verbose bool) (roleNames []string)
	MustRoleByNameExists(name string, verbose bool) (exists bool)
	RoleByNameExists(name string, verbose bool) (exists bool, err error)
}
