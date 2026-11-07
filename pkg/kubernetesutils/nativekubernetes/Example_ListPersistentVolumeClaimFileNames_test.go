package nativekubernetes_test

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This example shows how to list all files inside a persistent volume claim.
// A temporary pod with the PVC mounted is used to do so. Hidden files are
// included and the returned file names are sorted alphabetically.
func Test_Example_ListPersistentVolumeClaimFileNames(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+kindutils.SharedClusterName)
	require.NoError(t, err)

	const namespaceName = "example-list-files-in-pvc"
	const pvcName = "example-pvc"

	// Ensure the namespace exists:
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	// Ensure the PVC exists:
	_, err = clientset.CoreV1().PersistentVolumeClaims(namespaceName).Create(ctx, &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: pvcName},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("64Mi"),
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		// Leftover from a previous test run. It can be reused since it is bound to a volume.
		logging.LogInfoByCtxf(ctx, "Reason why its ok to ignore this: PVC '%s' already exists: %v", pvcName, err)
	}
	// ... prepare test environment finished.
	// -----

	// Write some files into the PVC using a temporary pod with the PVC mounted:
	_, err = nativekubernetes.RunCommandInTemporaryPod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "busybox:stable",
			PodName:                  "example-write-into-pvc",
			DeleteAlreadyExistingPod: true,
			PersistentVolumeClaims: map[string]kubernetesparameteroptions.PersistentVolumeClaimSource{
				"/mnt/pvc": {PVCName: pvcName},
			},
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"/bin/sh", "-c", `cd "$1" && mkdir -p subdir && echo 'Hello from pvc' > testfile.txt && echo 'Hello from subdir' > subdir/testfile.txt && echo 'hidden' > .hidden`, "sh", "/mnt/pvc"},
			},
		},
	)
	require.NoError(t, err)

	// List all files inside the PVC:
	fileNames, err := nativekubernetes.ListPersistentVolumeClaimFileNames(ctx, clientset, namespaceName, pvcName)
	require.NoError(t, err)

	// The output has to be sorted to be deterministic:
	require.True(t, sort.StringsAreSorted(fileNames))

	// Hidden files are included:
	require.EqualValues(
		t,
		[]string{".hidden", "subdir", "subdir/testfile.txt", "testfile.txt"},
		fileNames,
	)
}
