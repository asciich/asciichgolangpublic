package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeletePod(ctx context.Context, clientset *kubernetes.Clientset, podName string, namespace string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if podName == "" {
		return tracederrors.TracedErrorEmptyString("podName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Delete pod '%s' in namepace '%s' started.", podName, namespace)

	deletePolicy := metav1.DeletePropagationBackground
	err := clientset.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err == nil {
		logging.LogChangedByCtxf(ctx, "Pod '%s' in namepsace '%s' deleted.", podName, namespace)
	} else {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Pod '%s' already absent in namespace '%s'.", podName, namespace)
		} else {
			return tracederrors.TracedErrorf("Failed to delete pod '%s' in namespace '%s': %w", podName, namespace, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Delete pod '%s' in namepace '%s' finished.", podName, namespace)

	return nil
}
