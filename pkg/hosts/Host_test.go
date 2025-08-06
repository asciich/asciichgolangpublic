package hosts

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestHost_CheckReachable(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				err := host.CheckReachable(verbose)
				require.NoError(t, err)
			},
		)
	}
}

func TestHostGetHostName(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				host := MustGetHostByHostname(tt.hostname)
				require.EqualValues(
					tt.hostname,
					mustutils.Must(host.GetHostName()),
				)
			},
		)
	}
}

func TestHostGetHostDescripion(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)
	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				host := MustGetHostByHostname(tt.hostname)
				require.EqualValues(
					tt.hostname,
					mustutils.Must(host.GetHostDescription()),
				)
			},
		)
	}
}

func TestHostRunCommand(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host := MustGetHostByHostname(tt.hostname)
				ipsString, err := host.RunCommandAndGetStdoutAsString(
					getCtx(),
					&parameteroptions.RunCommandOptions{
						Command: []string{"hostname", "-i"},
					},
				)
				require.NoError(t, err)

				ips := strings.Split(strings.TrimSpace(ipsString), " ")

				require.Contains(t, ips, "192.168.10.32")
			},
		)
	}
}

func TestHost_GetDirectoryByPath(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname       string
		dirPath        string
		expectedExists bool
	}{
		{"localhost", "/home/", true},
		{"localhost", "/home/does_not_exist", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				directory, err := host.GetDirectoryByPath(tt.dirPath)
				require.NoError(t, err)

				_, ok := directory.(*files.CommandExecutorDirectory)
				require.True(t, ok)

				exists, err := directory.Exists(verbose)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedExists, exists)
			},
		)
	}
}

// Connections to the local host should use Bash and not SSH by default.
func TestHost_LocalHostUsesBashCommandExecutorByDefault(t *testing.T) {
	host := MustGetLocalHost()

	commandExecutorHost, ok := host.(*CommandExecutorHost)
	require.True(t, ok)

	commandExecutor := commandExecutorHost.MustGetCommandExecutor()

	_, ok = commandExecutor.(*commandexecutorbashoo.BashService)
	require.True(t, ok)
}
