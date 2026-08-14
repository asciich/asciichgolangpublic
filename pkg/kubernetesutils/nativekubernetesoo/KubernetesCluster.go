package nativekubernetesoo

import (
	"context"
	"sort"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeconfigutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"

	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type NativeKubernetesCluster struct {
	name   string
	config *rest.Config

	// client caches:
	clientSetCache     *kubernetes.Clientset
	dynamicClientCache *dynamic.DynamicClient
}

func GetClusterByName(ctx context.Context, clusterName string) (*NativeKubernetesCluster, error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	config, err := nativekubernetes.GetConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	return &NativeKubernetesCluster{
		name:   clusterName,
		config: config,
	}, nil
}

func GetDefaultCluster(ctx context.Context) (*NativeKubernetesCluster, error) {
	config, err := nativekubernetes.GetConfig(ctx, "")
	if err != nil {
		return nil, err
	}

	return &NativeKubernetesCluster{
		config: config,
	}, nil
}

func (n *NativeKubernetesCluster) CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace kubernetesinterfaces.Namespace, err error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	err = nativekubernetes.CreateNamespace(ctx, clientSet, namespaceName)
	if err != nil {
		return nil, err
	}

	return n.GetNamespaceByName(namespaceName)
}

func (n *NativeKubernetesCluster) GetDynamicClient() (*dynamic.DynamicClient, error) {
	config, err := n.GetConfig()
	if err != nil {
		return nil, err
	}

	if n.dynamicClientCache == nil {
		var err error
		n.dynamicClientCache, err = dynamic.NewForConfig(config)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Error creating kubernetes dynamic client: %w", err)
		}

	}

	return n.dynamicClientCache, nil
}

func (n *NativeKubernetesCluster) GetConfig() (*rest.Config, error) {
	if n.config == nil {
		return nil, tracederrors.TracedError("config not set")
	}

	return n.config, nil
}

func (n *NativeKubernetesCluster) GetClientSet() (*kubernetes.Clientset, error) {
	config, err := n.GetConfig()
	if err != nil {
		return nil, err
	}

	if n.clientSetCache == nil {
		var err error
		n.clientSetCache, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Error creating Kubernetes clientset: %w", err)
		}

	}

	return n.clientSetCache, nil
}

func (n *NativeKubernetesCluster) DeleteNamespaceByName(ctx context.Context, namespaceName string) (err error) {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	exists, err := n.NamespaceByNameExists(ctx, namespaceName)
	if err != nil {
		return err
	}

	if exists {
		clientset, err := n.GetClientSet()
		if err != nil {
			return err
		}

		deletePolicy := metav1.DeletePropagationForeground // This ensures child objects are deleted before the namespace
		deleteOptions := metav1.DeleteOptions{
			PropagationPolicy:  &deletePolicy,
			GracePeriodSeconds: nil, // Use default graceful termination period
		}

		err = clientset.CoreV1().Namespaces().Delete(ctx, namespaceName, deleteOptions)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete kubernetes namespace '%s': %w", namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' deleted.", namespaceName)

		err = n.WaitUntilNamespaceDeleted(ctx, namespaceName)
		if err != nil {
			return err
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already absent. Skip delete.", namespaceName)
	}

	return nil
}

func (n *NativeKubernetesCluster) GetKubectlContext(ctx context.Context) (contextName string, err error) {
	clusterName, err := n.GetName()
	if err != nil {
		return "", err
	}

	return kubeconfigutils.GetContextNameByClusterName(ctx, clusterName)
}
func (n *NativeKubernetesCluster) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return n.name, nil
}

func (n *NativeKubernetesCluster) GetNamespaceByName(name string) (namespace kubernetesinterfaces.Namespace, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return &NativeNamespace{
		name:              name,
		kubernetesCluster: n,
	}, nil
}

func (n *NativeKubernetesCluster) GetReplicaSetByNames(namespaceName string, replicaSetName string) (kubernetesinterfaces.ReplicaSet, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetReplicaSetByName(replicaSetName)
}

func (n *NativeKubernetesCluster) GetDeploymentByNames(namespaceName string, deploymentName string) (kubernetesinterfaces.Deployment, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetDeploymentByName(deploymentName)
}

