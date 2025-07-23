package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_WhoAmI(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				userInfo, err := kubernetes.WhoAmI(ctx)
				require.NoError(t, err)
				require.NotNil(t, userInfo)
				require.EqualValues(t, "kubernetes-admin", userInfo.Username)
			},
		)
	}
}
