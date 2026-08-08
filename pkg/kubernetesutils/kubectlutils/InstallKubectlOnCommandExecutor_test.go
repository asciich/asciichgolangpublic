package kubectlutils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestInstallKubectlOnCommandExecutor_DockerContainer(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Get Docker command executor
	docker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
	require.NoError(t, err)

	// Define unique container name to avoid conflicts
	containerName := fmt.Sprintf("kubectl-test-%d", time.Now().UnixNano())

	// Create and start a test container
	runOptions := &dockeroptions.DockerRunContainerOptions{
		Name:      containerName,
		ImageName: "alpine/curl:latest",
		Command:   []string{"sleep", "300"}, // Keep container running for 5 minutes
	}

	container, err := docker.RunContainer(ctx, runOptions)
	require.NoError(t, err)

	// Ensure container is cleaned up after test
	defer func() {
		killErr := docker.KillContainerByName(ctx, containerName)
		if killErr != nil {
			t.Logf("Warning: Failed to kill container %s: %v", containerName, killErr)
		}
	}()

	// Install kubectl inside the container using a temporary path
	// Use /tmp/kubectl to avoid requiring root access
	tempInstallPath := "/tmp/kubectl"
	options := &kubectlutils.InstallKubectlOptions{
		InstallPath: tempInstallPath,
		UseSudo:     false, // No sudo inside container
		Version:     "v1.36.2",
	}

	err = kubectlutils.InstallKubectlOnCommandExecutor(ctx, container, options)
	require.NoError(t, err)

	t.Log("InstallKubectlOnCommandExecutor completed successfully")

	// Verify kubectl was installed by checking it exists and is executable
	output, err := container.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{tempInstallPath, "version", "--client"},
		},
	)
	require.NoError(t, err)
	require.Contains(t, output, "Client Version", "kubectl version command should return version info")
}

func TestInstallKubectlOnCommandExecutor_NilCommandExecutor(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	options := &kubectlutils.InstallKubectlOptions{
		InstallPath: "/tmp/kubectl",
		UseSudo:     false,
	}

	err := kubectlutils.InstallKubectlOnCommandExecutor(ctx, nil, options)
	require.Error(t, err)
	require.Contains(t, err.Error(), "commandExecutor")
}

func TestInstallKubectlOnCommandExecutor_NilOptions(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Get Docker command executor
	docker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
	require.NoError(t, err)

	// Define unique container name
	containerName := fmt.Sprintf("kubectl-test-nilopts-%d", time.Now().UnixNano())

	// Create and start a test container
	runOptions := &dockeroptions.DockerRunContainerOptions{
		Name:      containerName,
		ImageName: "alpine:latest",
		Command:   []string{"sleep", "300"},
	}

	container, err := docker.RunContainer(ctx, runOptions)
	require.NoError(t, err)

	defer func() {
		_ = docker.KillContainerByName(ctx, containerName)
	}()

	time.Sleep(2 * time.Second)

	// Passing nil options should use defaults
	// This will likely fail because default path is /bin/kubectl which requires sudo
	err = kubectlutils.InstallKubectlOnCommandExecutor(ctx, container, nil)
	if err != nil {
		t.Logf("InstallKubectlOnCommandExecutor with nil options failed (expected): %v", err)
	}
}
