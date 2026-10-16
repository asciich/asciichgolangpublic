package hostreconciler

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// DirectorySource reads the desired state of a host from a local config directory.
//
// Expected structure:
//
//	<config-dir>/
//	└── host.yaml
type DirectorySource struct {
	dirPath string
}

// Compile time check to ensure DirectorySource implements the Source interface.
var _ Source = (*DirectorySource)(nil)

// NewDirectorySource creates a new DirectorySource for the given config directory.
//
// The existence of the directory is not checked here on purpose. This allows to create
// the source before the directory is populated (e.g. by a git clone). The check is
// performed when the configuration is read.
func NewDirectorySource(path string) (*DirectorySource, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return &DirectorySource{
		dirPath: path,
	}, nil
}

// GetDirPath returns the path of the config directory.
func (d *DirectorySource) GetDirPath() (string, error) {
	if d.dirPath == "" {
		return "", tracederrors.TracedErrorf("dirPath not set")
	}

	return d.dirPath, nil
}

// GetHostConfigFilePath returns the path of the host.yaml inside the config directory.
func (d *DirectorySource) GetHostConfigFilePath() (string, error) {
	dirPath, err := d.GetDirPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(dirPath, HostConfigFileName), nil
}

// GetHostName returns the desired hostname defined in the host.yaml of the config directory.
//
// Does not require root privileges.
func (d *DirectorySource) GetHostName(ctx context.Context) (string, error) {
	path, err := d.GetHostConfigFilePath()
	if err != nil {
		return "", err
	}

	config, err := d.loadHostConfig(ctx)
	if err != nil {
		return "", err
	}

	ret := config.Hostname
	if ret == "" {
		return "", tracederrors.TracedErrorf("No 'hostname' set in host config '%s'.", path)
	}

	logging.LogInfoByCtxf(ctx, "Desired hostname according to host config '%s' is '%s'.", path, ret)

	return ret, nil
}

// loadHostConfig reads and parses the host.yaml of the config directory.
//
// The parsing is strict:
//   - Unknown fields are rejected to detect typos early.
//   - Only a single YAML document is allowed.
func (d *DirectorySource) loadHostConfig(ctx context.Context) (*hostConfig, error) {
	path, err := d.GetHostConfigFilePath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Load host config '%s' started.", path)

	// TODO: Replace by the corresponding nativefiles function once verified.
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, tracederrors.TracedErrorf("Host config '%s' does not exist: %w", path, err)
		}

		return nil, tracederrors.TracedErrorf("Unable to open host config '%s': %w", path, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	ret := &hostConfig{}
	err = decoder.Decode(ret)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, tracederrors.TracedErrorf("Host config '%s' is empty.", path)
		}

		return nil, tracederrors.TracedErrorf("Unable to parse host config '%s': %w", path, err)
	}

	var additionalDocument yaml.Node
	err = decoder.Decode(&additionalDocument)
	if err == nil {
		return nil, tracederrors.TracedErrorf("Host config '%s' contains more than one YAML document. Only a single document is supported.", path)
	}

	if !errors.Is(err, io.EOF) {
		return nil, tracederrors.TracedErrorf("Unable to parse host config '%s': %w", path, err)
	}

	logging.LogInfoByCtxf(ctx, "Load host config '%s' finished.", path)

	return ret, nil
}

// readHostConfig reads and parses the host.yaml file of the config directory.
func (d *DirectorySource) readHostConfig(ctx context.Context) (*hostConfig, error) {
	if d.dirPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("dirPath")
	}

	configPath := filepath.Join(d.dirPath, HostConfigFileName)

	file, err := os.Open(configPath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to open host config '%s': %w", configPath, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	ret := &hostConfig{}
	err = decoder.Decode(ret)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, tracederrors.TracedErrorf("Unable to parse host config '%s': %w", configPath, err)
	}

	return ret, nil
}

// GetNetworkInterfaces returns the network interfaces defined in 'network.interfaces' of the host.yaml file.
// Returns an empty slice if no network interfaces are configured.
func (d *DirectorySource) GetNetworkInterfaces(ctx context.Context) ([]NetworkInterface, error) {
	config, err := d.readHostConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Return a copy so callers can not modify the parsed config.
	ret := append([]NetworkInterface{}, config.Network.Interfaces...)

	logging.LogInfoByCtxf(ctx, "Read '%d' network interfaces from '%s'.", len(ret), HostConfigFileName)

	return ret, nil
}

// GetBeforeScript returns the entries of 'before_script' of the host.yaml file.
// Returns an empty slice if no before_script is configured.
func (d *DirectorySource) GetBeforeScript(ctx context.Context) ([]string, error) {
	config, err := d.readHostConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Return a copy so callers can not modify the parsed config.
	ret := append([]string{}, config.BeforeScript...)

	logging.LogInfoByCtxf(ctx, "Read '%d' before_script entries from '%s'.", len(ret), HostConfigFileName)

	return ret, nil
}

