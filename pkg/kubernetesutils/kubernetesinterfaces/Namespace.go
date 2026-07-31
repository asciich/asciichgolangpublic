package kubernetesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

type Namespace interface {
	ConfigMapByNameExists(ctx context.Context, name string) (exits bool, err error)
	Create(ctx context.Context) (err error)
	CreateConfigMap(ctx context.Context, name string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap ConfigMap, err error)
	CreateDeployment(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (Deployment, error)
	CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (Object, error)
	CreatePod(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (Pod, error)
	CreateReplicaSet(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (ReplicaSet, error)
	CreateRole(ctx context.Context, createOptions *kubernetesparameteroptions.CreateRoleOptions) (createdRole Role, err error)
	CreateSecret(ctx context.Context, name string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret Secret, err error)
	DeleteConfigMapByName(ctx context.Context, name string) (err error)
	DeleteDeploymentByName(ctx context.Context, name string) (err error)
	DeletePodByName(ctx context.Context, name string) (err error)
	DeleteReplicaSetByName(ctx context.Context, name string) (err error)
	DeleteRoleByName(ctx context.Context, name string) (err error)
	DeleteSecretByName(ctx context.Context, name string) (err error)
	DeploymentByNameExists(ctx context.Context, deploymentName string) (bool, error)
	Exists(ctx context.Context) (bool, error)
	GetClusterName() (clusterName string, err error)
	GetConfigMapByName(name string) (configMap ConfigMap, err error)
	GetDeploymentByName(name string) (Deployment, error)
	GetKubernetesCluster() (KubernetesCluster, error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetObjectByNames(objectName string, objectType string) (Object, error)
	GetPodByName(name string) (Pod, error)
	GetReplicaSetByName(name string) (ReplicaSet, error)
	GetRoleByName(name string) (Role, error)
	GetSecretByName(name string) (Secret, error)
	ListConfigMapNames(ctx context.Context) ([]string, error)
	ListDeploymentNames(ctx context.Context) ([]string, error)
	ListObjectNames(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objectNames []string, err error)
	ListPodNames(ctx context.Context) ([]string, error)
	ListReplicaSetNames(ctx context.Context) ([]string, error)
	ListRoleNames(ctx context.Context) ([]string, error)
	ListSecrets(ctx context.Context) ([]Secret, error)
	ListSecretNames(ctx context.Context) ([]string, error)
	PodByNameExists(ctx context.Context, podName string) (bool, error)
	ReplicaSetByNameExists(ctx context.Context, replicaSetName string) (bool, error)
	RoleByNameExists(ctx context.Context, name string) (exists bool, err error)
	SecretByNameExists(ctx context.Context, name string) (exits bool, err error)
	WatchConfigMap(ctx context.Context, name string, onCreate func(ConfigMap), onUpdate func(ConfigMap), onDelete func(ConfigMap)) error
	WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, options *kubernetesparameteroptions.WaitForPodsOptions) error
	CronJobByNameExists(ctx context.Context, cronJobName string) (bool, error)
	CreateCronJob(ctx context.Context, cronJobName string, schedule string, image string, command []string, labels map[string]string) (CronJob, error)
	DeleteCronJobByName(ctx context.Context, cronJobName string) (err error)
	GetCronJobByName(name string) (CronJob, error)
	ListCronJobNames(ctx context.Context) ([]string, error)
}
