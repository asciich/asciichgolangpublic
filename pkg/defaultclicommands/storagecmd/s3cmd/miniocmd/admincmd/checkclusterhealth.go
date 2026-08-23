package admincmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/datetime"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewCheckClusterHealthCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-cluster-health",
		Short: "Check the whole cluster health.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			adminClient := options.GetAdminClient(ctx, cmd)

			healthResult := mustutils.Must(nativeminioclient.CheckClusterHealth(ctx, adminClient))

			// Log cluster status
			logging.LogInfoByCtxf(ctx, "Cluster Health Status:")
			logging.LogInfoByCtxf(ctx, "  Total Servers: %d", healthResult.TotalServers)
			logging.LogInfoByCtxf(ctx, "  Online Servers: %d", healthResult.OnlineServers)
			logging.LogInfoByCtxf(ctx, "  Offline Servers: %d", healthResult.OfflineServers)
			logging.LogInfoByCtxf(ctx, "  Total Disks: %d", healthResult.TotalDisks)
			logging.LogInfoByCtxf(ctx, "  Online Disks: %d", healthResult.OnlineDisks)
			logging.LogInfoByCtxf(ctx, "  Offline Disks: %d", healthResult.OfflineDisks)

			// Log all nodes and their status
			logging.LogInfoByCtxf(ctx, "Nodes Status:")
			serverStatusResult := mustutils.Must(nativeminioclient.CheckServerOnlineStatus(ctx, adminClient))
			for _, serverStatus := range serverStatusResult.ServerStatuses {
				status := "OFFLINE"
				if serverStatus.IsOnline {
					status = "ONLINE"
				}
				uptimeStr := mustutils.Must(datetime.FormatDurationAsString(&serverStatus.Uptime))
				logging.LogInfoByCtxf(ctx, "  Node %s: %s (Uptime: %s)", serverStatus.Endpoint, status, uptimeStr)
			}

			// Log all disks and their status
			logging.LogInfoByCtxf(ctx, "Disks Status:")
			diskStatusResults := mustutils.Must(nativeminioclient.CheckDiskStatus(ctx, adminClient))
			for _, serverDiskResult := range diskStatusResults {
				logging.LogInfoByCtxf(ctx, "  Server %s:", serverDiskResult.ServerEndpoint)
				for _, diskStatus := range serverDiskResult.DiskStatuses {
					status := "OFFLINE"
					if diskStatus.IsOnline {
						status = "ONLINE"
					}
					logging.LogInfoByCtxf(ctx, "    Disk %s: %s (State: %s, Used: %.2f%%)",
						diskStatus.DrivePath, status, diskStatus.State, diskStatus.UsedPercent)
				}
			}

			// Show LogGood message when everything is ok and exit 0
			if healthResult.IsHealthy {
				logging.LogGoodByCtxf(ctx, "Cluster health check passed. All servers and disks are online.")
			} else {
				// LogFatal with an error message (this automatically exits != 0)
				errorMessage := "Cluster health check failed."
				if len(healthResult.ServerWarnings) > 0 {
					errorMessage += fmt.Sprintf(" Server issues: %v", healthResult.ServerWarnings)
				}
				if len(healthResult.DiskWarnings) > 0 {
					errorMessage += fmt.Sprintf(" Disk issues: %v", healthResult.DiskWarnings)
				}
				logging.LogFatal(errorMessage)
			}
		},
	}

	return cmd
}