func (n *NativeKubernetesCluster) GetPodByNames(namespaceName string, podName string) (kubernetesinterfaces.Pod, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetPodByName(podName)
}

func (n *NativeKubernetesCluster) GetObjectByNames(objectName string, objectType string, namespaceName string) (object kubernetesinterfaces.Object, err error) {
	if objectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectName")
	}

	if objectType == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectType")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetObjectByNames(objectName, objectType)
}
func (n *NativeKubernetesCluster) ListNamespaces(ctx context.Context) (namespaces []kubernetesinterfaces.Namespace, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) ListNamespaceNames(ctx context.Context) (namespaceNames []string, err error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error listing namespaces: %w", err)
	}

	namespaceNames = []string{}
	if len(namespaces.Items) == 0 {
		return nil, tracederrors.TracedErrorf("No namespaces found.")
	} else {
		for _, ns := range namespaces.Items {
			namespaceNames = append(namespaceNames, ns.Name)
		}
	}

	sort.Strings(namespaceNames)

	logging.LogInfoByCtxf(ctx, "Found %d namespaces.", len(namespaceNames))

	return namespaceNames, nil
}
func (n *NativeKubernetesCluster) ListObjects(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objects []kubernetesinterfaces.Object, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) ListObjectNames(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objectNames []string, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) NamespaceByNameExists(ctx context.Context, namespaceName string) (exist bool, err error) {
	if namespaceName == "" {
		return false, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return false, err
	}

	_, err = clientset.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Kubernetes namespace '%s' does not exist.", namespaceName)
			return false, nil
		}
		return false, tracederrors.TracedErrorEmptyString("failed to get namespace: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Kubernetes namespace '%s' exists.", namespaceName)
	return true, nil
}

func (n *NativeKubernetesCluster) WaitUntilNamespaceDeleted(ctx context.Context, namepaceName string) (err error) {
	if namepaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	timeout := time.Second * 60

	logging.LogInfoByCtxf(ctx, "Wait for kubernetes namespace '%s' to be deleted started (timeout = %s).", namepaceName, timeout)

	ctx, _ = context.WithTimeout(ctx, timeout)
	tStart := time.Now()
	for {
		if ctx.Err() != nil {
			return tracederrors.TracedErrorf("Wait until namespace '%s' deleted failed: %w", namepaceName, ctx.Err())
		}

		exists, err := n.NamespaceByNameExists(ctx, namepaceName)
		if err != nil {
			return err
		}

		if exists {
			waitTime := time.Second * 1
			elapsed := time.Since(tStart)
			logging.LogInfoByCtxf(ctx, "Wait another %s until the kubernetes namespace '%s' is deleted (%s/%s).", waitTime, namepaceName, elapsed, timeout)
			time.Sleep(waitTime)
		} else {
			break
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait for kubernetes namespace '%s' to be deleted finished.", namepaceName)

	return nil
}

func (n *NativeKubernetesCluster) WaitUntilNamespaceCreated(ctx context.Context, namepaceName string) (err error) {
	if namepaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	timeout := time.Second * 15

	logging.LogInfoByCtxf(ctx, "Wait for kubernetes namespace '%s' to be created started (timeout = %s).", namepaceName, timeout)

	ctx, _ = context.WithTimeout(ctx, timeout)
	tStart := time.Now()
	for {
		if ctx.Err() != nil {
			return tracederrors.TracedErrorf("Wait until namespace '%s' created failed: %w", namepaceName, ctx.Err())
		}

		exists, err := n.NamespaceByNameExists(ctx, namepaceName)
		if err != nil {
			return err
		}

		if exists {
			break
		} else {
			waitTime := time.Second * 1
			elapsed := time.Since(tStart)
			logging.LogInfoByCtxf(ctx, "Wait another %s until the kubernetes namespace '%s' is created (%s/%s).", waitTime, namepaceName, elapsed, timeout)
			time.Sleep(waitTime)
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait for kubernetes namespace '%s' to be created finished.", namepaceName)

	return nil
}

func (n *NativeKubernetesCluster) CreateConfigMap(ctx context.Context, namespaceName string, configMapName string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdSecret kubernetesinterfaces.ConfigMap, err error) {
	namespace, err := n.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	return namespace.CreateConfigMap(ctx, configMapName, options)
}

func (n *NativeKubernetesCluster) CreateSecret(ctx context.Context, namespaceName string, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret kubernetesinterfaces.Secret, err error) {
	namespace, err := n.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateSecret(ctx, secretName, options)
}

func (n *NativeKubernetesCluster) SecretByNameExists(ctx context.Context, namespaceName string, secretName string) (exists bool, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.SecretByNameExists(ctx, secretName)
}

// CheckSecretByNameExists checks if a secret exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeKubernetesCluster) CheckSecretByNameExists(ctx context.Context, namespaceName string, secretName string) error {
	exists, err := n.SecretByNameExists(ctx, namespaceName, secretName)
	if err != nil {
		return err
	}
	if !exists {
		return tracederrors.TracedErrorf("Secret '%s' does not exist in namespace '%s'", secretName, namespaceName)
	}
	return nil
}

// CheckNamespaceByNameExists checks if a namespace exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeKubernetesCluster) CheckNamespaceByNameExists(ctx context.Context, namespaceName string) error {
	exists, err := n.NamespaceByNameExists(ctx, namespaceName)
	if err != nil {
		return err
	}
	if !exists {
		return tracederrors.TracedErrorf("Namespace '%s' does not exist", namespaceName)
	}
	return nil
}

// CheckPodByNameExists checks if a pod exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeKubernetesCluster) CheckPodByNameExists(ctx context.Context, namespaceName string, podName string) error {
	exists, err := n.PodByNameExists(ctx, namespaceName, podName)
	if err != nil {
		return err
	}
	if !exists {
		return tracederrors.TracedErrorf("Pod '%s' does not exist in namespace '%s'", podName, namespaceName)
	}
	return nil
}

// CheckReplicaSetByNameExists checks if a replicaSet exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeKubernetesCluster) CheckReplicaSetByNameExists(ctx context.Context, namespaceName string, replicaSetName string) error {
	exists, err := n.ReplicaSetByNameExists(ctx, namespaceName, replicaSetName)
	if err != nil {
		return err
	}
	if !exists {
		return tracederrors.TracedErrorf("ReplicaSet '%s' does not exist in namespace '%s'", replicaSetName, namespaceName)
	}
	return nil
}

