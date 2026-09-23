# OpenWebUI Utilities

This package provides utilities for working with OpenWebUI, a web interface for interacting with LLMs.

## Features

- Start OpenWebUI as a Docker container
- Query OpenWebUI version and health status
- Configure OpenWebUI with custom settings

## Usage

### Starting OpenWebUI Container

```go
import "gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/aiutils/openwebuiutils"

ctx := context.Background()

options := &openwebuiutils.StartOpenWebUIContainerOptions{
    ContainerName: "openwebui",
    Port:          8080,
    DataPath:      "/path/to/data",  // Optional: for data persistence
    OllamaBaseUrl: "http://localhost:11434",  // Optional: connect to Ollama
}

container, err := openwebuiutils.StartOpenWebUIAsDockerContainer(ctx, options)
if err != nil {
    // handle error
}
```

### Querying OpenWebUI

```go
openwebui, err := openwebuiutils.NewOpenWebUI("http://localhost:8080")
if err != nil {
    // handle error
}

// Get version
version, err := openwebui.GetVersion(ctx)

// Check health
healthy, err := openwebui.GetHealthStatus(ctx)
```

## Configuration Options

- `ContainerName`: Name for the Docker container
- `Port`: Host port to bind OpenWebUI to
- `ReachableByOtherMachines`: If true, binds to 0.0.0.0 instead of 127.0.0.1
- `DataPath`: Path to persist OpenWebUI data (optional)
- `OllamaBaseUrl`: URL of Ollama instance to connect to (optional)
- `AdditionalEnvVars`: Additional environment variables (optional)

## Docker Image

The package uses the official OpenWebUI Docker image: `ghcr.io/open-webui/open-webui:latest`
