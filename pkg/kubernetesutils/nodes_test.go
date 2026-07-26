package kubernetesutils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_ListNodeNames(t *testing.T) {
	tests := []struct {
		implementationName string
		expectedNodeNames  []string
	}{
		{"commandExecutorKubernetes", []string{"kubernetesutils-control-plane"}},
		{"nativeKubernetes", []string{"kubernetesutils-control-plane"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				cluster := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				nodeNames, err := cluster.ListNodeNames(ctx)
				require.NoError(t, err)
				require.Len(t, nodeNames, 1)

				// We can't check for the full name since there are multiple KinD cluster running in parallel.
				// This ends in generated names like e.g. "kubernetesutils-vgkda-control-plane".
				nodeName := nodeNames[0]
				require.True(t, strings.HasPrefix(nodeName, "kubernetesutils-"))
				require.True(t, strings.HasSuffix(nodeName, "-control-plane"))
			})
	}
}