// CheckDeploymentByNameExists checks if a deployment exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeKubernetesCluster) CheckDeploymentByNameExists(ctx context.Context, namespaceName string, deploymentName string) error {
	exists, err := n.DeploymentByNameExists(ctx, namespaceName, deploymentName)
	if err != nil {
		return err
	}
	if !exists {
		return tracederrors.TracedErrorf("Deployment '%s' does not exist in namespace '%s'", deploymentName, namespaceName)
	}
	return nil
}

// CheckCronJobByNameExists checks if a cronJob exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeKubernetesCluster) CheckCronJobByNameExists(ctx context.Context, namespaceName string, cronJobName string) error {
	exists, err := n.CronJobByNameExists(ctx, namespaceName, cronJobName)
	if err != nil {
		return err
	}
	if !exists {
		return tracederrors.TracedErrorf("CronJob '%s' does not exist in namespace '%s'", cronJobName, namespaceName)
	}
	return nil
}

func (n *NativeKubernetesCluster) DeleteSecretByName(ctx context.Context, namespaceName string, secretName string) (err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteSecretByName(ctx, secretName)
}

func (n *NativeKubernetesCluster) ConfigMapByNameExists(ctx context.Context, namespaceName string, configmapName string) (exists bool, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.ConfigMapByNameExists(ctx, configmapName)
}

func (n *NativeKubernetesCluster) DeleteConfigMapByName(ctx context.Context, namespaceName string, configmapName string) (err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteConfigMapByName(ctx, configmapName)
}

func (n *NativeKubernetesCluster) GetDiscoveryClient() (discovery.DiscoveryInterface, error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return clientset.Discovery(), nil
}

