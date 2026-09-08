package nativekubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeleteDaemonSet(ctx context.Context, clientset *kubernetes.Clientset, daemonSetName string, namespace string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if daemonSetName == "" {
		return tracederrors.TracedErrorEmptyString("daemonSetName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Delete DaemonSet '%s' in namespace '%s' started.", daemonSetName, namespace)

	deletePolicy := metav1.DeletePropagationBackground
	err := clientset.AppsV1().DaemonSets(namespace).Delete(ctx, daemonSetName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err == nil {
		err := WaitForDaemonSetDeleted(ctx, clientset, namespace, daemonSetName, time.Second*30)
		if err != nil {
			return err
		}
		logging.LogChangedByCtxf(ctx, "DaemonSet '%s' in namespace '%s' deleted.", daemonSetName, namespace)
	} else {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "DaemonSet '%s' already absent in namespace '%s'.", daemonSetName, namespace)
		} else {
			return tracederrors.TracedErrorf("Failed to delete DaemonSet '%s' in namespace '%s': %w", daemonSetName, namespace, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Delete DaemonSet '%s' in namespace '%s' finished.", daemonSetName, namespace)

	return nil
}

func CreateDaemonSet(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, options *kubernetesparameteroptions.KubernetesRunCommandOptions) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	daemonSetName, err := options.GetDaemonSetName()
	if err != nil {
		return err
	}

	imageName, err := options.GetImageName()
	if err != nil {
		return err
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return err
	}

	command, err := options.GetCommand()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create DaemonSet '%s' in namespace '%s' using container image '%s' started.", daemonSetName, namespaceName, imageName)

	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: daemonSetName,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": daemonSetName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": daemonSetName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    containerName,
							Image:   imageName,
							Command: command,
							Stdin:   true,
							TTY:     true,
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	logging.LogInfoByCtxf(ctx, "Going to start DaemonSet '%s' in namespace '%s' using container image '%s'.", daemonSetName, namespaceName, imageName)
	_, err = clientset.AppsV1().DaemonSets(namespaceName).Create(ctx, daemonSet, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) && options.DeleteAlreadyExistingDaemonSet {
			logging.LogInfoByCtxf(ctx, "Going to delete already existing DaemonSet '%s' in namespace '%s' before running command.", daemonSetName, namespaceName)
			err = DeleteDaemonSet(ctx, clientset, daemonSetName, namespaceName)
			if err != nil {
				return err
			}
			_, err = clientset.AppsV1().DaemonSets(namespaceName).Create(ctx, daemonSet, metav1.CreateOptions{})
			if err != nil {
				return tracederrors.TracedErrorf("Error creating DaemonSet: %w", err)
			}
		} else {
			return tracederrors.TracedErrorf("Error creating DaemonSet: %w", err)
		}
	}

	if options.WaitForDaemonSetAvailable {
		err = WaitForDaemonSetAvailable(ctx, clientset, namespaceName, daemonSetName, time.Minute*2)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Create DaemonSet '%s' in namespace '%s' using container image '%s' finished.", daemonSetName, namespaceName, imageName)

	return nil
}

func DaemonSetExists(ctx context.Context, clientSet *kubernetes.Clientset, daemonSetName string, namespace string) (bool, error) {
	if clientSet == nil {
		return false, tracederrors.TracedErrorNil("clientSet")
	}

	if daemonSetName == "" {
		return false, tracederrors.TracedErrorEmptyString("daemonSetName")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Check if DaemonSet '%s' in namespace '%s' exists.", daemonSetName, namespace)

	_, err := clientSet.AppsV1().DaemonSets(namespace).Get(ctx, daemonSetName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' does not exist.", daemonSetName, namespace)
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Failed to get DaemonSet '%s' in namespace '%s' to check if exists: %w", daemonSetName, namespace, err)
	}

	logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' exists.", daemonSetName, namespace)
	return true, nil
}

func WaitForDaemonSetDeleted(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, daemonSetName string, timeout time.Duration) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if daemonSetName == "" {
		return tracederrors.TracedErrorEmptyString("daemonSetName")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for DaemonSet '%s' in namespace '%s' to be deleted started.", daemonSetName, namespaceName)

	_, err := clientset.AppsV1().DaemonSets(namespaceName).Get(ctx, daemonSetName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' is already deleted.", daemonSetName, namespaceName)
		return nil
	}

	w, err := clientset.AppsV1().DaemonSets(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", daemonSetName),
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' is already deleted.", daemonSetName, namespaceName)
		} else {
			return fmt.Errorf("failed to set up watch for DaemonSet %s: %w", daemonSetName, err)
		}
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case _, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for DaemonSet '%s' in namespace '%s' to be deleted", daemonSetName, namespaceName)
			}

			_, err := clientset.AppsV1().DaemonSets(namespaceName).Get(ctx, daemonSetName, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' is now deleted.", daemonSetName, namespaceName)
				return nil
			}

			logging.LogInfoByCtxf(ctx, "Still waiting for DaemonSet '%s' in namespace '%s' to be deleted.", daemonSetName, namespaceName)
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for DaemonSet '%s' in namespace '%s' to be deleted", daemonSetName, namespaceName)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}

func WaitForDaemonSetAvailable(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, daemonSetName string, timeout time.Duration) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if daemonSetName == "" {
		return tracederrors.TracedErrorEmptyString("daemonSetName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for DaemonSet '%s' in namespace '%s' to be available started.", daemonSetName, namespace)

	w, err := clientSet.AppsV1().DaemonSets(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", daemonSetName),
	})
	if err != nil {
		return fmt.Errorf("failed to set up watch for DaemonSet %s: %w", daemonSetName, err)
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for DaemonSet '%s' in namespace '%s' to be available", daemonSetName, namespace)
			}

			daemonSet, ok := event.Object.(*appsv1.DaemonSet)
			if !ok {
				continue
			}

			// Check if the desired number of pods are scheduled and available on all eligible nodes
			if daemonSet.Status.DesiredNumberScheduled > 0 &&
				daemonSet.Status.NumberAvailable == daemonSet.Status.DesiredNumberScheduled {
				logging.LogInfoByCtxf(ctx, "Wait for DaemonSet '%s' in namespace '%s' to be available finished. All desired pods are available.", daemonSetName, namespace)
				return nil
			}
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for DaemonSet '%s' in namespace '%s' to be available", daemonSetName, namespace)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}

func ListDaemonSets(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	daemonSetList, err := clientset.AppsV1().DaemonSets(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list DaemonSets in namespace '%s'.", namespaceName)
	}

	daemonSetNames := []string{}
	for _, ds := range daemonSetList.Items {
		daemonSetNames = append(daemonSetNames, ds.Name)
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' DaemonSets in namespace '%s'.", len(daemonSetNames), namespaceName)

	return daemonSetNames, nil
}

func ListDaemonSetNames(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	return ListDaemonSets(ctx, clientset, namespaceName)
}
