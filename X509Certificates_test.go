package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX509CertificatesCreateRootCaIntoTemporaryDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				certificates := X509Certificates()
				tempDirectory := certificates.MustCreateRootCaIntoDirectory(&X509CreateCertificateOptions{
					UseTemporaryDirectory: true,
					Verbose:               verbose,
					CommonName:            "test-ca.asciich.ch",
					CountryName:           "CH",
					Locality:              "Zurich",
				})

				assert.True(tempDirectory.MustExists(verbose))
				crtFile := tempDirectory.MustGetFileInDirectory("rootCA.crt")
				crtCertFile := MustGetX509CertificateFileFromFile(crtFile)
				assert.True(crtFile.MustExists(verbose))
				assert.True(crtCertFile.MustIsX509Certificate(verbose))
				assert.True(crtCertFile.MustIsX509RootCertificate(verbose))

				keyFile := tempDirectory.MustGetFileInDirectory("rootCA.key")
				assert.True(keyFile.MustExists(verbose))
			},
		)
	}
}

func TestX509CertificatesCreateIntermediateCertificateIntoTemporaryDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				certificates := X509Certificates()
				tempDirectory := certificates.MustCreateIntermediateCertificateIntoDirectory(&X509CreateCertificateOptions{
					UseTemporaryDirectory: true,
					Verbose:               verbose,
					CommonName:            "test-intermediata.asciich.ch",
					CountryName:           "CH",
					Locality:              "Zurich",
				})

				assert.True(tempDirectory.MustExists(verbose))

				keyFile := tempDirectory.MustGetFileInDirectory("intermediateCertificate.key")
				assert.True(keyFile.MustExists(verbose))
			},
		)
	}
}

func TestX509CertificateCreateAndSignIntermediateCertificate(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Not implemented on Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

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

				assert.True(rootTempDirectory.MustExists(verbose))
				rootCrtFile := MustGetX509CertificateFileFromFile(
					rootTempDirectory.MustGetFileInDirectory("rootCA.crt"),
				)

				assert.True(rootCrtFile.MustExists(verbose))
				assert.True(rootCrtFile.MustIsX509Certificate(verbose))
				assert.True(rootCrtFile.MustIsX509RootCertificate(verbose))
				assert.False(rootCrtFile.MustIsX509IntermediateCertificate())
				assert.True(rootCrtFile.MustIsX509v3())

				rootKeyFile := rootTempDirectory.MustGetFileInDirectory("rootCA.key")
				assert.True(rootKeyFile.MustExists(verbose))

				// Create intermediate certificate
				intermediateTempDirectory := certificates.MustCreateIntermediateCertificateIntoDirectory(&X509CreateCertificateOptions{
					UseTemporaryDirectory: true,
					Verbose:               verbose,
					CommonName:            "test-intermediate.asciich.ch",
					CountryName:           "CH",
					Locality:              "Zurich",
				})

				assert.True(intermediateTempDirectory.MustExists(verbose))

				intermediateKeyFile := intermediateTempDirectory.MustGetFileInDirectory("intermediateCertificate.key")
				assert.True(intermediateKeyFile.MustExists(verbose))

				intermediateCertFile := MustGetX509CertificateFileFromFile(
					intermediateTempDirectory.MustGetFileInDirectory("intermediateCertificate.crt"),
				)
				assert.False(intermediateCertFile.MustExists(verbose))

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

				assert.True(intermediateCertFile.MustExists(verbose))
				assert.False(intermediateCertFile.MustIsX509RootCertificate(verbose))
				assert.True(intermediateCertFile.MustIsX509IntermediateCertificate())
				assert.True(intermediateCertFile.MustIsX509v3())
				assert.True(intermediateCertFile.MustIsX509CertificateSignedByCertificateFile(rootCrtFile, verbose))
			},
		)
	}
}

// Ensure only expired certificates are included into the testdata directory for security reasons.
func TestX509Certificates_NoTestdataCertificateUnexpired(t *testing.T) {
	const verbose bool = true

	type TestCase struct {
		pathToCheck string
	}

	tests := []TestCase{}

	repoRoot := Git().MustGetRepositoryRootPathByPath(".", verbose)
	pathsToCheck := Bash().MustRunOneLinerAndGetStdoutAsLines(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				toCheck := MustGetX509CertificateFileFromPath(tt.pathToCheck)

				assert.True(toCheck.MustIsX509Certificate(verbose))
				assert.True(toCheck.MustIsExpired(verbose))
			},
		)
	}
}
