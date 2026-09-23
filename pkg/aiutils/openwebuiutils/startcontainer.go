package openwebuiutils

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

type StartOpenWebUIContainerOptions struct {
	ContainerName string

	Port int

	// By default it binds only to 127.0.0.1 to be reachable from the local machine.
	// If set to true other machines in the network can reach OpenWebUI as well.
	ReachableByOtherMachines bool

	// Optional: Path to persist OpenWebUI data (litellm.db, uploads, etc.)
	DataPath string

	// Optional: OLLAMA_BASE_URL to connect to an Ollama instance
	OllamaBaseUrl string

	// Optional: Additional environment variables
	AdditionalEnvVars map[string]string
}

func (o *StartOpenWebUIContainerOptions) GetContainerName() (string, error) {
	if o.ContainerName == "" {
		return "", tracederrors.TracedError("ContainerName not set")
	}

	return o.ContainerName, nil
}

func (o *StartOpenWebUIContainerOptions) GetPort() (int, error) {
	if o.Port <= 0 {
		return 0, tracederrors.TracedErrorf("Port '%d' is not set or is invalid.", o.Port)
	}

	return o.Port, nil
}

func (o *StartOpenWebUIContainerOptions) GetReachableByOtherMachines() bool {
	return o.ReachableByOtherMachines
}

func (o *StartOpenWebUIContainerOptions) GetBindAddress() (string, error) {
	port, err := o.GetPort()
	if err != nil {
		return "", err
	}

	prefix := "127.0.0.1"
	if o.ReachableByOtherMachines {
		prefix = "0.0.0.0"
	}

	return prefix + ":" + strconv.Itoa(port), nil
}

func (o *StartOpenWebUIContainerOptions) GetDataPath() (string, error) {
	if o.DataPath == "" {
		return "", tracederrors.TracedError("DataPath not set")
	}

	// For now, just return as-is. Could add absolute path resolution if needed.
	return o.DataPath, nil
}

func (o *StartOpenWebUIContainerOptions) GetOllamaBaseUrl() string {
	return o.OllamaBaseUrl
}

func (o *StartOpenWebUIContainerOptions) GetAdditionalEnvVars() map[string]string {
	if o.AdditionalEnvVars == nil {
		return make(map[string]string)
	}
	return o.AdditionalEnvVars
}

// StartOpenWebUIAsDockerContainer starts an OpenWebUI container with the given options.
// OpenWebUI is a web interface for interacting with LLMs, similar to OpenHands but focused on chat UI.
func StartOpenWebUIAsDockerContainer(ctx context.Context, options *StartOpenWebUIContainerOptions) (containerinterfaces.Container, error) {
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

	logging.LogInfoByCtxf(ctx, "Start OpenWebUI container '%s' bound to port %s started.", containerName, bindAddress)

	// Build environment variables
	envVars := map[string]string{
		"WEBUI_NAME": "OpenWebUI",
	}

	// Add Ollama base URL if provided
	if ollamaUrl := options.GetOllamaBaseUrl(); ollamaUrl != "" {
		envVars["OLLAMA_BASE_URL"] = ollamaUrl
	}

	// Add additional environment variables
	for k, v := range options.GetAdditionalEnvVars() {
		envVars[k] = v
	}

	// Build mount options
	mounts := []string{}

	// Add data persistence mount if DataPath is provided
	if options.DataPath != "" {
		dataPath, err := options.GetDataPath()
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, dataPath+":/app/backend/data")
	}

	container, err := dockerutils.RunContainer(ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:                 containerName,
			ImageName:            "ghcr.io/open-webui/open-webui:latest",
			KeepStoppedContainer: false,
			Ports:                []string{bindAddress + ":8080"},
			AdditionalEnvVars:    envVars,
			Mounts:               mounts,
			SkipIfAlreadyRunning: true,
		},
	)
	if err != nil {
		return nil, err
	}

	url := "http://127.0.0.1:" + strconv.Itoa(port)

	// Wait for OpenWebUI to be ready
	err = httputils.WaitUntilStatusCodeOK(ctx, url, time.Minute*3)
	if err != nil {
		return nil, err
	}

	openwebui, err := NewOpenWebUI(url)
	if err != nil {
		return nil, err
	}

	// Verify the service is healthy
	healthy, err := openwebui.GetHealthStatus(ctx)
	if err != nil || !healthy {
		return nil, tracederrors.TracedErrorf("OpenWebUI health check failed: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Start OpenWebUI container '%s' bound to port %s finished.", containerName, bindAddress)

	return container, nil
}
