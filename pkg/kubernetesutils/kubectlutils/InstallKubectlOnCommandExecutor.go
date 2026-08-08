package kubectlutils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/commandexecutorinstall"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// InstallKubectlOnCommandExecutor installs kubectl on a target system using the provided CommandExecutor.
// This allows installing kubectl on remote systems or in containers.
//
// Parameters:
//   - ctx: Context for the operation
//   - commandExecutor: CommandExecutor to use for running installation commands
//   - options: Configuration options for the installation
//
// The function:
//  1. Downloads kubectl from the official Kubernetes release URL
//  2. Validates the SHA256 checksum
//  3. Installs to the specified path with configured permissions
//
// Example usage with a Docker container:
//
//	docker, _ := commandexecutordocker.GetLocalCommandExecutorDocker()
//	container, _ := docker.GetContainerByName("my-container")
//	err := InstallKubectlOnCommandExecutor(ctx, container, &InstallKubectlOptions{
//	    InstallPath: "/usr/local/bin/kubectl",
//	    UseSudo:     false,
//	    Version:     "v1.36.2",
//	})
func InstallKubectlOnCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InstallKubectlOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		options = DefaultInstallKubectlOptions()
	}

	// Set defaults for any unset options
	if options.InstallPath == "" {
		options.InstallPath = "/bin/kubectl"
	}
	if options.Version == "" {
		options.Version = "v1.36.2"
	}
	// UseSudo defaults to true if not set

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubectl on '%s' started.", hostDescription)

	sha256Sum := getSha256SumForVersion(options.Version)
	if sha256Sum == "" {
		return tracederrors.TracedErrorf("unsupported kubectl version: %s", options.Version)
	}

	srcUrl := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", options.Version)

	err = commandexecutorinstall.Install(ctx, commandExecutor, &installoptions.InstallOptions{
		SrcUrl:      srcUrl,
		InstallPath: options.InstallPath,
		Mode:        "u=rwx,g=rx,o=rx",
		UseSudo:     options.UseSudo,
		Sha256Sum:   sha256Sum,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubectl on '%s' finished.", hostDescription)

	return nil
}
