package fluxinterfaces

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Represents flux deployed in a kubernetes cluster.
type FluxDeployment interface {
	DeleteGitRepository(ctx context.Context, name string, namespace string) error
	DeleteKustomization(ctx context.Context, name string, namespace string) error
	GetGitRepositoryStatusMessage(ctx context.Context, name string, namespace string) (string, error)
	GetKustomizationStatusMessage(ctx context.Context, name string, namespace string) (string, error)
	GitRepositoryExists(ctx context.Context, name string, namespace string) (bool, error)
	KustomizationExists(ctx context.Context, name string, namespace string) (bool, error)
	WatchGitRepository(ctx context.Context, name string, namespace string, create func(*unstructured.Unstructured), update func(*unstructured.Unstructured), delete func(*unstructured.Unstructured)) error
	WatchKustomization(ctx context.Context, name string, namespace string, create func(*unstructured.Unstructured), update func(*unstructured.Unstructured), delete func(*unstructured.Unstructured)) error
}
