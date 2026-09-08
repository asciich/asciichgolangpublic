package kubectlutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
)

func TestInstallKubectl_customPath(t *testing.T) {
	// Create a context with verbose logging
	ctx := contextutils.WithVerbose(context.TODO())

	// Install kubectl to a custom path without sudo
	// Useful for user-local installations or testing
	options := &kubectlutils.InstallKubectlOptions{
		InstallPath: "~/.local/bin/kubectl2", // Install to user's local bin. for this example it's on purporse installed as kubectl2 to not provoke race conditions in CI.
		UseSudo:     false,                   // No sudo required
		Version:     "v1.36.2",               // Specify version
	}

	err := kubectlutils.InstallKubectl(ctx, options)
	require.NoError(t, err)

	// kubectl is now installed at ~/.local/bin/kubectl2
}
