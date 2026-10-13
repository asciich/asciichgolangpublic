package commandexecutordockercomposeutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockercomposeutils/dockercomposeoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// UpByComposeYamlPath ensures all services defined in the compose file 'yamlPath' are running
// on the machine reached by 'commandExecutor'. This is the idempotent equivalent of:
//
//	docker compose -f <yamlPath> up --detach [--remove-orphans]
//
// If 'options.RemoveOrphans' is set, containers of services no longer defined in the compose file are removed.
// If all services are already running (and no orphans need to be removed) nothing is changed.
//
// Requires the Docker Compose v2 plugin ('docker compose') on the target machine.
// No root privileges are required as long as the executing user has access to the docker daemon.
func UpByComposeYamlPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string, options *dockercomposeoptions.UpOptions) error {
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

	logging.LogInfoByCtxf(ctx, "Docker compose up for '%s' started.", yamlPath)

	isUp, err := IsUpByComposeYamlPath(ctx, commandExecutor, yamlPath)
	if err != nil {
		return err
	}

	orphans, err := ListOrphanServiceNamesByComposeYamlPath(ctx, commandExecutor, yamlPath)
	if err != nil {
		return err
	}

	removeOrphans := options.RemoveOrphans && len(orphans) > 0

	if !options.RemoveOrphans && len(orphans) > 0 {
		logging.LogInfoByCtxf(ctx, "Docker compose file '%s' has orphan services '%s' which are kept since RemoveOrphans is not set.", yamlPath, strings.Join(orphans, ", "))
	}

	if isUp && !removeOrphans {
		logging.LogInfoByCtxf(ctx, "All services of docker compose file '%s' are already up.", yamlPath)
	} else {
		command := []string{"docker", "compose", "-f", yamlPath, "up", "--detach"}
		if options.RemoveOrphans {
			command = append(command, "--remove-orphans")
		}

		_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: command,
		})
		if err != nil {
			return tracederrors.TracedErrorf("Docker compose up for '%s' failed: %w", yamlPath, err)
		}

		if !isUp {
			logging.LogChangedByCtxf(ctx, "Services of docker compose file '%s' started.", yamlPath)
		}

		if removeOrphans {
			logging.LogChangedByCtxf(ctx, "Removed orphan services '%s' of docker compose file '%s'.", strings.Join(orphans, ", "), yamlPath)
		}
	}

	logging.LogInfoByCtxf(ctx, "Docker compose up for '%s' finished.", yamlPath)

	return nil
}
