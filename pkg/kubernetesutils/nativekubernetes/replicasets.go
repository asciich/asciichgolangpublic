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
	"k8s.io/client-go/rest"
)

func DeleteReplicaSet(ctx context.Context, clientset *kubernetes.Clientset, replicaSetName string, namespace string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if replicaSetName == "" {
		return tracederrors.TracedErrorEmptyString("replicaSetName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Delete ReplicaSet '%s' in namespace '%s' started.", replicaSetName, namespace)

	deletePolicy := metav1.DeletePropagationBackground
	err := clientset.AppsV1().ReplicaSets(namespace).Delete(ctx, replicaSetName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err == nil {
		err := WaitForReplicaSetDeleted(ctx, clientset, namespace, replicaSetName, time.Second*30)
		if err != nil {
			return err
		}
		logging.LogChangedByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' deleted.", replicaSetName, namespace)
	} else {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' already absent in namespace '%s'.", replicaSetName, namespace)
		} else {
			return tracederrors.TracedErrorf("Failed to delete ReplicaSet '%s' in namespace '%s': %w", replicaSetName, namespace, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Delete ReplicaSet '%s' in namespace '%s' finished.", replicaSetName, namespace)

	return nil
}

func CreateReplicaSet(ctx context.Context, config *rest.Config, options *kubernetesparameteroptions.RunCommandOptions) error {
	if config == nil {
		return tracederrors.TracedErrorNil("config")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	namespace, err := options.GetNamespaceName()
	if err != nil {
		return err
	}

	replicaSetName, err := options.GetReplicaSetName()
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

	replicas := options.GetReplicas()

	logging.LogInfoByCtxf(ctx, "Create ReplicaSet '%s' in namespace '%s' using container image '%s' started.", replicaSetName, namespace, imageName)

	clientset, err := GetClientSetFromRestConfig(ctx, config)
	if err != nil {
		return err
	}

	replicaSet := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: replicaSetName,
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": replicaSetName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": replicaSetName,
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

	logging.LogInfoByCtxf(ctx, "Going to start ReplicaSet '%s' in namespace '%s' using container image '%s'.", replicaSetName, namespace, imageName)
	_, err = clientset.AppsV1().ReplicaSets(namespace).Create(ctx, replicaSet, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) && options.DeleteAlreadyExistingReplicaSet {
			logging.LogInfoByCtxf(ctx, "Going to delete already existing ReplicaSet '%s' in namespace '%s' before running command.", replicaSetName, namespace)
			err = DeleteReplicaSet(ctx, clientset, replicaSetName, namespace)
			if err != nil {
				return err
			}
			_, err = clientset.AppsV1().ReplicaSets(namespace).Create(ctx, replicaSet, metav1.CreateOptions{})
			if err != nil {
				return tracederrors.TracedErrorf("Error creating ReplicaSet: %w", err)
			}
		} else {
			return tracederrors.TracedErrorf("Error creating ReplicaSet: %w", err)
		}
	}

	if options.WaitForReplicaSetAvailable {
		err = WaitForReplicaSetAvailable(ctx, clientset, namespace, replicaSetName, time.Minute*1)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Create ReplicaSet '%s' in namespace '%s' using container image '%s' finished.", replicaSetName, namespace, imageName)

	return nil
}

func ReplicaSetExists(ctx context.Context, clientSet *kubernetes.Clientset, replicaSetName string, namespace string) (bool, error) {
	if clientSet == nil {
		return false, tracederrors.TracedErrorNil("clientSet")
	}

	if replicaSetName == "" {
		return false, tracederrors.TracedErrorEmptyString("replicaSetName")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Check if ReplicaSet '%s' in namespace '%s' exists.", replicaSetName, namespace)

	_, err := clientSet.AppsV1().ReplicaSets(namespace).Get(ctx, replicaSetName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' does not exist.", replicaSetName, namespace)
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Failed to get ReplicaSet '%s' in namespace '%s' to check if exists: %w", replicaSetName, namespace, err)
	}

	logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' exists.", replicaSetName, namespace)
	return true, nil
}

func WaitForReplicaSetDeleted(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, replicaSetName string, timeout time.Duration) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if replicaSetName == "" {
		return tracederrors.TracedErrorEmptyString("replicaSetName")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for ReplicaSet '%s' in namespace '%s' to be deleted started.", replicaSetName, namespaceName)

	_, err := clientset.AppsV1().ReplicaSets(namespaceName).Get(ctx, replicaSetName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' is already deleted.", replicaSetName, namespaceName)
		return nil
	}

	w, err := clientset.AppsV1().ReplicaSets(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", replicaSetName),
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' is already deleted.", replicaSetName, namespaceName)
		} else {
			return fmt.Errorf("failed to set up watch for ReplicaSet %s: %w", replicaSetName, err)
		}
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case _, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for ReplicaSet '%s' in namespace '%s' to be deleted", replicaSetName, namespaceName)
			}

			_, err := clientset.AppsV1().ReplicaSets(namespaceName).Get(ctx, replicaSetName, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' is now deleted.", replicaSetName, namespaceName)
				return nil
			}

			logging.LogInfoByCtxf(ctx, "Still waiting for ReplicaSet '%s' in namespace '%s' to be deleted.", replicaSetName, namespaceName)
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for ReplicaSet '%s' in namespace '%s' to be deleted", replicaSetName, namespaceName)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}

func WaitForReplicaSetAvailable(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, replicaSetName string, timeout time.Duration) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if replicaSetName == "" {
		return tracederrors.TracedErrorEmptyString("replicaSetName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for ReplicaSet '%s' in namespace '%s' to be available started.", replicaSetName, namespace)

	w, err := clientSet.AppsV1().ReplicaSets(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", replicaSetName),
	})
	if err != nil {
		return fmt.Errorf("failed to set up watch for ReplicaSet %s: %w", replicaSetName, err)
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for ReplicaSet '%s' in namespace '%s' to be available", replicaSetName, namespace)
			}

			replicaSet, ok := event.Object.(*appsv1.ReplicaSet)
			if !ok {
				continue
			}

			// Check if the desired number of replicas are available
			if replicaSet.Status.AvailableReplicas == *replicaSet.Spec.Replicas {
				logging.LogInfoByCtxf(ctx, "Wait for ReplicaSet '%s' in namespace '%s' to be available finished. All replicas are available.", replicaSetName, namespace)
				return nil
			}

			// If ReplicaSet is in a failed state, exit early
			if replicaSet.Status.Conditions != nil {
				for _, condition := range replicaSet.Status.Conditions {
					if condition.Type == appsv1.ReplicaSetReplicaFailure && condition.Status == corev1.ConditionTrue {
						return tracederrors.TracedErrorf("ReplicaSet '%s' in namespace '%s' has replica failure: %s", replicaSetName, namespace, condition.Message)
					}
				}
			}
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for ReplicaSet '%s' in namespace '%s' to be available", replicaSetName, namespace)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}

func ListReplicaSets(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	replicaSetList, err := clientset.AppsV1().ReplicaSets(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list ReplicaSets in namespace '%s'.", namespaceName)
	}

	replicaSetNames := []string{}
	for _, rs := range replicaSetList.Items {
		replicaSetNames = append(replicaSetNames, rs.Name)
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' ReplicaSets in namespace '%s'.", len(replicaSetNames), namespaceName)

	return replicaSetNames, nil
}
