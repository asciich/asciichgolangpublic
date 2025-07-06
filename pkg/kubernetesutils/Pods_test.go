package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func Test_PodsRunSingleCommand(t *testing.T) {
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

			},
		)
	}
}
