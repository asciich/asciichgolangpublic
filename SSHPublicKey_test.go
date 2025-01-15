package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert := assert.New(t)

				sshPublicKey := new(SSHPublicKey)
				sshPublicKey.MustSetFromString(tt.keyMaterial)

				keyMaterial := sshPublicKey.MustGetKeyMaterialAsString()
				assert.EqualValues(tt.expectedKeyMaterial, keyMaterial)
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
				assert := assert.New(t)

				sshPublicKey := new(SSHPublicKey)
				sshPublicKey.MustSetFromString(tt.keyMaterial)

				assert.EqualValues(tt.expectedKeyMaterial, sshPublicKey.MustGetKeyMaterialAsString())
				assert.EqualValues(tt.expectedUserName, sshPublicKey.MustGetKeyUserName())
				assert.EqualValues(tt.expectedUserHost, sshPublicKey.MustGetKeyHostName())
			},
		)
	}
}
