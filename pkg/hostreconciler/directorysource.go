package hostreconciler

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

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
