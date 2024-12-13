package asciichgolangpublic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostCheckReachableBySsh(t *testing.T) {
	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogInfo("Currently not available in CI/CD pipeline")
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				host.MustCheckReachableBySsh(verbose)
			},
		)
	}
}

func TestHostGetHostname(t *testing.T) {
	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogInfo("Currently not available in CI/CD pipeline")
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				host := MustGetHostByHostname(tt.hostname)
				assert.EqualValues(
					tt.hostname,
					host.MustGetHostname(),
				)
			},
		)
	}
}

func TestHostRunCommand(t *testing.T) {
	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogInfo("Currently not available in CI/CD pipeline")
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				ipsString := host.MustRunCommandAndGetStdoutAsString(
					&RunCommandOptions{
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

func TestHostIsACommandExecutor(t *testing.T) {
	assert := assert.New(t)

	const hostName = "hostname"

	var host CommandExecutor = MustGetHostByHostname(hostName)

	assert.EqualValues(
		hostName,
		host.MustGetHostDescription(),
	)
}

func TestHost_GetDirectoryByPath(t *testing.T) {
	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogInfo("Currently not available in CI/CD pipeline")
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				directory := host.MustGetDirectoryByPath(tt.dirPath)

				_, ok := directory.(*CommandExecutorDirectory)
				assert.True(ok)

				assert.EqualValues(
					tt.expectedExists,
					directory.MustExists(verbose),
				)
			},
		)
	}
}
