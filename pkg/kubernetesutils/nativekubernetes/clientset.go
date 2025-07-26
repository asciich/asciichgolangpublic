package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeconfigutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func GetInClusterConfig(ctx context.Context) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error getting in-cluster config: %w", err)
	}

	return config, nil
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
