package dockerutils

import (
"context"

"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// BuildContainerImage is a convenience function to build a Docker image on localhost.
func BuildContainerImage(ctx context.Context, options *dockeroptions.BuildContainerOptions) (containerinterfaces.Image, error) {
if options == nil {
return nil, tracederrors.TracedErrorNil("options")
}

docker, err := GetDockerOnLocalHost()
if err != nil {
return nil, err
}

return docker.BuildImage(ctx, options)
}
