package kvmutilsinterfaces

import "context"

type Volume interface {
	Delete(ctx context.Context) error

	GetName() (string, error)

	SetStoragePool(storagePool StoragePool) error
}
