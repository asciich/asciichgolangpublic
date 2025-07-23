package sshutils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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

				sshPublicKey := sshutils.MustLoadPublicKeyFromString(tt.keyMaterial)

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

				sshPublicKey := sshutils.MustLoadPublicKeyFromString(tt.keyMaterial)

				require.EqualValues(tt.expectedKeyMaterial, sshPublicKey.MustGetKeyMaterialAsString())
				require.EqualValues(tt.expectedUserName, sshPublicKey.MustGetKeyUserName())
				require.EqualValues(tt.expectedUserHost, sshPublicKey.MustGetKeyHostName())
			},
		)
	}
}

func Test_SetFromString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		key := sshutils.NewSSHPublicKey()
		require.Error(t, key.SetFromString(""))
	})

	t.Run("ed25519", func(t *testing.T) {
		key := sshutils.NewSSHPublicKey()
		input := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEB7W3jJgHEzL4kteQ4MlLPosP2zaqRRKEydm7ic5HKN user@host1234"
		require.NoError(t, key.SetFromString(input))
		require.EqualValues(t, "user", key.MustGetKeyUserName())
		require.EqualValues(t, "host1234", mustutils.Must(key.GetKeyUserHost()))
		require.EqualValues(t, "ssh-ed25519", mustutils.Must(key.GetKeyType()))
		require.EqualValues(t, input, mustutils.Must(key.GetAsPublicKeyLine()))
	})

	t.Run("rsa", func(t *testing.T) {
		key := sshutils.NewSSHPublicKey()
		input := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDgjxbTko4CIj7UZXpztBthlkwMV528uIreAVC9WiI0fX6g4QQ+EoTgcLAItMMZprBztkMxoWcu+YNIEPR0SwF8vvyct5ENVkCNOEa1fRfYct8u6ETQxdewmiUlfnIECD3j0c2REny9GYs9qgjUa+MOBfgCDajUTWB37S0cooaEK8Dz6sla/ESFCe1i+c2NKoFzLRoMvh5Oty45TXHa+QG/YDn2PfZxKTdZyXCYDpXJb4QDZFiOkiz+HTFpk9F8+RrwSOne7DI+wSx+VWMtxCF0t1tfiNZo+7DpR63jWtBIBp9Xt4ztGekX2D7ufhZ/XPiuhaXD1E17GebyGLsL8GLJgM/k1jX4UIvVueIgFjFIjuzPyjIr/1KPAUittGd/VyRT7R06UrKXvRHuPTxAHpwUwtDWnVEhlbm8+nxZaviOCUOKan9pT3TF8Uoay1CltPG49aQm6iVbdke/N4h+JgxDirjzJlqVyfvu3X3dven3ibcgiJx9fmhPe6iOC8k3CQ0= user2@host2345"
		require.NoError(t, key.SetFromString(input))
		require.EqualValues(t, "user2", key.MustGetKeyUserName())
		require.EqualValues(t, "host2345", mustutils.Must(key.GetKeyUserHost()))
		require.EqualValues(t, "ssh-rsa", mustutils.Must(key.GetKeyType()))
		require.EqualValues(t, input, mustutils.Must(key.GetAsPublicKeyLine()))
	})

	t.Run("ecdsa", func(t *testing.T) {
		key := sshutils.NewSSHPublicKey()
		input := "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBOCP6lN77PFHaOuTA7qt4cV8Fn/5ZER0Ufino987ObNKiYzwRX5lZ5MrxseMf3+QjH0g1XMLtqREdl888OUPovU= user@host"
		require.NoError(t, key.SetFromString(input))
		require.EqualValues(t, "user", key.MustGetKeyUserName())
		require.EqualValues(t, "host", mustutils.Must(key.GetKeyUserHost()))
		require.EqualValues(t, "ecdsa-sha2-nistp256", mustutils.Must(key.GetKeyType()))
		require.EqualValues(t, input, mustutils.Must(key.GetAsPublicKeyLine()))
	})
}

func Test_Equals(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		require.False(t, key1.Equals(nil))
	})
	t.Run("empty equals", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key2 := sshutils.NewSSHPublicKey()
		require.True(t, key1.Equals(key2))
	})

	t.Run("keyMaterial differ", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyMaterial = "AAAabc"
		key2 := sshutils.NewSSHPublicKey()
		require.False(t, key1.Equals(key2))
	})

	t.Run("keyUserName differ", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyMaterial = "username"
		key2 := sshutils.NewSSHPublicKey()
		require.False(t, key1.Equals(key2))
	})

	t.Run("keyUserHost differ", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyUserName = "host"
		key2 := sshutils.NewSSHPublicKey()
		require.False(t, key1.Equals(key2))
	})

	t.Run("Only key material set and equal", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyMaterial = "AAAabc"
		key2 := sshutils.NewSSHPublicKey()
		key2.KeyMaterial = "AAAabc"
		require.True(t, key1.Equals(key2))
	})

	t.Run("Key materual and user set equal", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyMaterial = "AAAabc"
		key1.KeyUserHost = "user"
		key2 := sshutils.NewSSHPublicKey()
		key2.KeyMaterial = "AAAabc"
		key2.KeyUserHost = "user"
		require.True(t, key1.Equals(key2))
	})

	t.Run("Key materual and host set equal", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyMaterial = "AAAabc"
		key1.KeyUserHost = "host"
		key2 := sshutils.NewSSHPublicKey()
		key2.KeyMaterial = "AAAabc"
		key2.KeyUserHost = "host"
		require.True(t, key1.Equals(key2))
	})

	t.Run("Key materual, user and host set equal", func(t *testing.T) {
		key1 := sshutils.NewSSHPublicKey()
		key1.KeyMaterial = "AAAabc"
		key1.KeyUserHost = "user"
		key1.KeyUserHost = "host"
		key2 := sshutils.NewSSHPublicKey()
		key2.KeyMaterial = "AAAabc"
		key2.KeyUserHost = "user"
		key2.KeyUserHost = "host"
		require.True(t, key1.Equals(key2))
	})
}

func Test_GetCurrentUsersSshDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				sshDir, err := sshutils.GetCurrentUsersSshDirectory()
				require.NoError(t, err)
				require.True(t, strings.HasSuffix(sshDir.MustGetLocalPath(), "/.ssh"))
				require.True(t, pathsutils.IsAbsolutePath(sshDir.MustGetLocalPath()))
			},
		)
	}
}
