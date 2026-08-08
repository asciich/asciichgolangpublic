package kubectlutils_test

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
)

func ExampleInstallKubectl_customPath() {
	// Create a context with verbose logging
	ctx := contextutils.WithVerbose(context.TODO())

	// Install kubectl to a custom path without sudo
	// Useful for user-local installations or testing
	options := &kubectlutils.InstallKubectlOptions{
		InstallPath: "~/.local/bin/kubectl", // Install to user's local bin
		UseSudo:     false,                  // No sudo required
		Version:     "v1.36.2",              // Specify version
	}

	err := kubectlutils.InstallKubectl(ctx, options)
	if err != nil {
		panic(err)
	}

	// kubectl is now installed at ~/.local/bin/kubectl
}
