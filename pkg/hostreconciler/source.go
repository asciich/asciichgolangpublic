package hostreconciler

import "context"

// Source provides the desired state of a host.
//
// Implementations:
//   - DirectorySource: Reads the desired state from a local directory.
type Source interface {
	// GetHostName returns the desired hostname of the host.
	GetHostName(ctx context.Context) (string, error)
}
