package kubernetesutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_PodsRunSingleCommand_echoHelloWorld(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		// {"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const podName = "podname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Wait until default sa in created namespace exists.
				time.Sleep(10 * time.Second)

				output, err := kubernetes.RunCommandInTemporaryPod(
					ctx,
					&kubernetesparameteroptions.RunCommandOptions{
						Image:                    "ubuntu",
						Namespace:                namespaceName,
						PodName:                  podName,
						DeleteAlreadyExistingPod: true,
						Command:                  []string{"bash", "-c", "echo hello_world"},
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, "hello_world\n", stdout)

				stderr, err := output.GetStderrAsString()
				require.NoError(t, err)
				require.EqualValues(t, "", stderr)

				retVal, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, 0, retVal)
			},
		)
	}
}

