package containerinterfaces

import "context"

type Image interface {
	Exists(ctx context.Context) (bool, error)

	// Get the name of the image.
	GetName() (string, error)
}
