package commandexecutordockercomposeutils

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsUpByComposeYamlPath returns true if all services defined in the compose file 'yamlPath'
// are running on the machine reached by 'commandExecutor'.
//
// Requires the Docker Compose v2 plugin ('docker compose') on the target machine.
// Read-only: No root privileges are required.
func IsUpByComposeYamlPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if yamlPath == "" {
		return false, tracederrors.TracedErrorEmptyString("yamlPath")
	}

	err := checkDockerComposeAvailable(ctx, commandExecutor)
	if err != nil {
		return false, err
	}

	configured, err := getConfiguredServiceNames(ctx, commandExecutor, yamlPath)
	if err != nil {
		return false, err
	}

	if len(configured) == 0 {
		return false, tracederrors.TracedErrorf("No services defined in docker compose file '%s'.", yamlPath)
	}

	running, err := getRunningServiceNames(ctx, commandExecutor, yamlPath)
	if err != nil {
		return false, err
	}

	ret := true
	for _, service := range configured {
		if !slices.Contains(running, service) {
			ret = false
			break
		}
	}

	logging.LogInfoByCtxf(ctx, "Docker compose file '%s': '%d' of '%d' services running, is up: '%t'.", yamlPath, len(running), len(configured), ret)

	return ret, nil
}
