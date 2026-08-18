package genericx509utils

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const defaultPrivateKeySize = 4096
const defaultValidityDays = 365

func CreateRootCa(ctx context.Context, options *x509options.X509CreateCertificateOptions) (caCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	logging.LogInfoByCtx(ctx, "Create root CA certificate started.")

	privateKeySize := getPrivateKeySize(options)

	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate RSA private key: %w", err)
	}

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	subject := buildPkixName(options)

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		Issuer:                subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(defaultValidityDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            -1,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create root CA certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse created root CA certificate: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Create root CA certificate finished. Subject: '%s'.", cert.Subject.String())

	return &X509CertKeyPair{
		Cert: cert,
		Key:  privateKey,
	}, nil
}

func CreateIntermediateCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions) (intermediateCert *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	logging.LogInfoByCtx(ctx, "Create self-signed intermediate certificate started.")

	privateKeySize := getPrivateKeySize(options)

	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate RSA private key: %w", err)
	}

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	subject := buildPkixName(options)

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		Issuer:                subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(defaultValidityDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create intermediate certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse created intermediate certificate: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Create self-signed intermediate certificate finished. Subject: '%s'.", cert.Subject.String())

	return &X509CertKeyPair{
		Cert: cert,
		Key:  privateKey,
	}, nil
}

func CreateSelfSignedCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions) (selfSignedCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	logging.LogInfoByCtx(ctx, "Create self-signed certificate started.")

	privateKeySize := getPrivateKeySize(options)

	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate RSA private key: %w", err)
	}

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	subject := buildPkixName(options)

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		Issuer:                subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(defaultValidityDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	if options.CommonName != "" {
		template.DNSNames = []string{options.CommonName}
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create self-signed certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse created self-signed certificate: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Create self-signed certificate finished. Subject: '%s'.", cert.Subject.String())

	return &X509CertKeyPair{
		Cert: cert,
		Key:  privateKey,
	}, nil
}

func CreateSignedIntermediateCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions, rootCaCertAndKey *X509CertKeyPair) (intermediateCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if rootCaCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("rootCaCertAndKey")
	}

	logging.LogInfoByCtx(ctx, "Create signed intermediate certificate started.")

	rootCert, err := rootCaCertAndKey.GetX509Certificate()
	if err != nil {
		return nil, err
	}

	rootKey, err := rootCaCertAndKey.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKeySize := getPrivateKeySize(options)

	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate RSA private key: %w", err)
	}

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	subject := buildPkixName(options)

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(defaultValidityDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, rootCert, &privateKey.PublicKey, rootKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create signed intermediate certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse created signed intermediate certificate: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Create signed intermediate certificate finished. Subject: '%s', Issuer: '%s'.", cert.Subject.String(), cert.Issuer.String())

	return &X509CertKeyPair{
		Cert: cert,
		Key:  privateKey,
	}, nil
}

func CreateSignedEndEntityCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions, caCertAndKey *X509CertKeyPair) (endEntityCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if caCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("caCertAndKey")
	}

	logging.LogInfoByCtx(ctx, "Create signed end entity certificate started.")

	caCert, err := caCertAndKey.GetX509Certificate()
	if err != nil {
		return nil, err
	}

	caKey, err := caCertAndKey.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKeySize := getPrivateKeySize(options)

	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate RSA private key: %w", err)
	}

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return nil, err
	}

	subject := buildPkixName(options)

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(defaultValidityDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	if options.CommonName != "" {
		template.DNSNames = []string{options.CommonName}
	}

	if len(options.AdditionalSans) > 0 {
		template.DNSNames = append(template.DNSNames, options.AdditionalSans...)
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &privateKey.PublicKey, caKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create signed end entity certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse created signed end entity certificate: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Create signed end entity certificate finished. Subject: '%s', Issuer: '%s'.", cert.Subject.String(), cert.Issuer.String())

	return &X509CertKeyPair{
		Cert: cert,
		Key:  privateKey,
	}, nil
}

// --- Helper functions ---

func getPrivateKeySize(options *x509options.X509CreateCertificateOptions) int {
	if options.PrivateKeySize > 0 {
		return options.PrivateKeySize
	}

	return defaultPrivateKeySize
}

func buildPkixName(options *x509options.X509CreateCertificateOptions) pkix.Name {
	name := pkix.Name{}

	if options.CountryName != "" {
		name.Country = []string{options.CountryName}
	}

	if options.Locality != "" {
		name.Locality = []string{options.Locality}
	}

	if options.Organization != "" {
		name.Organization = []string{options.Organization}
	}

	if options.CommonName != "" {
		name.CommonName = options.CommonName
	}

	return name
}

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate certificate serial number: %w", err)
	}

	return serialNumber, nil
}

func GeneratePrivateKey(ctx context.Context) (crypto.PrivateKey, error) {
	logging.LogInfoByCtx(ctx, "Generate private key started.")

	privateKey, err := rsa.GenerateKey(rand.Reader, defaultPrivateKeySize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to generate RSA private key: %w", err)
	}

	logging.LogInfoByCtx(ctx, "Generate private key finished.")

	return privateKey, nil
}