// ListKindNames retrieves a sorted list of all available resource kind names
// from the Kubernetes API server. It uses the discovery client to query the
// server's preferred resources across all API groups and versions, and returns
// their kind names in alphabetical order.
//
// Returns:
//   - []string: A sorted slice of unique resource kind names (e.g., "Pod", "Service", "Deployment").
//   - error: An error if the discovery client cannot be obtained or if the API server
//     cannot be queried for its available resources.
func (n *NativeKubernetesCluster) ListKindNames(ctx context.Context) ([]string, error) {
	discoveryClient, err := n.GetDiscoveryClient()
	if err != nil {
		return nil, err
	}

	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	apiKinds := []string{}
	for _, apiObjectList := range apiResourceLists {
		for _, apiObject := range apiObjectList.APIResources {
			apiKinds = append(apiKinds, apiObject.Kind)
		}
	}

	sort.Strings(apiKinds)

	return apiKinds, nil
}

func (n *NativeKubernetesCluster) CheckAccessible(ctx context.Context) error {
	clusterName, err := n.GetName()
	if err != nil {
		return err
	}

	_, err = n.WhoAmI(ctx)
	if err != nil {
		return tracederrors.TracedErrorf("Cluster '%s' is not reachable.", clusterName)
	}

	logging.LogInfoByCtxf(ctx, "Cluster '%s' is reachable.", clusterName)

	return err
}

func (n *NativeKubernetesCluster) GetUserNameByContextName(ctx context.Context, kubeContext string) (string, error) {
	return kubeconfigutils.GetUserNameByContextName(ctx, kubeContext)
}

func (n *NativeKubernetesCluster) WhoAmI(ctx context.Context) (*kubernetesimplementationindependend.UserInfo, error) {
	clusterName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	response, err := clientset.AuthenticationV1().SelfSubjectReviews().Create(ctx, &authenticationv1.SelfSubjectReview{}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	var username = response.Status.UserInfo.Username

	logging.LogInfoByCtxf(ctx, "Whoami: Kube config uses user '%s' to log in to cluster '%s'.", username, clusterName)

	return &kubernetesimplementationindependend.UserInfo{
		Username: username,
	}, nil
}

func (n *NativeKubernetesCluster) WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.WaitForPodsOptions) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.WaitUntilAllPodsInNamespaceAreRunning(ctx, options)
}

func (n *NativeKubernetesCluster) GetNamespaceByYamlString(yaml string) (kubernetesinterfaces.Namespace, error) {

	if yaml == "" {
		return nil, tracederrors.TracedErrorEmptyString("yaml")
	}

	objectYamls, err := kubernetesimplementationindependend.UnmarshalObjectYaml(yaml)
	if err != nil {
		return nil, err
	}

	nObjects := len(objectYamls)
	if nObjects != 1 {
		return nil, tracederrors.TracedErrorf("Exepected one yaml document to get namespace by yaml string but got '%d'.", nObjects)
	}

	return n.GetNamespaceByName(objectYamls[0].Namespace())
}

func (n *NativeKubernetesCluster) CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (kubernetesinterfaces.Object, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	namespace, err := n.GetNamespaceByYamlString(options.YamlString)
	if err != nil {
		return nil, err
	}

	return namespace.CreateObject(ctx, options)
}

