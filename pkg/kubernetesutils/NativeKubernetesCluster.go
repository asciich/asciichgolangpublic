package kubernetesutils

import (
	"context"
	"path/filepath"
	"sort"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type NativeKubernetesCluster struct {
	clientset *kubernetes.Clientset
}

// Get a client set based on the ~/.kube/config
func GetClientSetFromKubeconfig(ctx context.Context) (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, tracederrors.TracedError("Unable to find home directory for kubeconfig")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error building kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error creating Kubernetes clientset: %w", err)
	}

	logging.LogInfoByCtx(ctx, "Created kubernetes clientset from ~/.kube/config")

	return clientset, nil
}

func GetInClusterClientSet(ctx context.Context) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error getting in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error creating clientset: %w", err)
	}

	logging.LogInfoByCtx(ctx, "Created kubernetes in cluster clientset")

	return clientset, nil
}

// Get the kubernetes.Clientset to communicate with the kubernetes cluster.
//
// If in cluster authentication is available (e.g. running in a pod in the cluster) the returned clientset uses this method.
//
// Otherwise a clientset based on ~/.kube/config is returned.
func GetClientSet(ctx context.Context) (*kubernetes.Clientset, error) {
	if IsInClusterAuthenticationAvailable(ctx) {
		return GetInClusterClientSet(ctx)
	}

	return GetClientSetFromKubeconfig(ctx)
}

func GetNativeKubernetesClusterByName(ctx context.Context, clusterName string) (*NativeKubernetesCluster, error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	clientSet, err := GetClientSet(ctx)
	if err != nil {
		return nil, err
	}

	return &NativeKubernetesCluster{
		clientset: clientSet,
	}, nil
}

func (n *NativeKubernetesCluster) CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace Namespace, err error) {
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

func (n *NativeKubernetesCluster) GetClientSet() (*kubernetes.Clientset, error) {
	if n.clientset == nil {
		return nil, tracederrors.TracedError("Clientset not set")
	}

	return n.clientset, nil
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
	return "", tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) GetName() (name string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) GetNamespaceByName(name string) (namespace Namespace, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return &NativeNamespace{
		name:              name,
		kubernetesCluster: n,
	}, nil
}
func (n *NativeKubernetesCluster) GetResourceByNames(resourceName string, resourceType string, namespaceName string) (resource Resource, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) ListNamespaces(ctx context.Context) (namespaces []Namespace, err error) {
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
func (n *NativeKubernetesCluster) ListResources(options *parameteroptions.ListKubernetesResourcesOptions) (resources []Resource, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeKubernetesCluster) ListResourceNames(options *parameteroptions.ListKubernetesResourcesOptions) (resourceNames []string, err error) {
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

func (n *NativeKubernetesCluster) CreateSecret(ctx context.Context, namespaceName string, secretName string, options *CreateSecretOptions) (createdSecret Secret, err error) {
	namespace, err := n.GetNamespaceByName(namespaceName)
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
