package nativegnupg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/gnupgoptions"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/nativegnupg"
)

// This example shows how to:
//  1. Generate a key pair
//  2. Sign a temporary file
//  3. Valiadate the signature
func Test_Example_ValidateSignature(t *testing.T) {
	// enable verbose output
	ctx := contextutils.ContextVerbose()

	// Generate a new key pair:
	privateKey, publicKey, err := nativegnupg.GenerateKeyPair(ctx, &gnupgoptions.GenerateKeyPairOptions{
		Name: "reto",
		Email: "reto@example.com",
		Comment: "Example key",
		RSABits: 1024,
	})
	require.NoError(t, err)

	// Generate a temporary file to sign:
	tempFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tempFilePath, &filesoptions.DeleteOptions{})

	// Sign the temporary file:
	signature, err := nativegnupg.DetachSignFile(ctx, tempFilePath, privateKey)
	require.NoError(t, err)

	// Validate the signature:
	err = nativegnupg.CheckFileSignatureValid(ctx, tempFilePath, signature, [][]byte{publicKey})
	require.NoError(t, err)
}
