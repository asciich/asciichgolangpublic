package fluxinterfaces

import "context"

// Represents flux deployed in a kubernetes cluster.
type FluxDeployment interface {
	GitRepositoryExists(ctx context.Context, name string, namespace string) (bool, error)
}
