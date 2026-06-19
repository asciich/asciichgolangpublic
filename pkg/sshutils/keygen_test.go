package sshutils_test

import (
	"context"
	"crypto/ed25519"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
	"golang.org/x/crypto/ssh"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// tempPaths returns two temporary file paths (private, public) inside t.TempDir().
// The files are NOT created — only the paths are returned.
func tempPaths(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "id_ed25519"), filepath.Join(dir, "id_ed25519.pub")
}

// TestGenerateSshKeyPair_HappyPath verifies that both files are created and
// contain a valid, matching Ed25519 key pair.
func TestGenerateSshKeyPair_HappyPath(t *testing.T) {
	ctx := getCtx()

	privPath, pubPath := tempPaths(t)

	err := sshutils.GenerateSshKeyPair(ctx, privPath, pubPath)
	require.NoError(t, err)

	// --- private key file ---
	privData, err := os.ReadFile(privPath)
	require.NoError(t, err, "cannot read private key file")

	signer, err := ssh.ParsePrivateKey(privData)
	require.NoError(t, err, "failed to parse private key")
	require.Equal(t, ssh.KeyAlgoED25519, signer.PublicKey().Type(), "unexpected key algorithm")

	// --- public key file ---
	pubData, err := os.ReadFile(pubPath)
	require.NoError(t, err, "cannot read public key file")

	parsedPub, _, _, _, err := ssh.ParseAuthorizedKey(pubData)
	require.NoError(t, err, "failed to parse public key")
	require.Equal(t, ssh.KeyAlgoED25519, parsedPub.Type(), "unexpected public key algorithm")

	// --- keys must form a matching pair ---
	// Extract the raw ed25519.PublicKey from both sides and compare.
	cryptoPub := signer.PublicKey().(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
	parsedCryptoPub := parsedPub.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
	require.True(t, cryptoPub.Equal(parsedCryptoPub), "public key from private key file does not match the public key file")
}

// TestGenerateSshKeyPair_FilePermissions checks that the private key is
// written with mode 0600 and the public key with mode 0644.
func TestGenerateSshKeyPair_FilePermissions(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	privPath, pubPath := tempPaths(t)

	err := sshutils.GenerateSshKeyPair(ctx, privPath, pubPath)
	require.NoError(t, err)

	for _, tc := range []struct {
		name     string
		path     string
		wantPerm os.FileMode
	}{
		{"private key", privPath, 0600},
		{"public key", pubPath, 0644},
	} {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.path)
			require.NoError(t, err, "stat failed")
			require.Equal(t, tc.wantPerm, info.Mode().Perm(), "unexpected file permissions")
		})
	}
}

// TestGenerateSshKeyPair_CreatesParentDirs verifies that missing parent
// directories are created automatically.
func TestGenerateSshKeyPair_CreatesParentDirs(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	base := t.TempDir()
	privPath := filepath.Join(base, "a", "b", "id_ed25519")
	pubPath := filepath.Join(base, "a", "b", "id_ed25519.pub")

	err := sshutils.GenerateSshKeyPair(ctx, privPath, pubPath)
	require.NoError(t, err)

	_, err = os.Stat(privPath)
	require.NoError(t, err, "private key not found after dir creation")

	_, err = os.Stat(pubPath)
	require.NoError(t, err, "public key not found after dir creation")
}

// TestGenerateSshKeyPair_CancelledContext ensures the function respects an
// already-cancelled context and returns an error without writing any files.
func TestGenerateSshKeyPair_CancelledContext(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	privPath, pubPath := tempPaths(t)

	ctx, cancel := context.WithCancel(ctx)
	cancel() // cancel immediately

	err := sshutils.GenerateSshKeyPair(ctx, privPath, pubPath)
	require.Error(t, err, "expected an error for cancelled context")

	// Neither file should have been written.
	for _, p := range []string{privPath, pubPath} {
		_, statErr := os.Stat(p)
		require.True(t, os.IsNotExist(statErr), "file %q should not exist after cancelled context", p)
	}
}

// TestGenerateSshKeyPair_EmptyPaths checks that empty path arguments are
// rejected with an error before any key material is generated.
func TestGenerateSshKeyPair_EmptyPaths(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, tc := range []struct {
		name    string
		privKey string
		pubKey  string
	}{
		{"empty private path", "", "/tmp/id.pub"},
		{"empty public path", "/tmp/id_ed25519", ""},
		{"both empty", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := sshutils.GenerateSshKeyPair(ctx, tc.privKey, tc.pubKey)
			require.Error(t, err, "expected error for empty path")
		})
	}
}

// TestGenerateSshKeyPair_KeysAreUnique checks that two successive calls
// produce different key material (i.e., the RNG is not reused/seeded).
func TestGenerateSshKeyPair_KeysAreUnique(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	privPath1, pubPath1 := tempPaths(t)
	privPath2, pubPath2 := tempPaths(t)

	err := sshutils.GenerateSshKeyPair(ctx, privPath1, pubPath1)
	require.NoError(t, err, "first call failed")

	err = sshutils.GenerateSshKeyPair(ctx, privPath2, pubPath2)
	require.NoError(t, err, "second call failed")

	pub1, err := os.ReadFile(pubPath1)
	require.NoError(t, err, "cannot read first public key")

	pub2, err := os.ReadFile(pubPath2)
	require.NoError(t, err, "cannot read second public key")

	require.NotEqual(t, string(pub1), string(pub2), "two successive calls produced identical public keys — RNG may be broken")
}

// TestGenerateSshKeyPair_SshKeygenValidation uses the ssh-keygen binary to
// independently validate that the generated private key is well-formed and
// that the public key file matches the private key.
// The test is skipped automatically if ssh-keygen is not available on the host.
func TestGenerateSshKeyPair_SshKeygenValidation(t *testing.T) {
	sshKeygen, err := exec.LookPath("ssh-keygen")
	if err != nil {
		t.Skip("ssh-keygen not available on this system")
	}

	privPath, pubPath := tempPaths(t)

	err = sshutils.GenerateSshKeyPair(context.Background(), privPath, pubPath)
	require.NoError(t, err)

	// `ssh-keygen -y -f <privkey>` derives the public key directly from the
	// private key and prints it in authorized_keys format — without touching
	// the .pub file at all. We can therefore compare it byte-for-byte against
	// the .pub file we wrote to confirm both sides are consistent.
	cmd := exec.Command(sshKeygen, "-y", "-f", privPath)
	derivedPub, err := cmd.CombinedOutput()
	require.NoError(t, err, "ssh-keygen failed to read private key: %s", string(derivedPub))
	require.NotEmpty(t, derivedPub, "ssh-keygen produced no output for private key")

	writtenPub, err := os.ReadFile(pubPath)
	require.NoError(t, err, "cannot read public key file")

	// Both must contain the same key material. Trim trailing whitespace/newlines
	// since ssh-keygen may append a trailing newline that differs from ours.
	require.Equal(t,
		strings.TrimSpace(string(writtenPub)),
		strings.TrimSpace(string(derivedPub)),
		"public key file does not match the public key derived from the private key by ssh-keygen",
	)
}
