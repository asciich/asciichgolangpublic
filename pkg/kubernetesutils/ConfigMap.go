package kubernetesutils

import "context"

type ConfigMap interface {
	Exists(ctx context.Context) (bool, error)
	GetAllData(ctx context.Context) (map[string]string, error)
	GetAllLabels(ctx context.Context) (map[string]string, error)
	GetData(ctx context.Context, fieldName string) (string, error)
}
