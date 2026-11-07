package nativekubernetes

import (
	"context"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// pvcHelperImage contains find and sh to list the files inside the mounted PVC.
const pvcHelperImage = "busybox:stable"

// pvcHelperMountPath is the path inside the temporary pod where the PVC is mounted.
const pvcHelperMountPath = "/mnt/pvc"

// ListPersistentVolumeClaimFileNames returns the names of all files and directories
// inside the persistent volume claim 'pvcName' in namespace 'namespaceName',
// relative to the PVC root. Hidden files are included and the result is sorted
// alphabetically.
//
// The listing is done in a temporary pod with the PVC mounted. The pod is
// deleted afterwards.
func ListPersistentVolumeClaimFileNames(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, pvcName string) (fileNames []string, err error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if pvcName == "" {
		return nil, tracederrors.TracedErrorEmptyString("pvcName")
	}

	logging.LogInfoByCtxf(ctx, "List files in persistent volume claim '%s' in namespace '%s' started.", pvcName, namespaceName)

	// Fail early with a clear error instead of starting a pod for a missing PVC.
	_, err = clientset.CoreV1().PersistentVolumeClaims(namespaceName).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, tracederrors.TracedErrorf("Persistent volume claim '%s' does not exist in namespace '%s'. Unable to list its files.", pvcName, namespaceName)
		}
		return nil, tracederrors.TracedErrorf("Failed to get persistent volume claim '%s' in namespace '%s': %w", pvcName, namespaceName, err)
	}

	// busybox find does not know -printf, so the './' prefix is stripped in Go.
	// The mount path is passed as positional parameter to avoid any shell quoting problems.
	// $0 = "sh", $1 = mount point of the PVC.
	podName := "list-files-in-pvc-" + uuidutils.Generate(ctx)
	output, err := RunCommandInTemporaryPod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    pvcHelperImage,
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			PersistentVolumeClaims: map[string]kubernetesparameteroptions.PersistentVolumeClaimSource{
				pvcHelperMountPath: {PVCName: pvcName},
			},
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"/bin/sh", "-c", `cd "$1" && find . -mindepth 1 -print`, "sh", pvcHelperMountPath},
			},
		},
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Listing files in persistent volume claim '%s' in namespace '%s' failed: %w", pvcName, namespaceName, err)
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read the stdout of the temporary pod for persistent volume claim '%s' in namespace '%s': %w", pvcName, namespaceName, err)
	}

	ret := []string{}
	for _, line := range strings.Split(stdout, "\n") {
		fileName := strings.TrimPrefix(line, "./")
		if fileName == "" {
			continue
		}

		ret = append(ret, fileName)
	}

	sort.Strings(ret)

	logging.LogInfoByCtxf(ctx, "List files in persistent volume claim '%s' in namespace '%s' finished. Found '%d' files.", pvcName, namespaceName, len(ret))

	return ret, nil
}
