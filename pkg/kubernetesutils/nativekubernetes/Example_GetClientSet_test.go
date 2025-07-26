package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This example shows how to get a clientset (the native one of the K8s client-go library).
func Test_Example_GetClientSet(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Get the client set (by not specifying the cluster name the default one is returned.)
	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-kubernetesutils")
	require.NoError(t, err)

	// As an example we use the clientset to list the namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	// Load all namespace names into a string slice so it's easier to check if expected namespaces are present:
	namespaceNames := []string{}
	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}

	// We expect the default namespace:
	require.Contains(t, namespaceNames, "default")

	// Also the kube-system is expected:
	require.Contains(t, namespaceNames, "kube-system")
}