func (n *NativeKubernetesCluster) RunCommandInTemporaryPod(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (*commandoutput.CommandOutput, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.RunCommandInTemporaryPod(ctx, clientSet, namespaceName, options)
}

func (n *NativeKubernetesCluster) RunCommand(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (*commandoutput.CommandOutput, error) {
	config, err := n.GetConfig()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.RunCommand(ctx, config, namespaceName, options)
}

func (n *NativeKubernetesCluster) ReadSecret(ctx context.Context, namespaceName string, secretName string) (map[string][]byte, error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ReadSecret(ctx, clientset, namespaceName, secretName)
}

func (n *NativeKubernetesCluster) ListNodeNames(ctx context.Context) ([]string, error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return ListNodeNames(ctx, clientset)
}

func (n *NativeKubernetesCluster) DeleteReplicaSetByNames(ctx context.Context, namespaceName string, replicaSetName string) error {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteReplicaSetByName(ctx, replicaSetName)
}

func (n *NativeKubernetesCluster) DeleteDeploymentByNames(ctx context.Context, namespaceName string, deploymentName string) error {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteDeploymentByName(ctx, deploymentName)
}

func (n *NativeKubernetesCluster) DeletePodByNames(ctx context.Context, namespaceName string, podName string) error {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeletePodByName(ctx, podName)
}

func (n *NativeKubernetesCluster) CreatePod(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (kubernetesinterfaces.Pod, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreatePod(ctx, options)
}

func (n *NativeKubernetesCluster) CreateReplicaSet(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (kubernetesinterfaces.ReplicaSet, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateReplicaSet(ctx, options)
}

func (n *NativeKubernetesCluster) CreateDeployment(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (kubernetesinterfaces.Deployment, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateDeployment(ctx, options)
}

func (n *NativeKubernetesCluster) PodByNameExists(ctx context.Context, namespaceName string, podName string) (bool, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.PodByNameExists(ctx, podName)
}

func (n *NativeKubernetesCluster) ReplicaSetByNameExists(ctx context.Context, namespaceName string, replicaSetName string) (bool, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.ReplicaSetByNameExists(ctx, replicaSetName)
}

func (n *NativeKubernetesCluster) DeploymentByNameExists(ctx context.Context, namespaceName string, deploymentName string) (bool, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.DeploymentByNameExists(ctx, deploymentName)
}

func (n *NativeKubernetesCluster) CronJobByNameExists(ctx context.Context, namespaceName string, cronJobName string) (exists bool, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.CronJobByNameExists(ctx, cronJobName)
}

func (n *NativeKubernetesCluster) CreateCronJob(ctx context.Context, namespaceName string, cronJobName string, schedule string, image string, command []string, labels map[string]string) (createdCronJob kubernetesinterfaces.CronJob, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateCronJob(ctx, cronJobName, schedule, image, command, labels)
}

func (n *NativeKubernetesCluster) DeleteCronJobByName(ctx context.Context, namespaceName string, cronJobName string) (err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteCronJobByName(ctx, cronJobName)
}

func (n *NativeKubernetesCluster) ValidateSSHKeyInSecret(
	ctx context.Context,
	options *kubernetesparameteroptions.ValidateSshKeyInSecretOptions,
) (bool, error) {
	logging.LogInfoByCtxf(ctx, "Validate SSH key in secret '%s' started.", options.SecretName)

	clientSet, err := n.GetClientSet()
	if err != nil {
		return false, err
	}

	success, err := nativekubernetes.ValidateSSHKeyInSecret(ctx, clientSet, options)

	logging.LogInfoByCtxf(ctx, "Validate SSH key in secret '%s' finished.", options.SecretName)

	return success, err
}

func (n *NativeKubernetesCluster) CreateClusterRole(ctx context.Context, createOptions *kubernetesparameteroptions.CreateClusterRoleOptions) (createdClusterRole kubernetesinterfaces.ClusterRole, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}
	roleName, err := createOptions.GetName()
	if err != nil {
		return nil, err
	}
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}
	exists, err := n.ClusterRoleByNameExists(ctx, roleName)
	if err != nil {
		return nil, err
	}
	if exists {
		logging.LogInfoByCtxf(ctx, "ClusterRole '%s' already exists.", roleName)
	} else {
		verbs, err := createOptions.GetVerbs()
		if err != nil {
			return nil, err
		}
		resources, err := createOptions.GetResorces()
		if err != nil {
			return nil, err
		}
		apiGroups, err := createOptions.GetAPIGroups()
		if err != nil {
			return nil, err
		}
		role := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{Name: roleName, Labels: map[string]string{}},
			Rules:      []rbacv1.PolicyRule{{Verbs: verbs, Resources: resources, APIGroups: apiGroups}},
		}
		_, err = clientSet.RbacV1().ClusterRoles().Create(ctx, role, metav1.CreateOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to create ClusterRole '%s': %w", roleName, err)
		}
		logging.LogChangedByCtxf(ctx, "Created ClusterRole '%s'.", roleName)
	}
	return n.GetClusterRoleByName(roleName)
}

func (n *NativeKubernetesCluster) DeleteClusterRoleByName(ctx context.Context, roleName string) (err error) {
	if roleName == "" {
		return tracederrors.TracedErrorEmptyString("roleName")
	}
	clientSet, err := n.GetClientSet()
	if err != nil {
		return err
	}
	exists, err := n.ClusterRoleByNameExists(ctx, roleName)
	if err != nil {
		return err
	}
	if exists {
		err = clientSet.RbacV1().ClusterRoles().Delete(ctx, roleName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("failed to delete ClusterRole '%s': %w", roleName, err)
		}
		logging.LogChangedByCtxf(ctx, "ClusterRole '%s' deleted.", roleName)
	} else {
		logging.LogChangedByCtxf(ctx, "ClusterRole '%s' already absent.", roleName)
	}
	return nil
}

func (n *NativeKubernetesCluster) GetClusterRoleByName(roleName string) (kubernetesinterfaces.ClusterRole, error) {
	if roleName == "" {
		return nil, tracederrors.TracedErrorEmptyString("roleName")
	}
	ret := &NativeClusterRole{}

	err := ret.SetName(roleName)
	if err != nil {
		return nil, err
	}

	err = ret.SetKubernetesCluster(n)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (n *NativeKubernetesCluster) ClusterRoleByNameExists(ctx context.Context, roleName string) (exists bool, err error) {
	if roleName == "" {
		return false, tracederrors.TracedErrorEmptyString("roleName")
	}
	clientSet, err := n.GetClientSet()
	if err != nil {
		return false, err
	}
	roleList, err := clientSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to list ClusterRoles: %w", err)
	}
	for _, role := range roleList.Items {
		if role.Name == roleName {
			logging.LogInfoByCtxf(ctx, "ClusterRole '%s' exists.", roleName)
			return true, nil
		}
	}
	logging.LogInfoByCtxf(ctx, "ClusterRole '%s' does not exist.", roleName)
	return false, nil
}
func (n *NativeKubernetesCluster) CreateClusterRoleBinding(ctx context.Context, createOptions *kubernetesparameteroptions.CreateClusterRoleBindingOptions) (createdClusterRoleBinding kubernetesinterfaces.ClusterRoleBinding, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}
	bindingName, err := createOptions.GetName()
	if err != nil {
		return nil, err
	}
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}
	exists, err := n.ClusterRoleBindingByNameExists(ctx, bindingName)
	if err != nil {
		return nil, err
	}
	if exists {
		logging.LogInfoByCtxf(ctx, "ClusterRoleBinding '%s' already exists.", bindingName)
	} else {
		roleRef, err := createOptions.GetRoleRef()
		if err != nil {
			return nil, err
		}
		subjects, err := createOptions.GetSubjects()
		if err != nil {
			return nil, err
		}
		subjectKind, err := createOptions.GetSubjectKind()
		if err != nil {
			return nil, err
		}
		subjectNamespace, err := createOptions.GetSubjectNamespace()
		if err != nil {
			return nil, err
		}
		subjectList := make([]rbacv1.Subject, len(subjects))
		for i, subject := range subjects {
			subjectList[i] = rbacv1.Subject{
				Kind:      subjectKind,
				Name:      subject,
				Namespace: subjectNamespace,
			}
		}
		binding := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: bindingName},
			RoleRef:    rbacv1.RoleRef{Kind: "ClusterRole", Name: roleRef, APIGroup: "rbac.authorization.k8s.io"},
			Subjects:   subjectList,
		}
		_, err = clientSet.RbacV1().ClusterRoleBindings().Create(ctx, binding, metav1.CreateOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to create ClusterRoleBinding '%s': %w", bindingName, err)
		}
		logging.LogChangedByCtxf(ctx, "Created ClusterRoleBinding '%s'.", bindingName)
	}
	return n.GetClusterRoleBindingByName(bindingName)
}

