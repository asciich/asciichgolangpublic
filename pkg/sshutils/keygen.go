package sshutils

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/crypto/ssh"
)

// GenerateSshKeyPair generates an Ed25519 SSH key pair and writes the private
// key (OpenSSH PEM format) to privateKeyPath and the public key (authorized_keys
// format) to publicKeyPath.
// The function respects context cancellation before any I/O is performed.
func GenerateSshKeyPair(ctx context.Context, privateKeyPath string, publicKeyPath string) error {
	if err := ctx.Err(); err != nil {
		return tracederrors.TracedErrorf("keygen: context already done: %w", err)
	}

	if privateKeyPath == "" {
		return tracederrors.TracedErrorEmptyString("privateKeyPath")
	}

	if publicKeyPath == "" {
		return tracederrors.TracedErrorEmptyString("publicKeyPath")
	}

	logging.LogInfoByCtxf(ctx, "Generate SSH key pair started. Private key path is '%s' and public key path is '%s'.", privateKeyPath, publicKeyPath)

	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return tracederrors.TracedErrorf("keygen: failed to generate ed25519 key pair: %w", err)
	}

	privPEM, err := ssh.MarshalPrivateKey(privKey, "" /* no passphrase */)
	if err != nil {
		return tracederrors.TracedErrorf("keygen: failed to marshal private key: %w", err)
	}

	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return tracederrors.TracedErrorf("keygen: failed to create ssh public key: %w", err)
	}
	authorizedKey := ssh.MarshalAuthorizedKey(sshPubKey)

	if err := ctx.Err(); err != nil {
		return tracederrors.TracedErrorf("keygen: context cancelled before writing files: %w", err)
	}

	permPrivateKey := os.FileMode(0600)
	err = nativefiles.WriteBytes(ctx, privateKeyPath, pem.EncodeToMemory(privPEM), &filesoptions.WriteOptions{
		Perm: &permPrivateKey,
	})
	if err != nil {
		return err
	}

	permPublicKey := os.FileMode(0644)
	err = nativefiles.WriteBytes(ctx, publicKeyPath, authorizedKey, &filesoptions.WriteOptions{
		Perm: &permPublicKey,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Generate SSH key pair finished. Private key path is '%s' and public key path is '%s'.", privateKeyPath, publicKeyPath)

	return nil
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
