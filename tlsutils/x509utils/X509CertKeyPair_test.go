package x509utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func Test_CertKeyPairMatches(t *testing.T) {
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
				handler := getX509CertificateHandlerToTest(tt.implementationName)

				certKeyPair, err := handler.CreateRootCaCertificate(
					getCtx(),
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
						SerialNumber: "12345",
					},
				)
				require.NoError(t, err)

				matching, err := certKeyPair.IsKeyMatchingCert()
				require.NoError(t, err)
				require.True(t, matching)
			},
		)
	}
}
