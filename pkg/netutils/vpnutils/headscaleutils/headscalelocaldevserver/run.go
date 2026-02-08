package headscalelocaldevserver

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscaleoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func RunLocalDevServer(ctx context.Context, options *RunOptions) (headscale headscaleinterfaces.HeadScale, cancel func() error, err error) {
	if options == nil {
		options = &RunOptions{}
	}

	containerName := options.GetContainerName()
	port := options.GetPort()

	logging.LogInfoByCtxf(ctx, "Start local headscale container '%s' on port '%d' started.", containerName, port)

	// Use a minimal config:
	configPath := mustutils.Must(headscalegeneric.WriteMinimalConfigAsTemporaryFile(ctx))

	if options.RestartAlreadyRunningDevServer {
		err := nativedocker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
		if err != nil {
			return nil, nil, err
		}
	}

	container, err := nativedocker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			ImageName: "headscale/headscale",
			Name:      containerName,
			Ports:     []string{fmt.Sprintf("0.0.0.0:%d:8080", options.GetPort())},
			Command:   []string{"serve"},
			Mounts:    []string{configPath + ":/etc/headscale/config.yaml"},
		},
	)
	if err != nil {
		return nil, nil, err
	}

	logging.LogInfoByCtxf(ctx, "Start local headscale container '%s' on port '%d' finished.", containerName, port)

	cancel = func() error {
		logging.LogInfoByCtxf(ctx, "Remove headscale local dev container '%s' and config started.", containerName)
		err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
		if err != nil {
			return err
		}
		err = nativefiles.Delete(ctx, configPath, &filesoptions.DeleteOptions{})
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "Remove headscale local dev container '%s' and config finished.", containerName)

		return nil
	}

	headscale, err = commandexecutorheadscaleoo.New(container)
	if err != nil {
		return nil, nil, err
	}

	return headscale, cancel, nil
}
