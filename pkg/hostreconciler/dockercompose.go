package hostreconciler

import (
	"context"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockercomposeutils/commandexecutordockercomposeutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockercomposeutils/dockercomposeoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// DockerComposeDirName is the name of the directory inside a config directory
// which contains one subdirectory per docker compose deployment.
const DockerComposeDirName = "docker_compose"

// dockerComposeFileNames are the compose file names supported inside a deployment directory.
// These are the default file names docker compose looks up.
var dockerComposeFileNames = []string{
	"compose.yaml",
	"compose.yml",
	"docker-compose.yaml",
	"docker-compose.yml",
}

// Docker compose project names must only contain lowercase letters, digits, dashes and underscores
// and must start with a lowercase letter or digit.
var dockerComposeProjectNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// DockerComposeProject is a docker compose deployment to apply on the host.
//
// Every subdirectory of 'docker_compose' in the config directory is a deployment:
//
//	docker_compose/<Name>/compose.yaml
//
// It is returned by Source.GetDockerComposeProjects.
type DockerComposeProject struct {
	// Name of the docker compose project. Equals the name of the deployment directory.
	Name string

	// Absolute path to the compose file inside the deployment directory.
	ComposeFilePath string
}

// isHiddenEntryName returns true if 'name' is a hidden file or directory name.
func isHiddenEntryName(name string) bool {
	return strings.HasPrefix(name, ".")
}

// isYamlFileName returns true if 'fileName' has a '.yaml' or '.yml' extension.
func isYamlFileName(fileName string) bool {
	ext := filepath.Ext(fileName)

	return ext == ".yaml" || ext == ".yml"
}

// isDockerComposeFileName returns true if 'fileName' is one of the supported compose file names.
func isDockerComposeFileName(fileName string) bool {
	return slices.Contains(dockerComposeFileNames, fileName)
}

// checkDockerComposeProjectName returns an error if 'projectName' is not a valid docker compose project name.
func checkDockerComposeProjectName(projectName string) error {
	if projectName == "" {
		return tracederrors.TracedErrorEmptyString("projectName")
	}

	if !dockerComposeProjectNameRegex.MatchString(projectName) {
		return tracederrors.TracedErrorf(
			"Invalid docker compose project name '%s': Only lowercase letters, digits, '-' and '_' are allowed and it must start with a lowercase letter or digit.",
			projectName,
		)
	}

	return nil
}

// parseDockerComposeTopLevelName returns the optional top level 'name' of the docker compose file 'content'.
// An empty string is returned if no 'name' is set.
func parseDockerComposeTopLevelName(content string) (string, error) {
	parsed := struct {
		Name string `yaml:"name"`
	}{}

	err := yaml.Unmarshal([]byte(content), &parsed)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to parse docker compose file: %w", err)
	}

	return parsed.Name, nil
}

// validate returns an error if the docker compose project is incomplete or invalid.
func (d *DockerComposeProject) validate() error {
	if d.Name == "" {
		return tracederrors.TracedErrorEmptyString("Name")
	}

	err := checkDockerComposeProjectName(d.Name)
	if err != nil {
		return err
	}

	if d.ComposeFilePath == "" {
		return tracederrors.TracedErrorEmptyString("ComposeFilePath")
	}

	if !filepath.IsAbs(d.ComposeFilePath) {
		return tracederrors.TracedErrorf("Compose file path '%s' of docker compose project '%s' must be absolute.", d.ComposeFilePath, d.Name)
	}

	fileName := filepath.Base(d.ComposeFilePath)
	if !isDockerComposeFileName(fileName) {
		return tracederrors.TracedErrorf(
			"Unsupported compose file name '%s' for docker compose project '%s'. Supported: '%s'.",
			fileName, d.Name, strings.Join(dockerComposeFileNames, "', '"),
		)
	}

	// Docker compose derives the project name from the directory of the compose file.
	deploymentDirName := filepath.Base(filepath.Dir(d.ComposeFilePath))
	if deploymentDirName != d.Name {
		return tracederrors.TracedErrorf(
			"Compose file '%s' must be located in a directory named like the docker compose project '%s', but the directory is named '%s'.",
			d.ComposeFilePath, d.Name, deploymentDirName,
		)
	}

	return nil
}

// validateDockerComposeProjects returns an error if a docker compose project is invalid or the projects conflict with each other.
// An empty list is valid.
func validateDockerComposeProjects(projects []DockerComposeProject) error {
	nameToPath := map[string]string{}

	for i, project := range projects {
		err := project.validate()
		if err != nil {
			return tracederrors.TracedErrorf("Invalid docker compose project entry '%d' ('%s'): %w", i, project.Name, err)
		}

		otherPath, alreadyDefined := nameToPath[project.Name]
		if alreadyDefined {
			return tracederrors.TracedErrorf(
				"Docker compose project '%s' is defined more than once ('%s' and '%s').",
				project.Name, otherPath, project.ComposeFilePath,
			)
		}

		nameToPath[project.Name] = project.ComposeFilePath
	}

	return nil
}

// reconcileDockerComposeProjects ensures all services of every docker compose project in 'projects' are up on localhost.
//
// The projects are applied one by one in the given order using commandexecutordockercomposeutils.UpByComposeYamlPath
// with 'RemoveOrphans' set: Services removed from a compose file are removed from the host as well.
// This is safe since every project has its own unique name (the deployment directory name).
//
// Reconciliation is additive on project level: Removing a whole deployment directory does NOT stop its project.
//
// All projects are validated before any change is performed.
//
// Requires the Docker Compose v2 plugin and access to the docker daemon.
// This function is idempotent.
func reconcileDockerComposeProjects(ctx context.Context, projects []DockerComposeProject) error {
	err := validateDockerComposeProjects(projects)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		logging.LogInfoByCtxf(ctx, "No docker compose projects configured. Nothing to reconcile for docker compose.")
		return nil
	}

	logging.LogInfoByCtxf(ctx, "Reconcile '%d' docker compose projects started.", len(projects))

	executor := commandexecutorexecoo.Exec()

	for _, project := range projects {
		logging.LogInfoByCtxf(ctx, "Reconcile docker compose project '%s' from '%s'.", project.Name, project.ComposeFilePath)

		err = commandexecutordockercomposeutils.UpByComposeYamlPath(
			ctx,
			executor,
			project.ComposeFilePath,
			&dockercomposeoptions.UpOptions{
				RemoveOrphans: true,
			},
		)
		if err != nil {
			return tracederrors.TracedErrorf("Reconcile docker compose project '%s' from '%s' failed: %w", project.Name, project.ComposeFilePath, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Reconcile '%d' docker compose projects finished.", len(projects))

	return nil
}
