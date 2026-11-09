package commandexecutorkubernetes

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (c *CommandExecutorNamespace) CreatePod(ctx context.Context, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (kubernetesinterfaces.Pod, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	podName, err := options.GetPodName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create kubernetes pod '%s' started.", podName)

	kubectlContext, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return nil, err
	}

	command, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	if options.DeleteAlreadyExistingPod {
		deleteCommand := []string{
			"kubectl", "delete", "pod", podName,
			"--context", kubectlContext,
			"--namespace", namespaceName,
			"--ignore-not-found",
		}

		_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: deleteCommand,
		})
		if err != nil {
			return nil, err
		}
	}

	volumes, volumeMounts := podVolumesFromOptions(options)

	container := map[string]any{
		"name":    containerName,
		"image":   imageName,
		"command": command,
		"stdin":   true,
		"tty":     true,
		"env":     json.RawMessage(toJsonEnvVarFromSecrets(options.SecretEnvVars)),
	}
	if len(volumeMounts) > 0 {
		container["volumeMounts"] = volumeMounts
	}

	podSpec := map[string]any{
		"restartPolicy": "Never",
		"containers":    []map[string]any{container},
	}
	if len(volumes) > 0 {
		podSpec["volumes"] = volumes
	}

	// Build the pod from a manifest instead of using `kubectl run --overrides`:
	// `kubectl run` always creates a container named like the pod. A strategic
	// merge override with a different container name does not replace it but
	// adds a second container. The first one (plain image without command)
	// exits immediately, so the pod never becomes Ready.
	manifest, err := json.Marshal(map[string]any{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name":      podName,
			"namespace": namespaceName,
			// Same label as `kubectl run` sets, to stay compatible.
			"labels": map[string]string{"run": podName},
		},
		"spec": podSpec,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal manifest of pod '%s': %w", podName, err)
	}

	// `create` (not `apply`) keeps the behavior of `kubectl run`:
	// it fails with AlreadyExists if the pod is already there.
	createCommand := []string{
		"kubectl", "create", "-f", "-",
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command:     createCommand,
		StdinString: string(manifest),
	})
	if err != nil {
		return nil, err
	}

	if options.WaitForPodRunning {
		waitCommand := []string{
			"kubectl", "wait", "--for=condition=Ready", "pod", podName,
			"--context", kubectlContext,
			"--namespace", namespaceName,
			"--timeout", "60s",
		}

		_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: waitCommand,
		})
		if err != nil {
			return nil, err
		}
	}

	logging.LogChangedByCtxf(ctx, "Created pod '%s' in namespace '%s'.", podName, namespaceName)

	logging.LogInfoByCtxf(ctx, "Create kubernetes pod '%s' finished.", podName)

	return c.GetPodByName(podName)
}

// podVolumesFromOptions builds the pod volumes and the container volume mounts
// for the secret mounts and persistent volume claims given in the options.
//
//   - Secrets are mounted read-only.
//   - Persistent volume claims are mounted writable since the executed command
//     may modify their content.
//
// Volume names are index based ("secret-0", "pvc-0", ...) in the order of the
// sorted mount paths. This keeps the output deterministic and avoids invalid
// names: secret and PVC names may contain dots, volume names must not.
// The result is shaped for the JSON manifest passed to kubectl.
func podVolumesFromOptions(options *kubernetesparameteroptions.KubernetesRunCommandOptions) (volumes []map[string]any, volumeMounts []map[string]any) {
	volumes = []map[string]any{}
	volumeMounts = []map[string]any{}

	secretMountPaths := make([]string, 0, len(options.SecretMounts))
	for mountPath := range options.SecretMounts {
		secretMountPaths = append(secretMountPaths, mountPath)
	}
	sort.Strings(secretMountPaths)

	for i, mountPath := range secretMountPaths {
		secretSource := options.SecretMounts[mountPath]
		volumeName := "secret-" + itoa(i)

		volumes = append(volumes, map[string]any{
			"name": volumeName,
			"secret": map[string]any{
				"secretName": secretSource.SecretName,
			},
		})
		volumeMounts = append(volumeMounts, map[string]any{
			"name":      volumeName,
			"mountPath": mountPath,
			"readOnly":  true,
		})
	}

	pvcMountPaths := make([]string, 0, len(options.PersistentVolumeClaims))
	for mountPath := range options.PersistentVolumeClaims {
		pvcMountPaths = append(pvcMountPaths, mountPath)
	}
	sort.Strings(pvcMountPaths)

	for i, mountPath := range pvcMountPaths {
		pvcSource := options.PersistentVolumeClaims[mountPath]
		volumeName := "pvc-" + itoa(i)

		volumes = append(volumes, map[string]any{
			"name": volumeName,
			"persistentVolumeClaim": map[string]any{
				"claimName": pvcSource.PVCName,
			},
		})
		volumeMounts = append(volumeMounts, map[string]any{
			"name":      volumeName,
			"mountPath": mountPath,
		})
	}

	return volumes, volumeMounts
}

// itoa converts a small non-negative index to its decimal string.
func itoa(i int) string {
	b, _ := json.Marshal(i)
	return string(b)
}
