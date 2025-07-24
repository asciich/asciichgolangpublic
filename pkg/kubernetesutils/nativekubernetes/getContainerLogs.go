package nativekubernetes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

var ErrOnlyCombinedStreamAvailable = errors.New("specify the stream is not supported. Only combinded stream is availalbe by the kubernetes API, feature flag seems deactivated")

func getContainerLogsForStream(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string, containerName string, stream string) ([]byte, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientSet")
	}

	if podName == "" {
		return nil, tracederrors.TracedErrorEmptyString("podName")
	}

	if namespace == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespace")
	}

	// containerName is opional and therefore not checked.

	knownStreams := []string{"All", "Stdout", "Stderr"}
	if !slices.Contains(knownStreams, stream) {
		return nil, tracederrors.TracedErrorf("Unknown stream to get log: %s, known streams are: %s", knownStreams, stream)
	}

	logging.LogInfoByCtxf(ctx, "Get logs for stream '%s' of container '%s' in pod '%s' in namespace '%s' started.", stream, containerName, podName, namespace)

	podLogOptions := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    false, // Do not follow the stream, read complete output
	}

	if stream != "All" {
		podLogOptions.Stream = &stream
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, podLogOptions)

	podLogs, err := req.Stream(ctx)
	if err != nil {
		if apierrors.IsInvalid(err) && strings.Contains(err.Error(), "stream: Forbidden: may not be specified") {
			return nil, tracederrors.TracedErrorf("Unable to get log stream '%s' for container '%s' in pod '%s' in namespace '%s': %w", stream, containerName, podName, namespace, ErrOnlyCombinedStreamAvailable)
		}
		return nil, tracederrors.TracedErrorf("error opening log stream for container '%s' in pod '%s' in namespace '%s': %w", containerName, podName, namespace, err)
	}
	defer func() {
		if closeErr := podLogs.Close(); closeErr != nil {
			logging.LogErrorByCtxf(ctx, "Error closing log stream for container '%s' in pod '%s' in namespace '%s': %v", containerName, podName, namespace, closeErr)
		}
	}()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to copy log buffer for container '%s' in pod '%s' in namespace '%s': %w", containerName, podName, namespace, err)
	}

	logging.LogInfoByCtxf(ctx, "Get logs for stream '%s' of container '%s' in pod '%s' in namespace '%s' finished.", stream, containerName, podName, namespace)

	return buf.Bytes(), nil
}

func GetContainerStdoutLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string, containerName string) ([]byte, error) {
	return getContainerLogsForStream(ctx, clientset, namespace, podName, containerName, "Stdout")
}

func GetContainerStderrLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string, containerName string) ([]byte, error) {
	return getContainerLogsForStream(ctx, clientset, namespace, podName, containerName, "Stderr")
}

// 'All' combines stdout and stderr.
func GetContainerAllLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string, containerName string) ([]byte, error) {
	return getContainerLogsForStream(ctx, clientset, namespace, podName, containerName, "All")
}

func GetContainerLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string, containerName string) (stdout []byte, stderr []byte, err error) {
	stdout, err = GetContainerStdoutLogs(ctx, clientset, namespace, podName, containerName)
	if err != nil {
		if errors.Is(err, ErrOnlyCombinedStreamAvailable) {
			logging.LogInfoByCtxf(ctx, "Only combined stream available by the kubernetes API")
			stdout, err = GetContainerAllLogs(ctx, clientset, namespace, podName, containerName)
			if err != nil {
				return nil, nil, err
			}

			return stdout, []byte{}, nil
		}
		return nil, nil, err
	}

	stderr, err = GetContainerStderrLogs(ctx, clientset, namespace, podName, containerName)
	if err != nil {
		return nil, nil, err
	}

	return stdout, stderr, nil
}
