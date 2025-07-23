package nativekubernetes

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func RunCommandInTemporaryPod(ctx context.Context, config *rest.Config, options *kubernetesparameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
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

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	command, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' started.", podName, namespace, imageName)

	containerName := podName
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to get client set: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Going to start pod '%s' in namespace '%s' using container image '%s'.", podName, namespace, imageName)
	_, err = clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) && options.DeleteAlreadyExistingPod {
			logging.LogInfoByCtxf(ctx, "Going to delete pod already existing pod '%s' in namespace '%s' before running command.", podName, namespace)
			err = DeletePod(ctx, clientset, podName, namespace)
			if err != nil {
				return nil, err
			}
			_, err = clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
			if err != nil {
				return nil, tracederrors.TracedErrorf("Error creating Pod: %w", err)
			}
		} else {
			return nil, tracederrors.TracedErrorf("Error creating Pod: %w", err)
		}
	}

	// Ensure pod is deleted after executing the command
	defer func() {
		_ = DeletePod(ctx, clientset, podName, namespace)
	}()

	err = waitForPodSucceeded(ctx, clientset, namespace, podName, time.Minute*1)
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
