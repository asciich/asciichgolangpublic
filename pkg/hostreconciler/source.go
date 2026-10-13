package hostreconciler

import "context"

// Source provides the desired state of a host.
//
// Implementations:
//   - DirectorySource: Reads the desired state from a local directory.
type Source interface {
	// GetHostName returns the desired hostname of the host.
	GetHostName(ctx context.Context) (string, error)

	// GetNetworkInterfaces returns the desired network interface configurations of the host.
	// Returns an empty slice (and no error) if no network interfaces are configured.
	GetNetworkInterfaces(ctx context.Context) ([]NetworkInterface, error)
}
