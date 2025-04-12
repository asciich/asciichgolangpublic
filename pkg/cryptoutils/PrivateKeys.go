package cryptoutils

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func LoadPrivateKeyFromPEMString(pemEncoded string) (privateKey crypto.PrivateKey, err error) {
	if pemEncoded == "" {
		return nil, tracederrors.TracedErrorEmptyString("pemEncoded")
	}

	block, _ := pem.Decode([]byte(pemEncoded))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, tracederrors.TracedError("invalid private key.")
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		privateKey = key
	} else if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		privateKey = key
	} else {
		return nil, tracederrors.TracedErrorf("Unable to parse as key: %w", err)
	}

	return privateKey, nil
}

func EncodePrivateKeyAsPEMString(privateKey crypto.PrivateKey) (pemEncoded string, err error) {
	if privateKey == nil {
		return "", tracederrors.TracedErrorNil("privateKey")
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to marshal private key: '%w'", err)
	}

	var buf bytes.Buffer
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(io.Writer(&buf), pemBlock)
	if err != nil {
		return "", err
	}

	pemEncoded = buf.String()
	pemEncoded = stringsutils.EnsureEndsWithExactlyOneLineBreak(pemEncoded)

	const minLen = 50
	if len(pemEncoded) < minLen {
		return "", tracederrors.TracedErrorf("pemBytes has less than '%v' bytes which is not enough for a pem certificate", minLen)
	}

	return pemEncoded, nil
}

func GetPublicKeyFromPrivateKey(privateKey crypto.PrivateKey) (publicKey crypto.PublicKey, err error) {
	if privateKey == "" {
		return nil, tracederrors.TracedErrorNil("privateKey")
	}

	withPublic, ok := privateKey.(interface{ Public() crypto.PublicKey })
	if !ok {
		return nil, tracederrors.TracedErrorf("Unable to get publicKey out of privateKey. privateKey does not implement Public as expected (and done by all stdlib private keys).")
	}

	return withPublic.Public(), nil
}