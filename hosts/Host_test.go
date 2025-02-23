package hosts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestHost_CheckReachable(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				host.MustCheckReachable(verbose)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				host := MustGetHostByHostname(tt.hostname)
				require.EqualValues(
					tt.hostname,
					host.MustGetHostName(),
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				host := MustGetHostByHostname(tt.hostname)
				require.EqualValues(
					tt.hostname,
					host.MustGetHostDescription(),
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				ipsString := host.MustRunCommandAndGetStdoutAsString(
					&parameteroptions.RunCommandOptions{
						Command: []string{"hostname", "-i"},
						Verbose: verbose,
					},
				)

				ips := strings.Split(strings.TrimSpace(ipsString), " ")

				require.Contains(
					ips,
					"192.168.10.32",
				)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				directory := host.MustGetDirectoryByPath(tt.dirPath)

				_, ok := directory.(*files.CommandExecutorDirectory)
				require.True(ok)

				require.EqualValues(
					tt.expectedExists,
					directory.MustExists(verbose),
				)
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

	_, ok = commandExecutor.(*commandexecutor.BashService)
	require.True(t, ok)
}
