package nativekubernetes

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespace(ctx context.Context, clientSet *kubernetes.Clientset, namespaceName string) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	logging.LogInfoByCtxf(ctx, "Create kubernetes namespace '%s' started.", namespaceName)

	namespaceSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	var created = true
	_, err := clientSet.CoreV1().Namespaces().Create(ctx, namespaceSpec, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			logging.LogInfoByCtxf(ctx, "Kubernetes namespace '%s' already exists. Skip creation.", namespaceName)
			created = false
		} else {
			return tracederrors.TracedErrorf("Failed to create kubernetes namespace '%s': %w", namespaceName, err)
		}
	}

	if created {
		logging.LogChangedByCtxf(ctx, "Created kubernetes namespace '%s'.", namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Create kubernetes namespace '%s' finished.", namespaceName)

	return nil
}

func NamespaceExists(ctx context.Context, clientSet *kubernetes.Clientset, namespaceName string) (bool, error) {
	if clientSet == nil {
		return false, tracederrors.TracedErrorNil("clientSet")
	}

	if namespaceName == "" {
		return false, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	var exists = true
	_, err := clientSet.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			exists = false
		} else {
			return false, tracederrors.TracedErrorf("Failed to check if namespace '%s' exists: %w", namespaceName, err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Kubernetes namespace '%s' exists.", namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Kubernetes namespace '%s' does not exist.", namespaceName)
	}

	return exists, nil
}

func DeleteNamespace(ctx context.Context, clientSet *kubernetes.Clientset, namespaceName string) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	logging.LogInfoByCtxf(ctx, "Delete kubernetes namespace '%s' started.", namespaceName)

	deletePolicy := metav1.DeletePropagationForeground

	err := clientSet.CoreV1().Namespaces().Delete(
		ctx,
		namespaceName,
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)

	var deleted = true
	if err != nil {
		if errors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Kubernetes namespace '%s' is already absent. Skip delete.", namespaceName)
			deleted = false
		} else {
			return tracederrors.TracedErrorf("Failed to delete kubernetes namespace '%s': %w", namespaceName, err)
		}
	}

	err = WaitForNamespaceDeletion(ctx, clientSet, namespaceName)
	if err != nil {
		return err
	}

	if deleted {
		logging.LogChangedByCtxf(ctx, "Deleted kubernetes namespace '%s'.", namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete kubernetes namespace '%s' finished.", namespaceName)

	return nil
}

func WaitForNamespaceDeletion(ctx context.Context, clientSet *kubernetes.Clientset, namespaceName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	return wait.PollUntilContextCancel(ctx, 2*time.Second, true, func(ctx context.Context) (bool, error) {
		_, err := clientSet.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})

		if err != nil {
			if errors.IsNotFound(err) {
				logging.LogInfoByCtxf(ctx, "Wait for namespace '%s' deleted: Namespace successfully deleted.", namespaceName)
				return true, nil // Done waiting
			}
			return false, err // An actual error occurred
		}

		logging.LogInfoByCtxf(ctx, "Wait for namespace '%s' deleted: Namespace still exists, waiting...", namespaceName)
		return false, nil // Keep polling
	})
}
