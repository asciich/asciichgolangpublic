package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/datetime/durationparser"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ValidateSshKeyInSecretOptions contains options for validating an SSH key stored in a Kubernetes secret.
type ValidateSshKeyInSecretOptions struct {
	// Namespace containing the secret
	Namespace string

	// SecretName is the name of the Kubernetes secret containing the SSH private key
	SecretName string

	// SecretKey is the key name within the secret that contains the private key
	SecretKey string

	// TargetHost is the hostname or IP of the SSH server to test
	TargetHost string

	// TargetUser is the username for SSH authentication
	TargetUser string

	// TargetPort is the SSH port (typically 22)
	TargetPort int

	// SkipHostKeyValidation disables strict host key checking.
	// When true, the SSH client will not verify the server's host key.
	// Useful for testing or when connecting to servers with self-signed keys.
	// Default: false (host key validation is enabled)
	SkipHostKeyValidation bool

	// ConnectionTimeout is the timeout for SSH connection attempt.
	// Default: "10 seconds"
	ConnectionTimeout string

	// ConnectionAttempts is the number of connection attempts.
	// Default: 1
	ConnectionAttempts int
}

// Validate checks if all required fields are set and valid.
func (v *ValidateSshKeyInSecretOptions) Validate() error {
	if v == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if v.Namespace == "" {
		return tracederrors.TracedErrorEmptyString("Namespace")
	}

	if v.SecretName == "" {
		return tracederrors.TracedErrorEmptyString("SecretName")
	}

	if v.SecretKey == "" {
		return tracederrors.TracedErrorEmptyString("SecretKey")
	}

	if v.TargetHost == "" {
		return tracederrors.TracedErrorEmptyString("TargetHost")
	}

	if v.TargetUser == "" {
		return tracederrors.TracedErrorEmptyString("TargetUser")
	}

	if v.TargetPort <= 0 {
		return tracederrors.TracedErrorf("invalid TargetPort: %d", v.TargetPort)
	}

	// Set defaults
	if v.ConnectionTimeout == "" {
		v.ConnectionTimeout = "10 seconds"
	}

	if v.ConnectionAttempts <= 0 {
		v.ConnectionAttempts = 1
	}

	return nil
}

// GetConnectionTimeoutSeconds returns the connection timeout as seconds string for kubectl.
func (v *ValidateSshKeyInSecretOptions) GetConnectionTimeoutSeconds() (string, error) {
	if v.ConnectionTimeout == "" {
		return "10", nil
	}
	return durationparser.ToSecondsAsString(v.ConnectionTimeout)
}
