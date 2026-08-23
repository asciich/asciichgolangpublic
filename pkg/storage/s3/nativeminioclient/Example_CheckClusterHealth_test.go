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

func Test_Example_CheckClusterHealth_test(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Define admin credentials for the test environment
	const minioAdminUser = "minioadmin"
	minioAdminPassword, err := randomgenerator.GetRandomString(10)
	require.NoError(t, err)

	// Run minio in a docker container for testing.
	const containerName = "test-nativeminioclient-clusterhealth"
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

	// Check cluster health:
	clusterHealth, err := nativeminioclient.CheckClusterHealth(ctx, adminClient)
	require.NoError(t, err)

	// Verify cluster health is returned:
	require.NotNil(t, clusterHealth)

	// Print cluster health:
	fmt.Printf("Cluster Health Status\n")
	fmt.Printf("=====================\n")
	fmt.Printf("Is Healthy: %v\n", clusterHealth.IsHealthy)
	fmt.Printf("Total servers: %d\n", clusterHealth.TotalServers)
	fmt.Printf("Online servers: %d\n", clusterHealth.OnlineServers)
	fmt.Printf("Offline servers: %d\n", clusterHealth.OfflineServers)
	fmt.Printf("Total disks: %d\n", clusterHealth.TotalDisks)
	fmt.Printf("Online disks: %d\n", clusterHealth.OnlineDisks)
	fmt.Printf("Offline disks: %d\n", clusterHealth.OfflineDisks)

	if len(clusterHealth.ServerWarnings) > 0 {
		fmt.Println("\nServer Warnings:")
		for _, warning := range clusterHealth.ServerWarnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(clusterHealth.DiskWarnings) > 0 {
		fmt.Println("\nDisk Warnings:")
		for _, warning := range clusterHealth.DiskWarnings {
			fmt.Printf("  - %s\n", warning)
		}
	}
}
