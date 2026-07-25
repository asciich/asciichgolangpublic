package openhandsutils

import (
	"context"
	"strconv"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func StartAsDockerContainer(ctx context.Context, options *StartContainerOptions) (containerinterfaces.Container, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return nil, err
	}

	port, err := options.GetPort()
	if err != nil {
		return nil, err
	}

	bindAddress, err := options.GetBindAddress()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Start OpenHands container '%s' bound to port %s started.", containerName, bindAddress)

	workspacePath, err := options.GetWorkspacePath()
	if err != nil {
		return nil, err
	}

	container, err := dockerutils.RunContainer(ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:                 containerName,
			ImageName:            "ghcr.io/openhands/agent-canvas:latest",
			KeepStoppedContainer: false,
			Ports:                []string{bindAddress + ":8000"},
			AdditionalEnvVars: map[string]string{
				"RUNTIME": "local",
			},
			Mounts:               []string{workspacePath + ":/workspace"},
			SkipIfAlreadyRunning: true,
		},
	)
	if err != nil {
		return nil, err
	}

	err = httputils.WaitUntilStatusCodeOK(ctx, "http://127.0.0.1:"+strconv.Itoa(port), time.Minute*2)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Start OpenHands container '%s' bound to port %s finished.", containerName, bindAddress)

	return container, nil
}
