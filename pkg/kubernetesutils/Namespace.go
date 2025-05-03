package kubernetesutils

import "context"

type Namespace interface {
	Create(ctx context.Context) (err error)
	CreateRole(createOptions *CreateRoleOptions) (createdRole Role, err error)
	DeleteRoleByName(ctx context.Context, name string) (err error)
	GetClusterName() (clusterName string, err error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetResourceByNames(resourceName string, resourceType string) (resource Resource, err error)
	GetRoleByName(name string) (role Role, err error)
	ListRoleNames(ctx context.Context) (roleNames []string, err error)
	RoleByNameExists(ctx context.Context, name string) (exists bool, err error)
}
