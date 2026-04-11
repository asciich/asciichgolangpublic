package nativegnupg

import (
	"bytes"
	"context"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func DetachSignFile(ctx context.Context, path string, privateKey []byte) (signature []byte, err error) {
	err = contextutils.CheckContextStillAlive(ctx)
	if err != nil {
		return nil, err
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	if privateKey == nil {
		return nil, tracederrors.TracedErrorNil("privateKey")
	}

	logging.LogInfoByCtxf(ctx, "Create detached signature for '%s' started.", path)

	entity, err := entityByPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Use private key '%s' to sign '%s'.", getFingerprintFromEntity(entity), path)

	fileToSign, err := os.Open(path)
	if err != nil {
		tracederrors.TracedErrorf("Failed to open file for signing: %w", err)
	}
	defer fileToSign.Close()

	var sigBuf bytes.Buffer
	sigWriter, err := armor.Encode(&sigBuf, "PGP SIGNATURE", nil)
	if err != nil {
		tracederrors.TracedErrorf("Failed to create armor encoder for signature: %w", err)
	}

	if err := openpgp.DetachSign(sigWriter, entity, fileToSign, &packet.Config{}); err != nil {
		tracederrors.TracedErrorf("Failed to sign file: %w", err)
	}
	sigWriter.Close()

	signature = sigBuf.Bytes()

	logging.LogInfoByCtxf(ctx, "Create detached signature for '%s' using private key '%s' finished.", path, getFingerprintFromEntity(entity))

	return signature, nil
}

func CheckFileSignatureValid(ctx context.Context, path string, signature []byte, trustedKeys [][]byte) error {
	if path == "" {
		return tracederrors.TracedErrorNil("path")
	}

	if signature == nil {
		return tracederrors.TracedErrorNil("signature")
	}

	if trustedKeys == nil {
		return tracederrors.TracedErrorNil("trustedKeys")
	}

	logging.LogInfoByCtxf(ctx, "Validate file '%s' using GnuPG signature started.", path)

	keyRing, err := loadTrustedKeys(ctx, trustedKeys)
	if err != nil {
		return err
	}

	fileToVerify, err := os.Open(path)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open signature file: %w", err)
	}
	defer fileToVerify.Close()

	sigBuffer := bytes.NewBuffer(signature)

	signer, err := openpgp.CheckArmoredDetachedSignature(
		keyRing, // ← all trusted keys checked at once
		fileToVerify,
		sigBuffer,
		nil,
	)
	if err != nil {
		return tracederrors.TracedErrorf("Vailed to verify GnuPG signature for '%s': %w", path, err)
	}

	signerFingerprint := getFormatedFingerprint(signer.PrimaryKey.Fingerprint)

	logging.LogInfoByCtxf(ctx, "Validate file '%s' using GnuPG signature finished. Valid signature by '%s'", path, signerFingerprint)

	return nil
}

func loadTrustedKeys(ctx context.Context, trustedKeys [][]byte) (openpgp.EntityList, error) {
	var keyRing openpgp.EntityList

	for i, keyBytes := range trustedKeys {
		entities, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(keyBytes))
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to read trusted key at index %d: %w", i, err)
		}
		keyRing = append(keyRing, entities...)
	}

	if len(keyRing) == 0 {
		return nil, tracederrors.TracedErrorf("no valid trusted keys found")
	}

	logging.LogInfoByCtxf(ctx, "Loaded %d trusted GnuPG keys.", len(keyRing))

	return keyRing, nil
}
