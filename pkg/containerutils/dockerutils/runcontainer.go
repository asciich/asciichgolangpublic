package dockerutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// This is a convenience function to directly start a docker container on localhost.
func RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return nil, err
	}

	return docker.RunContainer(ctx, options)
}
