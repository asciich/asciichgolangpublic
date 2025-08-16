package nativekubernetes

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils/tarparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
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
		err := WaitForPodDeleted(ctx, clientset, namespace, podName, time.Second*30)
		if err != nil {
			return err
		}
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

func Exec(ctx context.Context, config *rest.Config, options *kubernetesparameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if config == nil {
		return nil, tracederrors.TracedErrorNil("config")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	namespace, err := options.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	podName, err := options.GetPodName()
	if err != nil {
		return nil, err
	}

	command, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Exec command in container '%s' of pod '%s' in namespace '%s' started.", containerName, podName, namespace)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create exec: %s", err)
	}
	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error executing command: %s", err)
	}

	stdoutBytes := stdout.Bytes()
	stderrBytes := stderr.Bytes()
	var retVal int

	output := &commandoutput.CommandOutput{
		Stdout:     &stdoutBytes,
		Stderr:     &stderrBytes,
		ReturnCode: &retVal,
	}

	logging.LogInfoByCtxf(ctx, "Exec command in container '%s' of pod '%s' in namespace '%s' finished.", containerName, podName, namespace)

	return output, nil
}

func CreatePod(ctx context.Context, clientset *kubernetes.Clientset, options *kubernetesparameteroptions.RunCommandOptions) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	namespace, err := options.GetNamespaceName()
	if err != nil {
		return err
	}

	podName, err := options.GetPodName()
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

	logging.LogInfoByCtxf(ctx, "Create pod '%s' in namespace '%s' using container image '%s' started.", podName, namespace, imageName)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
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
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	logging.LogInfoByCtxf(ctx, "Going to start pod '%s' in namespace '%s' using container image '%s'.", podName, namespace, imageName)
	_, err = clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) && options.DeleteAlreadyExistingPod {
			logging.LogInfoByCtxf(ctx, "Going to delete pod already existing pod '%s' in namespace '%s' before running command.", podName, namespace)
			err = DeletePod(ctx, clientset, podName, namespace)
			if err != nil {
				return err
			}
			_, err = clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
			if err != nil {
				return tracederrors.TracedErrorf("Error creating Pod: %w", err)
			}
		} else {
			return tracederrors.TracedErrorf("Error creating Pod: %w", err)
		}
	}

	if options.WaitForPodRunning {
		err = WaitForPodRunning(ctx, clientset, namespace, podName, time.Minute*1)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Create pod '%s' in namespace '%s' using container image '%s' finished.", podName, namespace, imageName)

	return nil
}

func RunCommandInTemporaryPod(ctx context.Context, clientset *kubernetes.Clientset, options *kubernetesparameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("config")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	namespace, err := options.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	podName, err := options.GetPodName()
	if err != nil {
		return nil, err
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return nil, err
	}

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' started.", podName, namespace, imageName)

	err = CreatePod(ctx, clientset, options)
	if err != nil {
		return nil, err
	}

	// Ensure pod is deleted after executing the command
	defer func() {
		_ = DeletePod(ctx, clientset, podName, namespace)
	}()

	err = WaitForPodSucceeded(ctx, clientset, namespace, podName, time.Minute*1)
	if err != nil {
		return nil, err
	}

	stdout, stderr, err := GetContainerLogs(ctx, clientset, namespace, podName, containerName)
	if err != nil {
		return nil, err
	}

	var retVal = 0
	output := &commandoutput.CommandOutput{
		ReturnCode: &retVal,
		Stdout:     &stdout,
		Stderr:     &stderr,
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' finished.", podName, namespace, imageName)

	return output, nil
}

func PodExists(ctx context.Context, clientSet *kubernetes.Clientset, podName string, namespace string) (bool, error) {
	if clientSet == nil {
		return false, tracederrors.TracedErrorNil("clientSet")
	}

	if podName == "" {
		return false, tracederrors.TracedErrorEmptyString("podName")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	_, err := clientSet.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return false, tracederrors.TracedErrorf("Failed to get pod '%s' in namespace '%s' to check if exists: %w", podName, namespace, err)
		}
	}

	exists := err == nil

	if exists {
		logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' exists.", podName, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' does not exist.", podName, namespace)
	}

	return exists, nil
}

func WaitForPodDeleted(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, podName string, timeout time.Duration) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientSet")
	}

	if podName == "" {
		return tracederrors.TracedErrorEmptyString("podName")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for pod '%s' in namespace '%s' to be deleted started.", podName, namespaceName)

	_, err := clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' is already deleted.", podName, err)
		return nil
	}

	w, err := clientset.CoreV1().Pods(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podName),
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' is already deleted.", podName, err)
		} else {
			return fmt.Errorf("failed to set up watch for pod %s: %w", podName, err)
		}
	}
	defer w.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case _, ok := <-w.ResultChan():
			if !ok {
				return tracederrors.TracedErrorf("watch channel closed unexpectedly when waiting for pod '%s' in namespace '%s' to be deleted", podName, namespaceName)
			}

			_, err := clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' is now deleted.", podName, namespaceName)
				return nil
			}

			logging.LogInfoByCtxf(ctx, "Still waiting for pod '%s' in namespace '%s' to be deleted.", podName, namespaceName)
		case <-timer.C:
			return tracederrors.TracedErrorf("timeout waiting for pod '%s' in namespace '%s' to be deleted", podName, namespaceName)
		case <-ctx.Done():
			return ctx.Err() // Context was cancelled
		}
	}
}

func WaitForPodRunning(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string, timeout time.Duration) error {
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

func WaitForPodSucceeded(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string, timeout time.Duration) error {
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

func CopyFileToPod(ctx context.Context, config *rest.Config, localFile string, destPath string, podName string, containerName string, namespaceName string) error {
	if config == nil {
		return tracederrors.TracedErrorNil("config")
	}

	if localFile == "" {
		return tracederrors.TracedErrorEmptyString("localFile")
	}

	if destPath == "" {
		return tracederrors.TracedErrorEmptyString("destPath")
	}

	if podName == "" {
		return tracederrors.TracedErrorEmptyString("podName")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	logging.LogInfoByCtxf(ctx, "Copy local file '%s' as '%s' into container '%s' of pod '%s' of namespace '%s' started.", localFile, destPath, containerName, podName, namespaceName)

	tarReader, err := tarutils.FileToTarReader(localFile, &tarparameteroptions.FileToTarOptions{
		OverrideFileName: filepath.Base(destPath),
	})
	if err != nil {
		return err
	}

	clientset, err := GetClientSetFromRestConfig(ctx, config)
	if err != nil {
		return err
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName).
		VersionedParams(&corev1.PodExecOptions{
			Command: []string{"tar", "xf", "-", "-C", filepath.Dir(destPath)},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  tarReader,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {

	}

	logging.LogInfoByCtxf(ctx, "Copy local file '%s' as '%s' into container '%s' of pod '%s' of namespace '%s' finished.", localFile, destPath, containerName, podName, namespaceName)

	return nil
}
