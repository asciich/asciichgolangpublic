package commandexecutoriscsi_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils/commandexecutorchecksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/storage/iscsiutils/commandexecutoriscsi"
)

func Test_EnsureInitiorNameConfig(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Use a temporary container for testing
	const containerName = "test-iscsi-initiator-config"

	// Ensure the container is absent before we start
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	// Start a container where we can test the initiator name configuration
	// We use a simple alpine container that stays running
	container, err := docker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:                 containerName,
			Command:              []string{"sleep", "60s"},
			ImageName:            "alpine:latest",
			KeepStoppedContainer: true,
		},
	)
	require.NoError(t, err)

	_, err = container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"apk", "add", "--no-cache", "sudo"},
	})
	require.NoError(t, err)

	// In any case we delete the container after this test
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	useSudo := []struct {
		UseSudo bool
	}{
		{true},
		{false},
	}

	for _, tt := range useSudo {
		t.Run(fmt.Sprintf("creates initiator name file when it does not exist, useSudo=%v", tt.UseSudo), func(t *testing.T) {
			const path = "/etc/iscsi/initiatorname.iscsi"

			// Ensure /etc/iscsi/ absent
			err := commandexecutorfile.DeleteDirectory(ctx, container, filepath.Dir(path), &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			exists, err := commandexecutorfile.Exists(ctx, container, path)
			require.NoError(t, err)
			require.False(t, exists)

			// Execute function under test.
			// It must create the missing file.
			changedCtx := contextutils.WithChangeIndicator(ctx)
			err = commandexecutoriscsi.EnsureInitiorNameConfig(changedCtx, container, &filesoptions.CreateOptions{UseSudo: tt.UseSudo})
			require.NoError(t, err)
			require.True(t, contextutils.IsChanged(changedCtx))

			exists, err = commandexecutorfile.Exists(ctx, container, path)
			require.NoError(t, err)
			require.True(t, exists)

			checksum, err := commandexecutorchecksumutils.GetMD5SumFromFileByPath(ctx, container, path)
			require.NoError(t, err)

			// Idempotence check: Running again should not change it:
			changedCtx = contextutils.WithChangeIndicator(ctx)
			err = commandexecutoriscsi.EnsureInitiorNameConfig(changedCtx, container, &filesoptions.CreateOptions{UseSudo: tt.UseSudo})
			require.NoError(t, err)
			require.False(t, contextutils.IsChanged(changedCtx))

			exists, err = commandexecutorfile.Exists(ctx, container, path)
			require.NoError(t, err)
			require.True(t, exists)

			checksum2, err := commandexecutorchecksumutils.GetMD5SumFromFileByPath(ctx, container, path)
			require.NoError(t, err)

			require.EqualValues(t, checksum, checksum2)

			// Check content:
			content, err := commandexecutorfile.ReadAsString(container, path)
			require.NoError(t, err)
			require.Contains(t, content, "InitiatorName=iqn.2016-04.com.open-iscsi:")
		})
	}
}
