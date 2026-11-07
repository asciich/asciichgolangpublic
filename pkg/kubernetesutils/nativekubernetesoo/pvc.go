package nativekubernetesoo

import (
	"context"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListPersistentVolumeClaimNames returns the sorted names of all persistent volume claims of the cluster.
//
// If namespaceName is empty the persistent volume claims of all namespaces are listed.
// Otherwise only the persistent volume claims of the given namespace are returned.
func (n *NativeKubernetesCluster) ListPersistentVolumeClaimNames(ctx context.Context, namespaceName string) (persistentVolumeClaimNames []string, err error) {
	if namespaceName == "" {
		logging.LogInfoByCtxf(ctx, "List persistent volume claims of all namespaces started.")
	} else {
		logging.LogInfoByCtxf(ctx, "List persistent volume claims of namespace '%s' started.", namespaceName)
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	claimList, err := clientSet.CoreV1().PersistentVolumeClaims(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list persistent volume claims: %w", err)
	}

	persistentVolumeClaimNames = []string{}
	for _, claim := range claimList.Items {
		persistentVolumeClaimNames = append(persistentVolumeClaimNames, claim.Name)
	}

	sort.Strings(persistentVolumeClaimNames)

	if namespaceName == "" {
		logging.LogInfoByCtxf(ctx, "List persistent volume claims of all namespaces finished. Found '%d' persistent volume claims.", len(persistentVolumeClaimNames))
	} else {
		logging.LogInfoByCtxf(ctx, "List persistent volume claims of namespace '%s' finished. Found '%d' persistent volume claims.", namespaceName, len(persistentVolumeClaimNames))
	}

	return persistentVolumeClaimNames, nil
}

// ListPersistentVolumeClaimFileNames returns the names of all files and directories
// inside the persistent volume claim 'pvcName' in namespace 'namespaceName',
// relative to the PVC root. Hidden files are included and the result is sorted
// alphabetically.
//
// The listing is done in a temporary pod with the PVC mounted. The pod is
// deleted afterwards.
func (n *NativeKubernetesCluster) ListPersistentVolumeClaimFileNames(ctx context.Context, namespaceName string, pvcName string) (fileNames []string, err error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListPersistentVolumeClaimFileNames(ctx, clientSet, namespaceName, pvcName)
}
