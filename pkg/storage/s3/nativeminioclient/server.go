package nativeminioclient

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/madmin-go/v3"
)

// ServerInfoResult contains the server information returned by ServerInfo
type ServerInfoResult struct {
	Servers []madmin.ServerProperties
}

// ServerOnlineStatus represents the online status of a server
type ServerOnlineStatus struct {
	Endpoint string
	IsOnline bool
	Uptime   time.Duration
}

// ServerOnlineStatusResult contains the results of server online status check
type ServerOnlineStatusResult struct {
	TotalServers   int
	OnlineServers  int
	OfflineServers int
	ServerStatuses []ServerOnlineStatus
}

// DiskStatus represents the status of a single disk
type DiskStatus struct {
	DrivePath   string
	State       string
	IsOnline    bool
	UsedSpace   uint64
	TotalSpace  uint64
	UsedPercent float64
}

// DiskStatusResult contains disk status for a server
type DiskStatusResult struct {
	ServerEndpoint string
	TotalDisks     int
	OnlineDisks    int
	OfflineDisks   int
	DiskStatuses   []DiskStatus
}

// ClusterHealthResult contains the overall cluster health status
type ClusterHealthResult struct {
	IsHealthy      bool
	TotalServers   int
	OnlineServers  int
	OfflineServers int
	TotalDisks     int
	OnlineDisks    int
	OfflineDisks   int
	ServerWarnings []string
	DiskWarnings   []string
}

// GetServerInfo retrieves server information for the entire cluster
func GetServerInfo(ctx context.Context, adminClient *madmin.AdminClient) (*ServerInfoResult, error) {
	info, err := adminClient.ServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	return &ServerInfoResult{
		Servers: info.Servers,
	}, nil
}

// CheckServerOnlineStatus checks the online status of all servers in the cluster
func CheckServerOnlineStatus(ctx context.Context, adminClient *madmin.AdminClient) (*ServerOnlineStatusResult, error) {
	info, err := adminClient.ServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	result := &ServerOnlineStatusResult{
		TotalServers:   len(info.Servers),
		ServerStatuses: make([]ServerOnlineStatus, 0, len(info.Servers)),
	}

	for _, server := range info.Servers {
		isOnline := server.State == "online"
		if isOnline {
			result.OnlineServers++
		} else {
			result.OfflineServers++
		}

		result.ServerStatuses = append(result.ServerStatuses, ServerOnlineStatus{
			Endpoint: server.Endpoint,
			IsOnline: isOnline,
			Uptime:   time.Duration(server.Uptime) * time.Second,
		})
	}

	return result, nil
}

// CheckDiskStatus checks the disk status for all servers in the cluster
func CheckDiskStatus(ctx context.Context, adminClient *madmin.AdminClient) ([]DiskStatusResult, error) {
	info, err := adminClient.ServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	results := make([]DiskStatusResult, 0, len(info.Servers))

	for _, server := range info.Servers {
		serverResult := DiskStatusResult{
			ServerEndpoint: server.Endpoint,
			DiskStatuses:   make([]DiskStatus, 0, len(server.Disks)),
		}

		for _, disk := range server.Disks {
			isOnline := disk.State == "ok"
			if isOnline {
				serverResult.OnlineDisks++
			} else {
				serverResult.OfflineDisks++
			}
			serverResult.TotalDisks++

			usedPercent := 0.0
			if disk.TotalSpace > 0 {
				usedPercent = float64(disk.UsedSpace) / float64(disk.TotalSpace) * 100
			}

			serverResult.DiskStatuses = append(serverResult.DiskStatuses, DiskStatus{
				DrivePath:   disk.DrivePath,
				State:       disk.State,
				IsOnline:    isOnline,
				UsedSpace:   disk.UsedSpace,
				TotalSpace:  disk.TotalSpace,
				UsedPercent: usedPercent,
			})
		}

		results = append(results, serverResult)
	}

	return results, nil
}

// CheckClusterHealth performs a comprehensive health check of the cluster
func CheckClusterHealth(ctx context.Context, adminClient *madmin.AdminClient) (*ClusterHealthResult, error) {
	info, err := adminClient.ServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	result := &ClusterHealthResult{
		TotalServers:   len(info.Servers),
		ServerWarnings: make([]string, 0),
		DiskWarnings:   make([]string, 0),
	}

	// Check server health
	for _, server := range info.Servers {
		if server.State == "online" {
			result.OnlineServers++
		} else {
			result.OfflineServers++
			result.ServerWarnings = append(result.ServerWarnings, fmt.Sprintf("Server %s is offline", server.Endpoint))
		}
	}

	// Check disk health
	for _, server := range info.Servers {
		for _, disk := range server.Disks {
			result.TotalDisks++
			if disk.State == "ok" {
				result.OnlineDisks++
			} else {
				result.OfflineDisks++
				result.DiskWarnings = append(result.DiskWarnings, fmt.Sprintf("Disk %s on server %s is offline", disk.DrivePath, server.Endpoint))
			}
		}
	}

	// Determine overall health
	result.IsHealthy = (result.OfflineServers == 0 && result.OfflineDisks == 0)

	return result, nil
}
