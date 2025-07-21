package nativekubernetes

import (
	"context"
	"sort"
	"time"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeconfigutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"

	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type NativeKubernetesCluster struct {
	name   string
	config *rest.Config

	// client caches:
	clientSetCache     *kubernetes.Clientset
	dynamicClientCache *dynamic.DynamicClient
}

func GetConfigFromKubeconfig(ctx context.Context, clusterName string) (*rest.Config, error) {
	kubeconfig, err := kubeconfigutils.GetDefaultKubeConfigPath(ctx)
	if err != nil {
		return nil, err
	}

	var config *rest.Config
	if clusterName == "" {
		logging.LogInfoByCtx(ctx, "clusterName not set. Loading config for default kubernetes cluster. If this fails the missing default cluster/ context could be the root cause.")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Error building kubeconfig: %w", err)
		}

		logging.LogInfoByCtx(ctx, "clusterName not set. Loaded config for default kubernetes cluster.")
	} else {
		kubeContext, err := kubeconfigutils.GetContextNameByClusterName(ctx, clusterName)
		if err != nil {
			return nil, err
		}

		kubeContextPath, err := kubeconfigutils.GetKubeConfigPath(ctx)
		if err != nil {
			return nil, err
		}

		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeContextPath},
			&clientcmd.ConfigOverrides{
				CurrentContext: kubeContext,
			},
		).ClientConfig()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Error building kubeconfig for cluster '%s': %w", clusterName, err)
		}

		logging.LogInfoByCtxf(ctx, "Loaded config for cluster '%s' kubernetes cluster.", clusterName)
	}

	return config, nil
}

func GetInClusterConfig(ctx context.Context) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error getting in-cluster config: %w", err)
	}

	return config, nil
}

// Get the rest.Config to communicate with the kubernetes cluster.
//
// If in cluster authentication is available (e.g. running in a pod in the cluster) the returned config uses this method.
//
// Otherwise a config based on ~/.kube/config is returned.
func GetConfig(ctx context.Context, clusterName string) (*rest.Config, error) {
	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		return GetInClusterConfig(ctx)
	}

	return GetConfigFromKubeconfig(ctx, clusterName)
}

// Get the kubernetes.Clientset to communicate with the kubernetes cluster.
//
// If in cluster authentication is available (e.g. running in a pod in the cluster) the returned clientset uses this method.
//
// Otherwise a clientset based on ~/.kube/config is returned.
func GetClientSet(ctx context.Context, clusterName string) (*kubernetes.Clientset, error) {
	config, err := GetConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create kubernetes clientset: %w", err)
	}

	return clientset, nil
}

func GetClusterByName(ctx context.Context, clusterName string) (*NativeKubernetesCluster, error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	config, err := GetConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	return &NativeKubernetesCluster{
		name:   clusterName,
		config: config,
	}, nil
}

func GetDefaultCluster(ctx context.Context) (*NativeKubernetesCluster, error) {
	config, err := GetConfig(ctx, "")
	if err != nil {
		return nil, err
	}

	return &NativeKubernetesCluster{
		config: config,
	}, nil
}

func (n *NativeKubernetesCluster) CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace kubernetesinterfaces.Namespace, err error) {
	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	exists, err := n.NamespaceByNameExists(ctx, namespaceName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already exists. Skip creation.", namespaceName)
	} else {
		clientset, err := n.GetClientSet()
		if err != nil {
			return nil, err
		}

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   namespaceName,
				Labels: map[string]string{},
			},
		}

		_, err = clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Error creating namespace '%s': %v", namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' created successfully.\n", namespaceName)

		err = n.WaitUntilNamespaceCreated(ctx, namespaceName)
		if err != nil {
			return nil, err
		}
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

		err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, deleteOptions)
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

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
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
			return tracederrors.TracedErrorf("Wait until namespace '%s' deleted failed: %w", namepaceName, err)
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
			return tracederrors.TracedErrorf("Wait until namespace '%s' created failed: %w", namepaceName, err)
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

	kubeContext, err := n.GetKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	username, err := kubeconfigutils.GetUserNameByContextName(ctx, kubeContext)
	if err != nil {
		return nil, err
	}

	response, err := clientset.AuthenticationV1().SelfSubjectReviews().Create(ctx, &authenticationv1.SelfSubjectReview{}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	username = response.Status.UserInfo.Username

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

func (n *NativeKubernetesCluster) RunCommandInTemporaryPod(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (*commandexecutor.CommandOutput, error) {
	config, err := n.GetConfig()
	if err != nil {
		return nil, err
	}

	return RunCommandInTemporaryPod(ctx, config, options)
}

func (n *NativeKubernetesCluster) ReadSecret(ctx context.Context, namespaceName string, secretName string) (map[string][]byte, error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return ReadSecret(ctx, clientset, namespaceName, secretName)
}

func (n *NativeKubernetesCluster) ListNodeNames(ctx context.Context) ([]string, error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return ListNodeNames(ctx, clientset)
}
