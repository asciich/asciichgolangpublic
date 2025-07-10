package kubernetesinterfaces

import "context"

// Represents a secret in kubernetes.
type Secret interface {
	Exists(ctx context.Context) (bool, error)
	Read(ctx context.Context) (map[string][]byte, error)
}
