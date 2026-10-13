package commandexecutordockercomposeutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockercomposeutils/dockercomposeoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// DownByComposeYamlPath ensures all containers of the compose file 'yamlPath' are stopped and removed
// on the machine reached by 'commandExecutor'. This is the idempotent equivalent of:
//
//	docker compose -f <yamlPath> down
//
// If no containers exist nothing is changed.
//
// Requires the Docker Compose v2 plugin ('docker compose') on the target machine.
// No root privileges are required as long as the executing user has access to the docker daemon.
func DownByComposeYamlPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string, options *dockercomposeoptions.DownOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if yamlPath == "" {
		return tracederrors.TracedErrorEmptyString("yamlPath")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	err := checkDockerComposeAvailable(ctx, commandExecutor)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Docker compose down for '%s' started.", yamlPath)

	containerIDs, err := getContainerIDs(ctx, commandExecutor, yamlPath)
	if err != nil {
		return err
	}

	if len(containerIDs) == 0 {
		logging.LogInfoByCtxf(ctx, "Docker compose file '%s' has no containers. Already down.", yamlPath)
	} else {
		_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"docker", "compose", "-f", yamlPath, "down"},
		})
		if err != nil {
			return tracederrors.TracedErrorf("Docker compose down for '%s' failed: %w", yamlPath, err)
		}

		logging.LogChangedByCtxf(ctx, "Removed '%d' containers of docker compose file '%s'.", len(containerIDs), yamlPath)
	}

	logging.LogInfoByCtxf(ctx, "Docker compose down for '%s' finished.", yamlPath)

	return nil
}
