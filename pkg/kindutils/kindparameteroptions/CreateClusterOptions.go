package kindparameteroptions

import "gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"

type CreateClusterOptions struct {
	// Name of the cluster to create:
	Name string

	// Numbers of workers to add in the cluster.
	// If unset 0 workers are created which results in a single container deployment.
	Workers int
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
