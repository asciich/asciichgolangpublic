package hosts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic"
)

func TestHost_CheckReachable(t *testing.T) {
	if asciichgolangpublic.ContinuousIntegration().IsRunningInContinuousIntegration() {
		asciichgolangpublic.LogInfo("Currently not available in CI/CD pipeline")
		return
	}

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				host.MustCheckReachable(verbose)
			},
		)
	}
}

func TestHostGetHostName(t *testing.T) {
	if asciichgolangpublic.ContinuousIntegration().IsRunningInContinuousIntegration() {
		asciichgolangpublic.LogInfo("Currently not available in CI/CD pipeline")
		return
	}

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				host := MustGetHostByHostname(tt.hostname)
				assert.EqualValues(
					tt.hostname,
					host.MustGetHostName(),
				)
			},
		)
	}
}

func TestHostGetHostDescripion(t *testing.T) {
	if asciichgolangpublic.ContinuousIntegration().IsRunningInContinuousIntegration() {
		asciichgolangpublic.LogInfo("Currently not available in CI/CD pipeline")
		return
	}

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				host := MustGetHostByHostname(tt.hostname)
				assert.EqualValues(
					tt.hostname,
					host.MustGetHostDescription(),
				)
			},
		)
	}
}

func TestHostRunCommand(t *testing.T) {
	if asciichgolangpublic.ContinuousIntegration().IsRunningInContinuousIntegration() {
		asciichgolangpublic.LogInfo("Currently not available in CI/CD pipeline")
		return
	}

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				ipsString := host.MustRunCommandAndGetStdoutAsString(
					&asciichgolangpublic.RunCommandOptions{
						Command: []string{"hostname", "-i"},
						Verbose: verbose,
					},
				)

				ips := strings.Split(strings.TrimSpace(ipsString), " ")

				assert.Contains(
					ips,
					"192.168.10.32",
				)
			},
		)
	}
}

func TestHost_GetDirectoryByPath(t *testing.T) {
	if asciichgolangpublic.ContinuousIntegration().IsRunningInContinuousIntegration() {
		asciichgolangpublic.LogInfo("Currently not available in CI/CD pipeline")
		return
	}

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
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				directory := host.MustGetDirectoryByPath(tt.dirPath)

				_, ok := directory.(*asciichgolangpublic.CommandExecutorDirectory)
				assert.True(ok)

				assert.EqualValues(
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

	_, ok = commandExecutor.(*asciichgolangpublic.BashService)
	require.True(t, ok)
}
