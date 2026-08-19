package x509utils_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

// ExampleReadCertificateFromFile demonstrates how to load an X509 certificate from a file.
//
// The function reads a PEM-encoded certificate file from disk and returns it as
// an *x509.Certificate for further processing.
//
// This example shows:
//  1. Creating a test certificate
//  2. Writing it to a temporary file
//  3. Loading it back using ReadCertificateFromFile
//  4. Verifying the loaded certificate matches the original
func ExampleReadCertificateFromFile() {
	ctx := contextutils.ContextVerbose()

	// Create a test certificate
	rootCaCertAndKey, err := genericx509utils.CreateRootCaCertificate(
		ctx,
		&x509options.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "ExampleOrg",
			PrivateKeySize: 1024, // Small key size for example speed
		},
	)
	if err != nil {
		fmt.Printf("Failed to create certificate: %v\n", err)
		return
	}

	// Encode the certificate as PEM string
	pemEncoded, err := genericx509utils.WriteCertificateAsPEMString(rootCaCertAndKey.Cert)
	if err != nil {
		fmt.Printf("Failed to encode certificate: %v\n", err)
		return
	}

	// Write to a temporary file (in real usage, you would read an existing file)
	tmpDir, err := os.MkdirTemp("", "x509utils-example-*")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	tmpFilePath := filepath.Join(tmpDir, "certificate.pem")
	err = os.WriteFile(tmpFilePath, []byte(pemEncoded), 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		return
	}

	// Load the certificate from file
	loadedCert, err := x509utils.ReadCertificateFromFile(ctx, tmpFilePath)
	if err != nil {
		fmt.Printf("Failed to load certificate from file: %v\n", err)
		return
	}

	// Verify the loaded certificate matches the original
	if loadedCert.Equal(rootCaCertAndKey.Cert) {
		fmt.Println("Certificate loaded successfully and matches original")
	}

	// Output: Certificate loaded successfully and matches original
}