func (n *NativeKubernetesCluster) DeleteClusterRoleBindingByName(ctx context.Context, bindingName string) (err error) {
	if bindingName == "" {
		return tracederrors.TracedErrorEmptyString("bindingName")
	}
	clientSet, err := n.GetClientSet()
	if err != nil {
		return err
	}
	exists, err := n.ClusterRoleBindingByNameExists(ctx, bindingName)
	if err != nil {
		return err
	}
	if exists {
		err = clientSet.RbacV1().ClusterRoleBindings().Delete(ctx, bindingName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("failed to delete ClusterRoleBinding '%s': %w", bindingName, err)
		}
		logging.LogChangedByCtxf(ctx, "ClusterRoleBinding '%s' deleted.", bindingName)
	} else {
		logging.LogChangedByCtxf(ctx, "ClusterRoleBinding '%s' already absent.", bindingName)
	}
	return nil
}

func (n *NativeKubernetesCluster) GetClusterRoleBindingByName(bindingName string) (kubernetesinterfaces.ClusterRoleBinding, error) {
	if bindingName == "" {
		return nil, tracederrors.TracedErrorEmptyString("bindingName")
	}
	ret := &NativeClusterRoleBinding{}

	err := ret.SetName(bindingName)
	if err != nil {
		return nil, err
	}

	err = ret.SetKubernetesCluster(n)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (n *NativeKubernetesCluster) ClusterRoleBindingByNameExists(ctx context.Context, bindingName string) (exists bool, err error) {
	if bindingName == "" {
		return false, tracederrors.TracedErrorEmptyString("bindingName")
	}
	clientSet, err := n.GetClientSet()
	if err != nil {
		return false, err
	}
	bindingList, err := clientSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to list ClusterRoleBindings: %w", err)
	}
	for _, binding := range bindingList.Items {
		if binding.Name == bindingName {
			logging.LogInfoByCtxf(ctx, "ClusterRoleBinding '%s' exists.", bindingName)
			return true, nil
		}
	}
	logging.LogInfoByCtxf(ctx, "ClusterRoleBinding '%s' does not exist.", bindingName)
	return false, nil
}

