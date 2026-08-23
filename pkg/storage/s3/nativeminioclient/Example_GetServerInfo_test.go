package nativeminioclient_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
)

func Test_Example_GetServerInfo_test(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Define admin credentials for the test environment
	const minioAdminUser = "minioadmin"
	minioAdminPassword, err := randomgenerator.GetRandomString(10)
	require.NoError(t, err)

	// Run minio in a docker container for testing.
	const containerName = "test-nativeminioclient-serverinfo"
	err = nativedocker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	container, err := nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
		Name:      containerName,
		ImageName: "quay.io/minio/minio",
		Command:   []string{"server", "/data", "--console-address", ":9001"},
		Ports:     []string{"9000:9000"},
		AdditionalEnvVars: map[string]string{
			"MINIO_ROOT_USER":     minioAdminUser,
			"MINIO_ROOT_PASSWORD": minioAdminPassword,
		},
		WaitForPortsOpen: true,
	})
	require.NoError(t, err)
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	time.Sleep(time.Second * 2)

	// Get the minio admin client:
	adminClient, err := nativeminioclient.NewAdminClient("localhost:9000", minioAdminUser, minioAdminPassword, &s3options.NewS3ClientOptions{})
	require.NoError(t, err)

	// Get server information:
	serverInfo, err := nativeminioclient.GetServerInfo(ctx, adminClient)
	require.NoError(t, err)

	// Verify server info is returned:
	require.NotNil(t, serverInfo)
	require.NotEmpty(t, serverInfo.Servers)

	// Print server information:
	fmt.Printf("Total servers: %d\n", len(serverInfo.Servers))
	for i, server := range serverInfo.Servers {
		fmt.Printf("Server %d:\n", i+1)
		fmt.Printf("  Endpoint: %s\n", server.Endpoint)
		fmt.Printf("  State: %s\n", server.State)
		fmt.Printf("  Uptime: %d seconds\n", server.Uptime)
		fmt.Printf("  Num Disks: %d\n", len(server.Disks))
	}
}
