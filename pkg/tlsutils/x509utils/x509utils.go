package x509utils

import (
	"context"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/nativex509utils"
)

func ReadCertificateFromFile(ctx context.Context, pathToRead string) (*x509.Certificate, error) {
	return nativex509utils.ReadCertificateFromFile(ctx, pathToRead)
}