// GetDockerComposeProjects returns the docker compose projects in the 'docker_compose' directory of the config directory.
//
// Expected structure:
//
//	<config-dir>/
//	└── docker_compose/
//	    └── <deployment>/
//	        └── compose.yaml
//
// Rules:
//   - Every (non hidden) subdirectory of 'docker_compose' is a deployment. Its name is the docker compose project name.
//   - Every deployment directory must contain exactly one of 'compose.yaml', 'compose.yml', 'docker-compose.yaml' or 'docker-compose.yml'.
//   - A top level 'name' in the compose file is optional. If set, it must equal the deployment directory name.
//   - '*.yaml' and '*.yml' files directly inside 'docker_compose' are rejected to avoid silently ignored deployments.
//   - Hidden entries and other files (e.g. 'README.md') directly inside 'docker_compose' are skipped.
//   - The projects are returned in lexical order of their deployment directory names.
//
// Returns an empty slice if the 'docker_compose' directory does not exist.
func (d *DirectorySource) GetDockerComposeProjects(ctx context.Context) ([]DockerComposeProject, error) {
	if d.dirPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("dirPath")
	}

	composeDirPath, err := filepath.Abs(filepath.Join(d.dirPath, DockerComposeDirName))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to get absolute path of docker compose directory in '%s': %w", d.dirPath, err)
	}

	// os.ReadDir returns the entries sorted by file name, which defines the order of the deployments.
	entries, err := os.ReadDir(composeDirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.LogInfoByCtxf(ctx, "Docker compose directory '%s' does not exist. No docker compose projects configured.", composeDirPath)
			return []DockerComposeProject{}, nil
		}

		return nil, tracederrors.TracedErrorf("Unable to list docker compose directory '%s': %w", composeDirPath, err)
	}

	ret := []DockerComposeProject{}
	for _, entry := range entries {
		entryName := entry.Name()

		if isHiddenEntryName(entryName) {
			logging.LogInfoByCtxf(ctx, "Skip hidden entry '%s' in docker compose directory '%s'.", entryName, composeDirPath)
			continue
		}

		if !entry.IsDir() {
			if isYamlFileName(entryName) {
				return nil, tracederrors.TracedErrorf(
					"Docker compose file '%s' found directly in '%s'. Every deployment needs its own directory, e.g. '%s/<deployment>/compose.yaml'.",
					entryName, composeDirPath, DockerComposeDirName,
				)
			}

			logging.LogInfoByCtxf(ctx, "Skip file '%s' in docker compose directory '%s': Only deployment directories are used.", entryName, composeDirPath)
			continue
		}

		project, err := readDockerComposeProject(filepath.Join(composeDirPath, entryName))
		if err != nil {
			return nil, err
		}

		ret = append(ret, project)
	}

	logging.LogInfoByCtxf(ctx, "Read '%d' docker compose projects from '%s'.", len(ret), composeDirPath)

	return ret, nil
}

// readDockerComposeProject reads the docker compose deployment in 'deploymentDirPath'.
func readDockerComposeProject(deploymentDirPath string) (DockerComposeProject, error) {
	if deploymentDirPath == "" {
		return DockerComposeProject{}, tracederrors.TracedErrorEmptyString("deploymentDirPath")
	}

	projectName := filepath.Base(deploymentDirPath)

	err := checkDockerComposeProjectName(projectName)
	if err != nil {
		return DockerComposeProject{}, tracederrors.TracedErrorf("Invalid deployment directory name '%s': %w", deploymentDirPath, err)
	}

	foundPaths := []string{}
	for _, fileName := range dockerComposeFileNames {
		composeFilePath := filepath.Join(deploymentDirPath, fileName)

		info, err := os.Stat(composeFilePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return DockerComposeProject{}, tracederrors.TracedErrorf("Unable to stat compose file '%s': %w", composeFilePath, err)
		}

		if !info.Mode().IsRegular() {
			return DockerComposeProject{}, tracederrors.TracedErrorf("Compose file '%s' is not a regular file.", composeFilePath)
		}

		foundPaths = append(foundPaths, composeFilePath)
	}

	if len(foundPaths) == 0 {
		return DockerComposeProject{}, tracederrors.TracedErrorf(
			"No compose file found in deployment directory '%s'. Expected exactly one of: '%s'.",
			deploymentDirPath, strings.Join(dockerComposeFileNames, "', '"),
		)
	}

	if len(foundPaths) > 1 {
		return DockerComposeProject{}, tracederrors.TracedErrorf(
			"Multiple compose files found in deployment directory '%s': '%s'. Keep exactly one.",
			deploymentDirPath, strings.Join(foundPaths, "', '"),
		)
	}

	composeFilePath := foundPaths[0]

	content, err := os.ReadFile(composeFilePath)
	if err != nil {
		return DockerComposeProject{}, tracederrors.TracedErrorf("Unable to read compose file '%s': %w", composeFilePath, err)
	}

	topLevelName, err := parseDockerComposeTopLevelName(string(content))
	if err != nil {
		return DockerComposeProject{}, tracederrors.TracedErrorf("Invalid compose file '%s': %w", composeFilePath, err)
	}

	if topLevelName != "" && topLevelName != projectName {
		return DockerComposeProject{}, tracederrors.TracedErrorf(
			"Top level 'name' '%s' in compose file '%s' differs from the deployment directory name '%s'. Remove 'name' or set it to '%s'.",
			topLevelName, composeFilePath, projectName, projectName,
		)
	}

	ret := DockerComposeProject{
		Name:            projectName,
		ComposeFilePath: composeFilePath,
	}

	return ret, nil
}

// GetKVMConfig returns the desired KVM configuration read from the 'kvm' section of the host.yaml file.
// Returns a KVMConfig with all actions disabled (and no error) if no kvm section is configured.
func (d *DirectorySource) GetKVMConfig(ctx context.Context) (*KVMConfig, error) {
	config, err := d.readHostConfig(ctx)
	if err != nil {
		return nil, err
	}

	ret := config.KVM

	return &ret, nil
}
