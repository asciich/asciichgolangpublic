package kindparameteroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

// DefaultNodeImage is the kindest/node image used when no explicit NodeImage is set.
// It matches the kind binary version pinned by kindutils.InstallKind.
// See: https://github.com/kubernetes-sigs/kind/releases/tag/v0.20.0
const DefaultNodeImage = "kindest/node:v1.27.3"

type CreateClusterOptions struct {
	// Name of the cluster to create:
	Name string

	// Numbers of workers to add in the cluster.
	// If unset 0 workers are created which results in a single container deployment.
	Workers int

	// Node image to create the cluster from, e.g. 'kindest/node:v1.27.3'.
	// If unset DefaultNodeImage is used.
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

func (c *CreateClusterOptions) GetNodeImage() string {
	if c.NodeImage == "" {
		return DefaultNodeImage
	}

	return c.NodeImage
}
