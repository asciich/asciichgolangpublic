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

func Test_Example_CheckServerOnlineStatus_test(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Define admin credentials for the test environment
	const minioAdminUser = "minioadmin"
	minioAdminPassword, err := randomgenerator.GetRandomString(10)
	require.NoError(t, err)

	// Run minio in a docker container for testing.
	const containerName = "test-nativeminioclient-serverstatus"
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

	// Check server online status:
	serverStatus, err := nativeminioclient.CheckServerOnlineStatus(ctx, adminClient)
	require.NoError(t, err)

	// Verify server status is returned:
	require.NotNil(t, serverStatus)
	require.Greater(t, serverStatus.TotalServers, 0)

	// Print server status:
	fmt.Printf("Total servers: %d\n", serverStatus.TotalServers)
	fmt.Printf("Online servers: %d\n", serverStatus.OnlineServers)
	fmt.Printf("Offline servers: %d\n", serverStatus.OfflineServers)
	fmt.Println("\nServer details:")
	for i, status := range serverStatus.ServerStatuses {
		fmt.Printf("  Server %d:\n", i+1)
		fmt.Printf("    Endpoint: %s\n", status.Endpoint)
		fmt.Printf("    Is Online: %v\n", status.IsOnline)
		fmt.Printf("    Uptime: %d seconds\n", status.Uptime)
	}
}
