package x509utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestX509CertificatesCreateIntermediateCertificateIntoTemporaryDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const verbose bool = true

				certificates := X509Certificates()
				tempDirectory, err := certificates.CreateIntermediateCertificateIntoDirectory(
					ctx,
					&X509CreateCertificateOptions{
						UseTemporaryDirectory: true,
						CommonName:            "test-intermediata.asciich.ch",
						CountryName:           "CH",
						Locality:              "Zurich",
					},
				)
				require.NoError(t, err)

				exists, err := tempDirectory.Exists(verbose)
				require.NoError(t, err)
				require.True(t, exists)

				keyFile, err := tempDirectory.GetFileInDirectory("intermediateCertificate.key")
				require.NoError(t, err)
				require.True(t, keyFile.MustExists(verbose))
			},
		)
	}
}

/* TODO move to gopass
func TestX509CertificateCreateAndSignIntermediateCertificate(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				// Create root CA:
				certificates := X509Certificates()
				rootTempDirectory := certificates.MustCreateRootCaIntoDirectory(&X509CreateCertificateOptions{
					UseTemporaryDirectory: true,
					Verbose:               verbose,
					CommonName:            "test-ca.asciich.ch",
					CountryName:           "CH",
					Locality:              "Zurich",
				})

				require.True(rootTempDirectory.MustExists(verbose))
				rootCrtFile := MustGetX509CertificateFileFromFile(
					rootTempDirectory.MustGetFileInDirectory("rootCA.crt"),
				)

				require.True(rootCrtFile.MustExists(verbose))
				require.True(rootCrtFile.MustIsX509Certificate(verbose))
				require.True(rootCrtFile.MustIsX509RootCertificate(verbose))
				require.False(rootCrtFile.MustIsX509IntermediateCertificate())
				require.True(rootCrtFile.MustIsX509v3())

				rootKeyFile := rootTempDirectory.MustGetFileInDirectory("rootCA.key")
				require.True(rootKeyFile.MustExists(verbose))

				// Create intermediate certificate
				intermediateTempDirectory := certificates.MustCreateIntermediateCertificateIntoDirectory(&X509CreateCertificateOptions{
					UseTemporaryDirectory: true,
					Verbose:               verbose,
					CommonName:            "test-intermediate.asciich.ch",
					CountryName:           "CH",
					Locality:              "Zurich",
				})

				require.True(intermediateTempDirectory.MustExists(verbose))

				intermediateKeyFile := intermediateTempDirectory.MustGetFileInDirectory("intermediateCertificate.key")
				require.True(intermediateKeyFile.MustExists(verbose))

				intermediateCertFile := MustGetX509CertificateFileFromFile(
					intermediateTempDirectory.MustGetFileInDirectory("intermediateCertificate.crt"),
				)
				require.False(intermediateCertFile.MustExists(verbose))

				// Sign intermediate certificate
				certificates.MustSignIntermediateCertificate(&X509SignCertificateOptions{
					KeyFileUsedForSigning:  rootKeyFile,
					CertFileUsedForSigning: rootCrtFile,
					KeyFileToSign:          intermediateKeyFile,
					OutputCertificateFile:  intermediateCertFile,
					CommonName:             "test-intermediate.asciich.ch",
					CountryName:            "CH",
					Locality:               "Zurich",
					Verbose:                verbose,
				})

				require.True(intermediateCertFile.MustExists(verbose))
				require.False(intermediateCertFile.MustIsX509RootCertificate(verbose))
				require.True(intermediateCertFile.MustIsX509IntermediateCertificate())
				require.True(intermediateCertFile.MustIsX509v3())
				require.True(intermediateCertFile.MustIsX509CertificateSignedByCertificateFile(rootCrtFile, verbose))
			},
		)
	}
}
*/

// Ensure only expired certificates are included into the testdata directory for security reasons.
/* TODO enable again
func TestX509Certificates_NoTestdataCertificateUnexpired(t *testing.T) {
	const verbose bool = true

	type TestCase struct {
		pathToCheck string
	}

	tests := []TestCase{}

	repoRoot := mustRepoRoot()
	pathsToCheck := commandexecutorbashoo.Bash().MustRunOneLinerAndGetStdoutAsLines(
		fmt.Sprintf(
			"grep -l -r 'CERTIFICATE' '%s/testdata'",
			repoRoot,
		),
		verbose,
	)

	for _, pathToCheck := range pathsToCheck {
		tests = append(tests, TestCase{pathToCheck: pathToCheck})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				toCheck := MustGetX509CertificateFileFromPath(tt.pathToCheck)

				require.True(toCheck.MustIsX509Certificate(verbose))
				require.True(toCheck.MustIsExpired(verbose))
			},
		)
	}
}
*/
