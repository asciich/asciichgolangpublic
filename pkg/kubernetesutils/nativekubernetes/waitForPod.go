package nativekubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func waitForPodRunning(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string, timeout time.Duration) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if podName == "" {
		return tracederrors.TracedErrorEmptyString("podName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for pod '%s' in namespace '%s' to be running started.", podName, namespace)

	w, err := clientSet.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podName),
	})
	if err != nil {
		return fmt.Errorf("failed to set up watch for pod %s: %w", podName, err)
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for pod '%s' in namespace '%s' to be running", podName, namespace)
			}

			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			if pod.Status.Phase == corev1.PodRunning {
				logging.LogInfoByCtxf(ctx, "Wait for pod '%s' in namespace '%s' to be running finished. Pod is now running", podName, namespace)
				return nil
			}
			// If pod is in a failed state, exit early
			if pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodSucceeded {
				return tracederrors.TracedErrorf("pod '%s' in namespace '%s' entered phase %s before running", podName, namespace, pod.Status.Phase)
			}
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for pod '%s' in namespace '%s' to be running", podName, namespace)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}

}

func waitForPodSucceeded(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string, timeout time.Duration) error {
	if clientSet == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if podName == "" {
		return tracederrors.TracedErrorEmptyString("podName")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for pod '%s' in namespace '%s' to be succeeded started.", podName, namespace)

	w, err := clientSet.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podName),
	})
	if err != nil {
		return fmt.Errorf("failed to set up watch for pod %s: %w", podName, err)
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for pod '%s' in namespace '%s' to be succeeded", podName, namespace)
			}

			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			if pod.Status.Phase == corev1.PodSucceeded {
				logging.LogInfoByCtxf(ctx, "Wait for pod '%s' in namespace '%s' to be succeeded finished. Pod is now succeeded.", podName, namespace)
				return nil
			}
			// If pod is in a failed state, exit early
			if pod.Status.Phase == corev1.PodFailed {
				return tracederrors.TracedErrorf("pod '%s' in namespace '%s' failed", podName, namespace)
			}
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for pod '%s' in namespace '%s' to be succeeded", podName, namespace)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}
