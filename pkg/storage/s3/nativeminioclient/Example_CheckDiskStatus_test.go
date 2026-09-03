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

func Test_Example_CheckDiskStatus_test(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Define admin credentials for the test environment
	const minioAdminUser = "minioadmin"
	minioAdminPassword, err := randomgenerator.GetRandomString(10)
	require.NoError(t, err)

	// Run minio in a docker container for testing.
	const containerName = "test-nativeminioclient-diskstatus"
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

	// Check disk status:
	diskStatuses, err := nativeminioclient.CheckDiskStatus(ctx, adminClient)
	require.NoError(t, err)

	// Verify disk status is returned:
	require.NotNil(t, diskStatuses)
	require.NotEmpty(t, diskStatuses)

	// Print disk status:
	for _, serverResult := range diskStatuses {
		fmt.Printf("Server: %s\n", serverResult.ServerEndpoint)
		fmt.Printf("  Total disks: %d\n", serverResult.TotalDisks)
		fmt.Printf("  Online disks: %d\n", serverResult.OnlineDisks)
		fmt.Printf("  Offline disks: %d\n", serverResult.OfflineDisks)
		fmt.Println("  Disk details:")
		for i, disk := range serverResult.DiskStatuses {
			fmt.Printf("    Disk %d:\n", i+1)
			fmt.Printf("      Drive Path: %s\n", disk.DrivePath)
			fmt.Printf("      State: %s\n", disk.State)
			fmt.Printf("      Is Online: %v\n", disk.IsOnline)
			fmt.Printf("      Used Space: %d bytes\n", disk.UsedSpace)
			fmt.Printf("      Total Space: %d bytes\n", disk.TotalSpace)
			fmt.Printf("      Used Percent: %.2f%%\n", disk.UsedPercent)
		}
		fmt.Println()
	}
}
