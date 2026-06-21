package sshutils

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/sshoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/crypto/ssh"
)

// GenerateSshKeyPair generates an Ed25519 SSH key pair and writes the private
// key (OpenSSH PEM format) to privateKeyPath and the public key (authorized_keys
// format) to publicKeyPath.
// The function respects context cancellation before any I/O is performed.
func GenerateSshKeyPair(ctx context.Context, options *sshoptions.GenerateKeyOptions) (*SSHKeyPair, error) {
	if err := ctx.Err(); err != nil {
		return nil, tracederrors.TracedErrorf("keygen: context already done: %w", err)
	}

	if options == nil {
		options = &sshoptions.GenerateKeyOptions{}
	}

	logging.LogInfoByCtxf(ctx, "Generate SSH key pair started.")

	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, tracederrors.TracedErrorf("keygen: failed to generate ed25519 key pair: %w", err)
	}

	privPEM, err := ssh.MarshalPrivateKey(privKey, "" /* no passphrase */)
	if err != nil {
		return nil, tracederrors.TracedErrorf("keygen: failed to marshal private key: %w", err)
	}

	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("keygen: failed to create ssh public key: %w", err)
	}
	authorizedKey := ssh.MarshalAuthorizedKey(sshPubKey)

	if err := ctx.Err(); err != nil {
		return nil, tracederrors.TracedErrorf("keygen: context cancelled before writing files: %w", err)
	}

	encodedPrivateKey := pem.EncodeToMemory(privPEM)

	privateKeyPath := options.PrivateKeyPath
	if len(privateKeyPath) > 0 {
		permPrivateKey := os.FileMode(0600)
		err = nativefiles.WriteBytes(ctx, privateKeyPath, encodedPrivateKey, &filesoptions.WriteOptions{
			Perm: &permPrivateKey,
		})
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Wrote generated SSH private key to '%s'.", privateKeyPath)
	}

	publicKeyPath := options.PublicKeyPath

	if len(publicKeyPath) > 0 {
		permPublicKey := os.FileMode(0644)
		err = nativefiles.WriteBytes(ctx, publicKeyPath, authorizedKey, &filesoptions.WriteOptions{
			Perm: &permPublicKey,
		})
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Wrote generated SSH public key to '%s'.", publicKeyPath)
	}

	logging.LogInfoByCtxf(ctx, "Generate SSH key pair finished. Private key path is '%s' and public key path is '%s'.", privateKeyPath, publicKeyPath)

	return &SSHKeyPair{
		PublicKey: &SSHPublicKey{
			KeyType:     SSH_KEY_TYPE_ED25519,
			KeyMaterial: string(authorizedKey),
		},
		PrivateKey: &SSHPrivateKey{
			KeyType:     SSH_KEY_TYPE_ED25519,
			KeyMaterial: string(encodedPrivateKey),
		},
	}, nil
}
