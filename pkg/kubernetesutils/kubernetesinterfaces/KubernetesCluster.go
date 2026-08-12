package kubernetesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

type KubernetesCluster interface {
	CheckAccessible(ctx context.Context) error
	CheckNamespaceByNameExists(ctx context.Context, namespaceName string) error
	CheckSecretByNameExists(ctx context.Context, namespaceName string, secretName string) error
	ConfigMapByNameExists(ctx context.Context, namespaceName string, configMapName string) (exists bool, err error)
	CreateConfigMap(ctx context.Context, namespaceName string, configMapName string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap ConfigMap, err error)
	CreateDeployment(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (Deployment, error)
	CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace Namespace, err error)
	CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (Object, error)
	CreatePod(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (Pod, error)
	CreateReplicaSet(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (ReplicaSet, error)
	CreateRole(ctx context.Context, namespaceName string, createOptions *kubernetesparameteroptions.CreateRoleOptions) (createdRole Role, err error)
	CreateClusterRole(ctx context.Context, createOptions *kubernetesparameteroptions.CreateClusterRoleOptions) (createdClusterRole ClusterRole, err error)
	CreateRoleBinding(ctx context.Context, namespaceName string, createOptions *kubernetesparameteroptions.CreateRoleBindingOptions) (createdRoleBinding RoleBinding, err error)
	CreateClusterRoleBinding(ctx context.Context, createOptions *kubernetesparameteroptions.CreateClusterRoleBindingOptions) (createdClusterRoleBinding ClusterRoleBinding, err error)
	CreateSecret(ctx context.Context, namespaceName string, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret Secret, err error)
	DeleteDeploymentByNames(ctx context.Context, namespaceName string, deploymentName string) error
	DeleteNamespaceByName(ctx context.Context, namespaceName string) (err error)
	DeletePodByNames(ctx context.Context, namespaceName string, podName string) error
	DeleteReplicaSetByNames(ctx context.Context, namespaceName string, replicaSetName string) error
	DeleteRoleByName(ctx context.Context, namespaceName string, roleName string) (err error)
	DeleteClusterRoleByName(ctx context.Context, roleName string) (err error)
	DeleteRoleBindingByName(ctx context.Context, namespaceName string, roleBindingName string) (err error)
	DeleteClusterRoleBindingByName(ctx context.Context, roleBindingName string) (err error)
	DeleteSecretByName(ctx context.Context, namespaceName string, secretName string) (err error)
	DeploymentByNameExists(ctx context.Context, namespaceName string, deploymentName string) (bool, error)
	CheckDeploymentByNameExists(ctx context.Context, namespaceName string, deploymentName string) error
	GetDeploymentByNames(namespaceName string, deploymentName string) (Deployment, error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetNamespaceByName(name string) (namespace Namespace, err error)
	GetObjectByNames(objectName string, kind string, namespaceName string) (object Object, err error)
	GetPodByNames(namespaceName string, podName string) (Pod, error)
	GetReplicaSetByNames(namespaceName string, replicaSetName string) (ReplicaSet, error)
	GetRoleByName(namespaceName string, roleName string) (Role, error)
	GetClusterRoleByName(roleName string) (ClusterRole, error)
	GetRoleBindingByName(namespaceName string, roleBindingName string) (RoleBinding, error)
	GetClusterRoleBindingByName(roleBindingName string) (ClusterRoleBinding, error)
	ListKindNames(ctx context.Context) ([]string, error)
	ListNamespaces(ctx context.Context) (namespaces []Namespace, err error)
	ListNamespaceNames(ctx context.Context) ([]string, error)
	ListNodeNames(ctx context.Context) ([]string, error)
	ListObjects(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objects []Object, err error)
	ListObjectNames(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objectNames []string, err error)
	NamespaceByNameExists(ctx context.Context, namespaceName string) (bool, error)
	PodByNameExists(ctx context.Context, namespaceName string, podName string) (bool, error)
	CheckPodByNameExists(ctx context.Context, namespaceName string, podName string) error
	ReadSecret(ctx context.Context, namespaceName string, secretName string) (map[string][]byte, error)
	ReplicaSetByNameExists(ctx context.Context, namespaceName string, replicaSetName string) (bool, error)
	CheckReplicaSetByNameExists(ctx context.Context, namespaceName string, replicaSetName string) error
	RoleByNameExists(ctx context.Context, namespaceName string, roleName string) (exists bool, err error)
	ClusterRoleByNameExists(ctx context.Context, roleName string) (exists bool, err error)
	ListClusterRoleNames(ctx context.Context) ([]string, error)
	RoleBindingByNameExists(ctx context.Context, namespaceName string, roleBindingName string) (exists bool, err error)
	ClusterRoleBindingByNameExists(ctx context.Context, roleBindingName string) (exists bool, err error)
	ListClusterRoleBindingNames(ctx context.Context) ([]string, error)
	RunCommandInTemporaryPod(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (*commandoutput.CommandOutput, error)
	SecretByNameExists(ctx context.Context, namespaceName string, secretName string) (exists bool, err error)
	ValidateSSHKeyInSecret(ctx context.Context, options *kubernetesparameteroptions.ValidateSshKeyInSecretOptions) (bool, error)
	WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.WaitForPodsOptions) error
	WhoAmI(ctx context.Context) (*kubernetesimplementationindependend.UserInfo, error)
	CronJobByNameExists(ctx context.Context, namespaceName string, cronJobName string) (exists bool, err error)
	CheckCronJobByNameExists(ctx context.Context, namespaceName string, cronJobName string) error
	CreateCronJob(ctx context.Context, namespaceName string, cronJobName string, schedule string, image string, command []string, labels map[string]string) (CronJob, error)
	DeleteCronJobByName(ctx context.Context, namespaceName string, cronJobName string) (err error)
}
