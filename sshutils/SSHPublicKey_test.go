package sshutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestSshPublicKeySetFromString(t *testing.T) {
	tests := []struct {
		keyMaterial         string
		expectedKeyMaterial string
	}{
		{"AAAabc", "AAAabc"},
		{"ssh-rsa AAAabc user@host", "AAAabc"},
		{"ssh-rsa AAAabc", "AAAabc"},
		{"AAAabc user@host", "AAAabc"},
		{"AAAabc\n", "AAAabc"},
		{"ssh-rsa AAAabc user@host\n", "AAAabc"},
		{"ssh-rsa AAAabc\n", "AAAabc"},
		{"AAAabc user@host\n", "AAAabc"},
		{"\nAAAabc", "AAAabc"},
		{"\nssh-rsa AAAabc user@host", "AAAabc"},
		{"\nssh-rsa AAAabc", "AAAabc"},
		{"\nAAAabc user@host", "AAAabc"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				sshPublicKey := MustLoadPublicKeyFromString(tt.keyMaterial)

				keyMaterial := sshPublicKey.MustGetKeyMaterialAsString()
				require.EqualValues(tt.expectedKeyMaterial, keyMaterial)
			},
		)
	}
}

func TestSshPublicKeySetFromStringUserCorrect(t *testing.T) {
	tests := []struct {
		keyMaterial         string
		expectedKeyMaterial string
		expectedUserName    string
		expectedUserHost    string
	}{
		{"ssh-rsa AAAabc user@host", "AAAabc", "user", "host"},
		{"ssh-rsa AAAabc user@host\n", "AAAabc", "user", "host"},
		{"AAAabc user@host\n", "AAAabc", "user", "host"},
		{"\nssh-rsa AAAabc user@host", "AAAabc", "user", "host"},
		{"\nAAAabc user@host", "AAAabc", "user", "host"},
		{"\nAAAabc user2@host", "AAAabc", "user2", "host"},
		{"\nAAAabc user2@host3", "AAAabc", "user2", "host3"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				sshPublicKey := MustLoadPublicKeyFromString(tt.keyMaterial)

				require.EqualValues(tt.expectedKeyMaterial, sshPublicKey.MustGetKeyMaterialAsString())
				require.EqualValues(tt.expectedUserName, sshPublicKey.MustGetKeyUserName())
				require.EqualValues(tt.expectedUserHost, sshPublicKey.MustGetKeyHostName())
			},
		)
	}
}

func Test_SetFromString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		key := NewSSHPublicKey()
		require.Error(t, key.SetFromString(""))
	})

	t.Run("ed25519", func(t *testing.T) {
		key := NewSSHPublicKey()
		require.NoError(t, key.SetFromString("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEB7W3jJgHEzL4kteQ4MlLPosP2zaqRRKEydm7ic5HKN user@host1234"))
		require.EqualValues(t, "user", key.MustGetKeyUserName())
		require.EqualValues(t, "host1234", key.MustGetKeyUserHost())
	})
}