func (n *NativeKubernetesCluster) ListClusterRoleNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}
	roleList, err := clientSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list ClusterRoles: %w", err)
	}
	names := make([]string, len(roleList.Items))
	for i, role := range roleList.Items {
		names[i] = role.Name
	}
	return names, nil
}

func (n *NativeKubernetesCluster) ListClusterRoleBindingNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}
	bindingList, err := clientSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list ClusterRoleBindings: %w", err)
	}
	names := make([]string, len(bindingList.Items))
	for i, binding := range bindingList.Items {
		names[i] = binding.Name
	}
	return names, nil
}

func (n *NativeKubernetesCluster) CreateRole(ctx context.Context, namespaceName string, createOptions *kubernetesparameteroptions.CreateRoleOptions) (createdRole kubernetesinterfaces.Role, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}
	return namespace.CreateRole(ctx, createOptions)
}

func (n *NativeKubernetesCluster) DeleteRoleByName(ctx context.Context, namespaceName string, roleName string) (err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}
	return namespace.DeleteRoleByName(ctx, roleName)
}

func (n *NativeKubernetesCluster) GetRoleByName(namespaceName string, roleName string) (kubernetesinterfaces.Role, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}
	return namespace.GetRoleByName(roleName)
}

func (n *NativeKubernetesCluster) RoleByNameExists(ctx context.Context, namespaceName string, roleName string) (exists bool, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}
	return namespace.RoleByNameExists(ctx, roleName)
}

func (n *NativeKubernetesCluster) CreateRoleBinding(ctx context.Context, namespaceName string, createOptions *kubernetesparameteroptions.CreateRoleBindingOptions) (createdRoleBinding kubernetesinterfaces.RoleBinding, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}
	return namespace.CreateRoleBinding(ctx, createOptions)
}

func (n *NativeKubernetesCluster) DeleteRoleBindingByName(ctx context.Context, namespaceName string, roleBindingName string) (err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}
	return namespace.DeleteRoleBindingByName(ctx, roleBindingName)
}

func (n *NativeKubernetesCluster) GetRoleBindingByName(namespaceName string, roleBindingName string) (kubernetesinterfaces.RoleBinding, error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}
	return namespace.GetRoleBindingByName(roleBindingName)
}

func (n *NativeKubernetesCluster) RoleBindingByNameExists(ctx context.Context, namespaceName string, roleBindingName string) (exists bool, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}
	return namespace.RoleBindingByNameExists(ctx, roleBindingName)
}
