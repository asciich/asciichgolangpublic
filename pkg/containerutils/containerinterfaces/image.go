package containerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
)

type Image interface {
	Exists(ctx context.Context) (bool, error)

	// Get the name of the image.
	GetName() (string, error)

	// Remove removes the image.
	Remove(ctx context.Context, options *dockeroptions.RemoveOptions) error
}
