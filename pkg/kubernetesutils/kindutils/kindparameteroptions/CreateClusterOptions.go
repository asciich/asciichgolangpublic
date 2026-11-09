package kindparameteroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

// DefaultNodeImage is the kindest/node image used when no explicit NodeImage is set.
// It matches the kind binary version pinned by kindutils.InstallKind.
// See: https://github.com/kubernetes-sigs/kind/releases/tag/v0.33.0
const DefaultNodeImage = "kindest/node:v1.37.0" // version matching kind v0.33.0

type CreateClusterOptions struct {
	// Name of the cluster to create:
	Name string

	// Numbers of workers to add in the cluster.
	// If unset 0 workers are created which results in a single container deployment.
	Workers int

	// Node image to create the cluster from, e.g. 'kindest/node:v1.27.3'.
	// If unset the `--image` parameter is omitted and the default node image
	// of the installed kind binary is used.
	NodeImage string
}

func (c *CreateClusterOptions) GetNumberOfWorkers() int {
	if c.Workers <= 0 {
		return 0
	}

	return c.Workers
}

func (c *CreateClusterOptions) GetName() (string, error) {
	if c.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return c.Name, nil
}

func (c *CreateClusterOptions) GetNodeImageOrDefaultIfUnset() string {
	if c.NodeImage == "" {
		return DefaultNodeImage
	}

	return c.NodeImage
}

func (c *CreateClusterOptions) GetNodeImageOrEmptyIfUnset() string {
	return c.NodeImage
}
