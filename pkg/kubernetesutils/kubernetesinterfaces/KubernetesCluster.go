package kubernetesinterfaces

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesimplementationindependend"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesparameteroptions"
)

type KubernetesCluster interface {
	CheckAccessible(ctx context.Context) error
	ConfigMapByNameExists(ctx context.Context, namespaceName string, configMapName string) (exists bool, err error)
	CreateConfigMap(ctx context.Context, namespaceName string, configMapName string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap ConfigMap, err error)
	CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace Namespace, err error)
	CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (Object, error)
	CreateSecret(ctx context.Context, namespaceName string, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret Secret, err error)
	DeleteNamespaceByName(ctx context.Context, namespaceName string) (err error)
	DeleteSecretByName(ctx context.Context, namespaceName string, secretName string) (err error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetNamespaceByName(name string) (namespace Namespace, err error)
	GetObjectByNames(objectName string, kind string, namespaceName string) (object Object, err error)
	ListNamespaces(ctx context.Context) (namespaces []Namespace, err error)
	ListNamespaceNames(ctx context.Context) ([]string, error)
	ListNodeNames(ctx context.Context) ([]string, error)
	ListObjects(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objects []Object, err error)
	ListObjectNames(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objectNames []string, err error)
	NamespaceByNameExists(ctx context.Context, namespaceName string) (exists bool, err error)
	ReadSecret(ctx context.Context, namespaceName string, secretName string) (map[string][]byte, error)
	RunCommandInTemporaryPod(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (*commandexecutor.CommandOutput, error)
	SecretByNameExists(ctx context.Context, namespaceName string, secretName string) (exists bool, err error)
	WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.WaitForPodsOptions) error
	WhoAmI(ctx context.Context) (*kubernetesimplementationindependend.UserInfo, error)
}
