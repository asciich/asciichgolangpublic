package x509utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getX509CertificateHandlerToTest(implementationName string) (handler X509CertificateHandler) {
	if implementationName == "NativeX509CertificateHandler" {
		return GetNativeX509CertificateHandler()
	}

	logging.LogFatalWithTracef("Unknown implementationName '%s'", implementationName)

	return
}

func TestX509Handler_CreateRootCertificate(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				handler := getX509CertificateHandlerToTest(tt.implementationName)

				cert, _ := mustutils.Must2(handler.CreateRootCertificate(
					&X509CreateCertificateOptions{
						CountryName: "CH",
						Locality: "Zurich",
						
						Verbose:     true,
					},
				))

				assert.True(mustutils.Must(IsCertificateRootCa(cert)))

				assert.EqualValues([]string{"CH"}, cert.Issuer.Country)
				assert.EqualValues([]string{"Zurich"}, cert.Issuer.Locality)
			},
		)
	}
}
