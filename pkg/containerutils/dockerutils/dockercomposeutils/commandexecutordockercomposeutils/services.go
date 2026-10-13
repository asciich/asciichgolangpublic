package commandexecutordockercomposeutils

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func checkDockerComposeAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sh", "-c", "docker compose version >/dev/null 2>&1 && echo yes || echo no"},
	})
	if err != nil {
		return tracederrors.TracedErrorf("Unable to check if 'docker compose' is available: %w", err)
	}

	result := strings.TrimSpace(stdout)
	switch result {
	case "yes":
		return nil
	case "no":
		return tracederrors.TracedError("'docker compose' is not available. Please install the Docker Compose v2 plugin (e.g. package 'docker-compose' or 'docker-compose-plugin').")
	default:
		return tracederrors.TracedErrorf("Unexpected output '%s' while checking if 'docker compose' is available.", result)
	}
}

func runCommandAndGetLines(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, command []string) ([]string, error) {
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Running '%s' failed: %w", strings.Join(command, " "), err)
	}

	ret := []string{}
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			ret = append(ret, line)
		}
	}

	return ret, nil
}

func runComposeAndGetLines(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string, args ...string) ([]string, error) {
	command := append([]string{"docker", "compose", "-f", yamlPath}, args...)

	return runCommandAndGetLines(ctx, commandExecutor, command)
}

func getConfiguredServiceNames(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string) ([]string, error) {
	return runComposeAndGetLines(ctx, commandExecutor, yamlPath, "config", "--services")
}

func getRunningServiceNames(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string) ([]string, error) {
	return runComposeAndGetLines(ctx, commandExecutor, yamlPath, "ps", "--services", "--status", "running")
}

func getContainerIDs(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string) ([]string, error) {
	return runComposeAndGetLines(ctx, commandExecutor, yamlPath, "ps", "--all", "--quiet")
}

// getProjectName returns the normalized compose project name of 'yamlPath'.
func getProjectName(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, yamlPath string) (string, error) {
	command := []string{"docker", "compose", "-f", yamlPath, "config", "--format", "json"}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return "", tracederrors.TracedErrorf("Running '%s' failed: %w", strings.Join(command, " "), err)
	}

	var config struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal([]byte(stdout), &config)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to parse docker compose config of '%s' as json: %w", yamlPath, err)
	}

	if config.Name == "" {
		return "", tracederrors.TracedErrorf("Unable to evaluate project name of docker compose file '%s'.", yamlPath)
	}

	return config.Name, nil
}

// getServiceNamesOfAllProjectContainers returns the sorted and unique service names of all containers
// (running or not) labeled with the compose project 'projectName'.
// This also includes containers of services no longer defined in the compose file (orphans).
func getServiceNamesOfAllProjectContainers(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, projectName string) ([]string, error) {
	lines, err := runCommandAndGetLines(ctx, commandExecutor, []string{
		"docker", "ps", "--all",
		"--filter", "label=com.docker.compose.project=" + projectName,
		"--format", `{{.Label "com.docker.compose.service"}}`,
	})
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, line := range lines {
		if !slices.Contains(ret, line) {
			ret = append(ret, line)
		}
	}

	slices.Sort(ret)

	return ret, nil
}
