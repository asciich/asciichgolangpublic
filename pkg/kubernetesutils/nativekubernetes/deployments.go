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

func DeleteDeployment(ctx context.Context, clientset *kubernetes.Clientset, deploymentName string, namespace string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if deploymentName == "" {
		return tracederrors.TracedErrorEmptyString("deploymentName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Delete Deployment '%s' in namespace '%s' started.", deploymentName, namespace)

	deletePolicy := metav1.DeletePropagationBackground
	err := clientset.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err == nil {
		err := WaitForDeploymentDeleted(ctx, clientset, namespace, deploymentName, time.Second*30)
		if err != nil {
			return err
		}
		logging.LogChangedByCtxf(ctx, "Deployment '%s' in namespace '%s' deleted.", deploymentName, namespace)
	} else {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Deployment '%s' already absent in namespace '%s'.", deploymentName, namespace)
		} else {
			return tracederrors.TracedErrorf("Failed to delete Deployment '%s' in namespace '%s': %w", deploymentName, namespace, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Delete Deployment '%s' in namespace '%s' finished.", deploymentName, namespace)

	return nil
}

func CreateDeployment(ctx context.Context, config *rest.Config, options *kubernetesparameteroptions.RunCommandOptions) error {
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

	deploymentName, err := options.GetDeploymentName()
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

	logging.LogInfoByCtxf(ctx, "Create Deployment '%s' in namespace '%s' using container image '%s' started.", deploymentName, namespace, imageName)

	clientset, err := GetClientSetFromRestConfig(ctx, config)
	if err != nil {
		return err
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
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

	logging.LogInfoByCtxf(ctx, "Going to start Deployment '%s' in namespace '%s' using container image '%s'.", deploymentName, namespace, imageName)
	_, err = clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) && options.DeleteAlreadyExistingDeployment {
			logging.LogInfoByCtxf(ctx, "Going to delete already existing Deployment '%s' in namespace '%s' before running command.", deploymentName, namespace)
			err = DeleteDeployment(ctx, clientset, deploymentName, namespace)
			if err != nil {
				return err
			}
			_, err = clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
			if err != nil {
				return tracederrors.TracedErrorf("Error creating Deployment: %w", err)
			}
		} else {
			return tracederrors.TracedErrorf("Error creating Deployment: %w", err)
		}
	}

	if options.WaitForDeploymentAvailable {
		err = WaitForDeploymentAvailable(ctx, clientset, namespace, deploymentName, time.Minute*2)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Create Deployment '%s' in namespace '%s' using container image '%s' finished.", deploymentName, namespace, imageName)

	return nil
}

func DeploymentExists(ctx context.Context, clientSet *kubernetes.Clientset, deploymentName string, namespace string) (bool, error) {
	if clientSet == nil {
		return false, tracederrors.TracedErrorNil("clientSet")
	}

	if deploymentName == "" {
		return false, tracederrors.TracedErrorEmptyString("deploymentName")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Check if Deployment '%s' in namespace '%s' exists.", deploymentName, namespace)

	_, err := clientSet.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' does not exist.", deploymentName, namespace)
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Failed to get Deployment '%s' in namespace '%s' to check if exists: %w", deploymentName, namespace, err)
	}

	logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' exists.", deploymentName, namespace)
	return true, nil
}

func WaitForDeploymentDeleted(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, deploymentName string, timeout time.Duration) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if deploymentName == "" {
		return tracederrors.TracedErrorEmptyString("deploymentName")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for Deployment '%s' in namespace '%s' to be deleted started.", deploymentName, namespaceName)

	_, err := clientset.AppsV1().Deployments(namespaceName).Get(ctx, deploymentName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' is already deleted.", deploymentName, namespaceName)
		return nil
	}

	w, err := clientset.AppsV1().Deployments(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", deploymentName),
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' is already deleted.", deploymentName, namespaceName)
		} else {
			return fmt.Errorf("failed to set up watch for Deployment %s: %w", deploymentName, err)
		}
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case _, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for Deployment '%s' in namespace '%s' to be deleted", deploymentName, namespaceName)
			}

			_, err := clientset.AppsV1().Deployments(namespaceName).Get(ctx, deploymentName, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' is now deleted.", deploymentName, namespaceName)
				return nil
			}

			logging.LogInfoByCtxf(ctx, "Still waiting for Deployment '%s' in namespace '%s' to be deleted.", deploymentName, namespaceName)
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for Deployment '%s' in namespace '%s' to be deleted", deploymentName, namespaceName)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}

func WaitForDeploymentAvailable(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, deploymentName string, timeout time.Duration) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if deploymentName == "" {
		return tracederrors.TracedErrorEmptyString("deploymentName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for Deployment '%s' in namespace '%s' to be available started.", deploymentName, namespace)

	w, err := clientSet.AppsV1().Deployments(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", deploymentName),
	})
	if err != nil {
		return fmt.Errorf("failed to set up watch for Deployment %s: %w", deploymentName, err)
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for Deployment '%s' in namespace '%s' to be available", deploymentName, namespace)
			}

			deployment, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				continue
			}

			// Check if the desired number of replicas are available
			if deployment.Status.AvailableReplicas == *deployment.Spec.Replicas {
				logging.LogInfoByCtxf(ctx, "Wait for Deployment '%s' in namespace '%s' to be available finished. All replicas are available.", deploymentName, namespace)
				return nil
			}

			// If Deployment is in a failed state, exit early
			if deployment.Status.Conditions != nil {
				for _, condition := range deployment.Status.Conditions {
					if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == corev1.ConditionTrue {
						return tracederrors.TracedErrorf("Deployment '%s' in namespace '%s' has replica failure: %s", deploymentName, namespace, condition.Message)
					}
				}
			}
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for Deployment '%s' in namespace '%s' to be available", deploymentName, namespace)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}
