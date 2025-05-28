package kubernetesutils

import "context"

// Represents a secret in kubernetes.
type Secret interface {
	Exists(ctx context.Context) (bool, error)
}
