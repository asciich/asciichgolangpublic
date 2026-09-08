package hostsutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestHost_IsReachable(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		expectedReachable  bool
	}{
		{"commandExecutorHost", true},
		{"nativeHost", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				host := getHostByImplementationName(t, tt.implementationName)

				isReachable, err := host.IsReachable(ctx)
				require.NoError(t, err)
				require.Equal(t, tt.expectedReachable, isReachable)
			},
		)
	}
}

func TestHost_WaitUntilReachable(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		renewHostKey       bool
	}{
		{"commandExecutorHost", false},
		{"commandExecutorHost", true},
		{"nativeHost", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				host := getHostByImplementationName(t, tt.implementationName)

				// localhost must always be reachable, so this must return without error.
				err := host.WaitUntilReachable(ctx, tt.renewHostKey)
				require.NoError(t, err)
			},
		)
	}
}
