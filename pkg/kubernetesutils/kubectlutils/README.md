# kubectlutils

This package handles the binary `kubectl` itself, not kubernetes (k8s).

## Functions

### InstallKubectl

Installs the kubectl binary from the official Kubernetes release URL.

**Features:**
- Downloads kubectl from https://dl.k8s.io
- Configurable install path (default: `/bin/kubectl`)
- Configurable version (default: `v1.36.2`)
- Configurable sudo usage (default: `true`)
- Sets permissions to `u=rwx,g=rx,o=rx`
- Validates SHA256 checksum

**Options:**
```go
type InstallKubectlOptions struct {
    InstallPath string  // Path where kubectl will be installed (default: /bin/kubectl)
    UseSudo     bool    // Use sudo for installation (default: true)
    Version     string  // kubectl version to install (default: v1.36.2)
}
```

**Requirements:**
- Network access to reach the Kubernetes release server
- Sudo privileges if installing to system directories (unless `UseSudo: false`)

**Example:**
```go
// Install with defaults
ctx := contextutils.WithVerbose(context.TODO())
err := kubectlutils.InstallKubectl(ctx, kubectlutils.DefaultInstallKubectlOptions())
if err != nil {
    panic(err)
}

// Install to custom path without sudo
options := &kubectlutils.InstallKubectlOptions{
    InstallPath: "~/.local/bin/kubectl",
    UseSudo:     false,
    Version:     "v1.36.2",
}
err := kubectlutils.InstallKubectl(ctx, options)
```

### InstallKubectlOnCommandExecutor

Installs kubectl on a target system using a CommandExecutor. This allows installing kubectl on remote systems or in Docker containers.

**Features:**
- Same features as `InstallKubectl`
- Works with any CommandExecutor implementation
- Supports Docker containers, SSH hosts, and other remote targets

**Example with Docker container:**
```go
// Get Docker executor
docker, _ := commandexecutordocker.GetLocalCommandExecutorDocker()

// Start a container
container, _ := docker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
    Name:      "my-container",
    ImageName: "alpine:latest",
    Command:   []string{"sleep", "300"},
})

// Install kubectl inside the container
err := kubectlutils.InstallKubectlOnCommandExecutor(ctx, container, &kubectlutils.InstallKubectlOptions{
    InstallPath: "/usr/local/bin/kubectl",
    UseSudo:     false,
    Version:     "v1.36.2",
})
```

## Testing

Run tests with:
```bash
go test -v ./...
```

Note: Tests require network access. The tests use temporary file paths to avoid overwriting system binaries.
Tests using Docker require Docker to be installed and running.

## Specifications

For specifications see [kubectlutils.spec.md](kubectlutils.spec.md)

## Examples

* [Install kubectl](Example_InstallKubectl_test.go)
