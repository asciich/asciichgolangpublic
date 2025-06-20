package kubernetesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

type Namespace interface {
	ConfigMapByNameExists(ctx context.Context, name string) (exits bool, err error)
	Create(ctx context.Context) (err error)
	CreateConfigMap(ctx context.Context, name string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap ConfigMap, err error)
	CreateRole(ctx context.Context, createOptions *kubernetesparameteroptions.CreateRoleOptions) (createdRole Role, err error)
	CreateSecret(ctx context.Context, name string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret Secret, err error)
	CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (Object, error)
	DeleteConfigMapByName(ctx context.Context, name string) (err error)
	DeleteRoleByName(ctx context.Context, name string) (err error)
	DeleteSecretByName(ctx context.Context, name string) (err error)
	GetClusterName() (clusterName string, err error)
	GetConfigMapByName(name string) (configMap ConfigMap, err error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetObjectByNames(objectName string, objectType string) (object Object, err error)
	GetRoleByName(name string) (role Role, err error)
	GetSecretByName(name string) (secret Secret, err error)
	ListRoleNames(ctx context.Context) (roleNames []string, err error)
	RoleByNameExists(ctx context.Context, name string) (exists bool, err error)
	SecretByNameExists(ctx context.Context, name string) (exits bool, err error)
	WatchConfigMap(ctx context.Context, name string, onCreate func(ConfigMap), onUpdate func(ConfigMap), onDelete func(ConfigMap)) error
	WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, options *kubernetesparameteroptions.WaitForPodsOptions) error
}
