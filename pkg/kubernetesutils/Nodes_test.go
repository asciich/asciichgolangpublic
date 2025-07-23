package kubernetesutils_test

import (
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

				cluster := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				nodeNames, err := cluster.ListNodeNames(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedNodeNames, nodeNames)
			})
	}
}
