package nativekubernetes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils/tarparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
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
		err := WaitForPodDeleted(ctx, clientset, namespace, podName, time.Second*100)
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

func Exec(ctx context.Context, config *rest.Config, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if config == nil {
		return nil, tracederrors.TracedErrorNil("config")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
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

	logging.LogInfoByCtxf(ctx, "Exec command in container '%s' of pod '%s' in namespace '%s' started.", containerName, podName, namespaceName)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     options.IsStinDataAvailable(),
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create exec: %s", err)
	}
	var stdout, stderr bytes.Buffer

	streamOptions := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if options.IsStinDataAvailable() {
		streamOptions.Stdin = bytes.NewReader(options.StdinBytes)
	}

	err = exec.StreamWithContext(ctx, streamOptions)
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

	logging.LogInfoByCtxf(ctx, "Exec command in container '%s' of pod '%s' in namespace '%s' finished.", containerName, podName, namespaceName)

	return output, nil
}

func CreatePod(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
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

	logging.LogInfoByCtxf(ctx, "Create pod '%s' in namespace '%s' using container image '%s' started.", podName, namespaceName, imageName)

	// Build environment variables from secrets
	envVars := []corev1.EnvVar{}
	if options.SecretEnvVars != nil {
		for envVarName, secretSource := range options.SecretEnvVars {
			envVars = append(envVars, corev1.EnvVar{
				Name: envVarName,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: secretSource.SecretName,
						},
						Key: secretSource.SecretKey,
					},
				},
			})
		}
	}

	// Build volumes and volume mounts from secrets
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}
	if options.SecretMounts != nil {
		for mountPath, secretSource := range options.SecretMounts {
			volumeName := "secret-" + secretSource.SecretName
			volumes = append(volumes, corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: secretSource.SecretName,
					},
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      volumeName,
				MountPath: mountPath,
				ReadOnly:  true,
			})
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:         containerName,
					Image:        imageName,
					Command:      command,
					Stdin:        true,
					TTY:          true,
					Env:          envVars,
					VolumeMounts: volumeMounts,
				},
			},
			Volumes:       volumes,
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	logging.LogInfoByCtxf(ctx, "Going to start pod '%s' in namespace '%s' using container image '%s'.", podName, namespaceName, imageName)
	_, err = clientset.CoreV1().Pods(namespaceName).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) && options.DeleteAlreadyExistingPod {
			logging.LogInfoByCtxf(ctx, "Going to delete pod already existing pod '%s' in namespace '%s' before running command.", podName, namespaceName)
			err = DeletePod(ctx, clientset, podName, namespaceName)
			if err != nil {
				return err
			}
			_, err = clientset.CoreV1().Pods(namespaceName).Create(ctx, pod, metav1.CreateOptions{})
			if err != nil {
				return tracederrors.TracedErrorf("Error creating Pod: %w", err)
			}
		} else {
			return tracederrors.TracedErrorf("Error creating Pod: %w", err)
		}
	}

	if options.WaitForPodRunning {
		err = WaitForPodRunning(ctx, clientset, namespaceName, podName, time.Minute*1)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Create pod '%s' in namespace '%s' using container image '%s' finished.", podName, namespaceName, imageName)

	return nil
}

func RunCommandInTemporaryPod(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
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

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' started.", podName, namespaceName, imageName)

	err = CreatePod(ctx, clientset, namespaceName, options)
	if err != nil {
		return nil, err
	}

	// Ensure pod is deleted after executing the command
	defer func() {
		_ = DeletePod(ctx, clientset, podName, namespaceName)
	}()

	err = WaitForPodSucceeded(ctx, clientset, namespaceName, podName, time.Minute*1)
	if err != nil {
		return nil, err
	}

	stdout, stderr, err := GetContainerLogs(ctx, clientset, namespaceName, podName, containerName)
	if err != nil {
		return nil, err
	}

	var retVal = 0
	output := &commandoutput.CommandOutput{
		ReturnCode: &retVal,
		Stdout:     &stdout,
		Stderr:     &stderr,
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' finished.", podName, namespaceName, imageName)

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

	// Verify that the pod exists
	_, err = clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return tracederrors.TracedErrorf("pod '%s' not found in namespace '%s'", podName, namespaceName)
		}
		return tracederrors.TracedErrorf("failed to get pod '%s' in namespace '%s': %w", podName, namespaceName, err)
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
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  tarReader,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return tracederrors.TracedErrorf("failed to copy file to pod: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Copy local file '%s' as '%s' into container '%s' of pod '%s' of namespace '%s' finished.", localFile, destPath, containerName, podName, namespaceName)

	return nil
}

// CopyFileFromPod copies a file from a container in a pod to the local filesystem
// Similar to: kubectl cp <namespace>/<pod>:<container> <local-path>
func CopyFileFromPod(ctx context.Context, config *rest.Config, podName string, namespaceName string, containerName string, srcPath string, destFile string) error {
	if config == nil {
		return tracederrors.TracedErrorNil("config")
	}
	if podName == "" {
		return tracederrors.TracedErrorEmptyString("podName")
	}
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}
	if containerName == "" {
		return tracederrors.TracedErrorEmptyString("containerName")
	}
	if srcPath == "" {
		return tracederrors.TracedErrorEmptyString("srcPath")
	}
	if destFile == "" {
		return tracederrors.TracedErrorEmptyString("destFile")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return tracederrors.TracedErrorf("failed to create kubernetes client: %w", err)
	}

	// Validate pod exists before attempting copy
	_, err = clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("pod '%s' not found in namespace '%s': %w", podName, namespaceName, err)
	}

	logging.LogInfoByCtxf(ctx, "Copy file '%s' from container '%s' of pod '%s' of namespace '%s' to local '%s' started.", srcPath, containerName, podName, namespaceName, destFile)

	// Create the destination directory if it doesn't exist
	destDir := filepath.Dir(destFile)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return tracederrors.TracedErrorf("failed to create destination directory '%s': %w", destDir, err)
	}

	// Create destination file
	destFileHandle, err := os.Create(destFile)
	if err != nil {
		return tracederrors.TracedErrorf("failed to create destination file '%s': %w", destFile, err)
	}
	defer destFileHandle.Close()

	// Use kubectl exec to cat the file and stream it back
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"cat", srcPath},
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return tracederrors.TracedErrorf("failed to create executor: %w", err)
	}

	var stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: destFileHandle,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return tracederrors.TracedErrorf("failed to copy file from pod: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Copy file '%s' from container '%s' of pod '%s' of namespace '%s' to local '%s' finished.", srcPath, containerName, podName, namespaceName, destFile)

	return nil
}

func ListPods(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	podList, err := clientset.CoreV1().Pods(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list pods in namespace '%s'.", namespaceName)
	}

	podNames := []string{}
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' pods in namespace '%s'.", len(podNames), namespaceName)

	return podNames, nil
}

func ListPodNames(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	return ListPods(ctx, clientset, namespaceName)
}
