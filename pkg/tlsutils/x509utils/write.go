package x509utils

import (
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func WriteCertificateAsPEMString(cert *x509.Certificate) (string, error) {
	return genericx509utils.WriteCertificateAsPEMString(cert)
}

func WriteCertificateAsPEMBytes(cert *x509.Certificate) ([]byte, error) {
	return genericx509utils.WriteCertificateAsPEMBytes(cert)
}

func WriteCertificatesAsPEMString(certs []*x509.Certificate) (string, error) {
	return genericx509utils.WriteCertificatesAsPEMString(certs)
}

func WriteCertificatesAsPEMBytes(certs []*x509.Certificate) ([]byte, error) {
	return genericx509utils.WriteCertificatesAsPEMBytes(certs)
}
