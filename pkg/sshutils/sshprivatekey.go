package sshutils

type SSHPrivateKey struct {
	// Type. E.g. "ssh-rsa" or "ssh-ed25519"
	KeyType string

	// The effective key material
	KeyMaterial string
}

