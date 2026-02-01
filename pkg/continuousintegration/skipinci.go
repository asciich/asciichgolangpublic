package continuousintegration

import "testing"

// Skips the test if running in CI.
func SkipInCi(t *testing.T, msg string) {
	if IsRunningInContinuousIntegration() {
		t.Skip(msg)
	}
}

// Skips the test if running in Github CI.
func SkipInGithubCi(t *testing.T, msg string) {
	if IsRunningInGithub() {
		t.Skip(msg)
	}
}
