package commandexecutordockercomposeutils

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ListOrphanServiceNamesByComposeYamlPath returns the sorted service names of containers which belong
// to the compose project of 'yamlPath' but whose service is no longer defined in the compose file.
//
// Requires the Docker Compose v2 plugin ('docker compose') on the target machine.
// Read-only: No root privileges are required.
func ListOrphanServiceNamesByComposeYamlPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string) ([]string, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if yamlPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("yamlPath")
	}

	err := checkDockerComposeAvailable(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List orphan services of docker compose file '%s' started.", yamlPath)

	configured, err := getConfiguredServiceNames(ctx, commandExecutor, yamlPath)
	if err != nil {
		return nil, err
	}

	projectName, err := getProjectName(ctx, commandExecutor, yamlPath)
	if err != nil {
		return nil, err
	}

	existing, err := getServiceNamesOfAllProjectContainers(ctx, commandExecutor, projectName)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, service := range existing {
		if !slices.Contains(configured, service) {
			ret = append(ret, service)
		}
	}

	logging.LogInfoByCtxf(ctx, "List orphan services of docker compose file '%s' finished. Found '%d' orphan services in project '%s'.", yamlPath, len(ret), projectName)

	return ret, nil
}
