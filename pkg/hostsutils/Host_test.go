package hostsutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhost"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/nativehost"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getHostByImplementationName(t *testing.T, implementationName string) hostsutilsinterfaces.Host {
	if implementationName == "commandExecutorHost" {
		host := commandexecutorhost.NewCommandExecutorHost()
		// Set up the bash command executor for localhost testing
		err := host.SetCommandExecutor(commandexecutorbashoo.Bash())
		require.NoError(t, err)
		return host
	}

	if implementationName == "nativeHost" {
		return nativehost.NewNativeHost()
	}

	t.Fatalf("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestHostGetHostName(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		expectedHostname   string
	}{
		{"commandExecutorHost", "localhost"},
		{"nativeHost", "localhost"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host := getHostByImplementationName(t, tt.implementationName)

				hostName, err := host.GetHostName()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedHostname, hostName)
			},
		)
	}
}

func TestHostGetHostDescription(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		expectedHostname   string
	}{
		{"commandExecutorHost", "localhost"},
		{"nativeHost", "localhost"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host := getHostByImplementationName(t, tt.implementationName)

				hostDescription, err := host.GetHostDescription()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedHostname, hostDescription)
			},
		)
	}
}

func TestHost_GetDirectoryByPath(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		dirPath            string
		expectedExists     bool
	}{
		{"commandExecutorHost", "/home/", true},
		{"commandExecutorHost", "/home/does_not_exist", false},
		{"nativeHost", "/home/", true},
		{"nativeHost", "/home/does_not_exist", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				host := getHostByImplementationName(t, tt.implementationName)

				directory, err := host.GetDirectoryByPath(ctx, tt.dirPath)
				require.NoError(t, err)

				exists, err := directory.Exists(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedExists, exists)
			},
		)
	}
}

// Connections to the local host should use Bash and not SSH by default.
func TestHost_LocalHostUsesBashCommandExecutorByDefault(t *testing.T) {
	host, err := hostsutils.GetLocalCommandExecutorHost()
	require.NoError(t, err)

	commandExecutorHost, ok := host.(*commandexecutorhost.CommandExecutorHost)
	require.True(t, ok)

	commandExecutor, err := commandExecutorHost.GetCommandExecutor()
	require.NoError(t, err)

	_, ok = commandExecutor.(*commandexecutorbashoo.BashService)
	require.True(t, ok)
}

func TestHost_LocalHostReturnsNativeHost(t *testing.T) {
	host, err := hostsutils.GetLocalHost()
	require.NoError(t, err)

	_, ok := host.(*nativehost.NativeHost)
	require.True(t, ok)
}

func Test_HostIsCommandExecutor(t *testing.T) {
	var host commandexecutorinterfaces.CommandExecutor
	var err error

	host, err = hostsutils.GetHostByHostname("example.com")
	require.NoError(t, err)
	require.NotNil(t, host)
}

func TestHost_CheckReachable(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorHost"},
		{"nativeHost"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host := getHostByImplementationName(t, tt.implementationName)

				err := host.CheckReachable(true)
				require.NoError(t, err)
			},
		)
	}
}

func TestHost_RunCommand(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		command            []string
		expectedOutput     string
	}{
		{"commandExecutorHost", []string{"echo", "hello"}, "hello"},
		{"nativeHost", []string{"echo", "hello"}, "hello"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				host := getHostByImplementationName(t, tt.implementationName)

				output, err := host.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)
				require.Contains(t, output, tt.expectedOutput)
			},
		)
	}
}

func TestHost_GetCPUArchitecture(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorHost"},
		{"nativeHost"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				host := getHostByImplementationName(t, tt.implementationName)

				arch, err := host.GetCPUArchitecture(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, arch)
			},
		)
	}
}

func TestHost_IsRunningOnLocalhost(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		expectedLocalhost  bool
	}{
		{"commandExecutorHost", true},
		{"nativeHost", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host := getHostByImplementationName(t, tt.implementationName)

				isLocalhost, err := host.IsRunningOnLocalhost()
				require.NoError(t, err)
				require.Equal(t, tt.expectedLocalhost, isLocalhost)
			},
		)
	}
}

func TestHost_GetDeepCopy(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorHost"},
		{"nativeHost"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host := getHostByImplementationName(t, tt.implementationName)

				copy := host.GetDeepCopyAsCommandExecutor()
				require.NotNil(t, copy)

				// Verify the copy has the same host description
				desc, err := copy.GetHostDescription()
				require.NoError(t, err)
				origDesc, err := host.GetHostDescription()
				require.NoError(t, err)
				require.Equal(t, origDesc, desc)
			},
		)
	}
}

func TestHost_GetSshPublicKeyOfUserAsString_ErrorForNonExistentUser(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorHost"},
		{"nativeHost"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				host := getHostByImplementationName(t, tt.implementationName)

				// Try to get SSH key for non-existent user
				_, err := host.GetSshPublicKeyOfUserAsString(ctx, "nonexistent_user_12345")
				require.Error(t, err)
			},
		)
	}
}
