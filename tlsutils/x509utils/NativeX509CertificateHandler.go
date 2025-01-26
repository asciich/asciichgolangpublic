package x509utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type NativeX509CertificateHandler struct {
}

func NewNativeX509CertificateHandler() (handler *NativeX509CertificateHandler) {
	return new(NativeX509CertificateHandler)
}

func GetNativeX509CertificateHandler() (Handler X509CertificateHandler) {
	return NewNativeX509CertificateHandler()
}

func (n *NativeX509CertificateHandler) CreateRootCaCertificate(options *X509CreateCertificateOptions) (cert *x509.Certificate, privateKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	countryName, err := options.GetCountryName()
	if err != nil {
		return nil, nil, err
	}

	locality, err := options.GetLocality()
	if err != nil {
		return nil, nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{countryName},
			Province:      []string{""},
			Locality:      []string{locality},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	generatedKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Unable to create private key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &generatedKey.PublicKey, generatedKey)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Unable to create private key: %w", err)
	}

	cert, err = LoadCertificateFromDerBytes(certBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, generatedKey, nil
}
