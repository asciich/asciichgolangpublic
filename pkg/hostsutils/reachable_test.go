package hostsutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhost"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestHost_IsReachable(t *testing.T) {
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
				err := host.WaitUntilReachable(ctx, &hostsutilsoptions.WaitUntilReachableOptions{
					RenewHostKey: tt.renewHostKey,
				})

				require.NoError(t, err)
			},
		)
	}
}

func TestHost_IsReachable_TimeoutForUnreachableHost(t *testing.T) {
	// 10.255.255.1 is in a non-routable range that typically black-holes
	// packets, so the SSH connection attempt hangs instead of failing fast.
	// This lets us verify the 5 second timeout actually kicks in.
	host, err := hostsutils.GetHostByHostname("10.255.255.1")
	require.NoError(t, err)

	// Ensure we exercise the SSH command executor (not bash/localhost).
	_, ok := host.(*commandexecutorhost.CommandExecutorHost)
	require.True(t, ok)

	ctx := getCtx()

	tStart := time.Now()
	isReachable, err := host.IsReachable(ctx)
	elapsed := time.Since(tStart)

	// The host is unreachable, so IsReachable must not report it as reachable.
	require.Error(t, err)
	require.False(t, isReachable)

	// The 5 second timeout must kick in: the call must not hang much longer
	// than the configured 5 seconds, and it must not return early either.
	require.GreaterOrEqual(t, elapsed, 5*time.Second)
	require.Less(t, elapsed, 10*time.Second)
}
