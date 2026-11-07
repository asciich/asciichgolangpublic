package nativekubernetesoo

import (
	"context"
	"sort"

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
