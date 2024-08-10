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

				const verbose = true

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
