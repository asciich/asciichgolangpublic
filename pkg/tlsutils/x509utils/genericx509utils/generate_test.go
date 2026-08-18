package genericx509utils_test

import (
	"crypto"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

func getDefaultRootCaOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "CH",
		Locality:       "Zurich",
		Organization:   "TestRootOrg",
		CommonName:     "TestRootCA",
		PrivateKeySize: 2048,
	}
}

func getDefaultIntermediateOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "CH",
		Locality:       "Bern",
		Organization:   "TestIntOrg",
		CommonName:     "TestIntermediateCA",
		PrivateKeySize: 2048,
	}
}

func getDefaultEndEntityOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "CH",
		Locality:       "Basel",
		Organization:   "TestEEOrg",
		CommonName:     "server.example.com",
		PrivateKeySize: 2048,
	}
}

func getDefaultSelfSignedOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "DE",
		Locality:       "Berlin",
		Organization:   "SelfSignedOrg",
		CommonName:     "selfsigned.example.com",
		PrivateKeySize: 2048,
	}
}

func Test_CreateRootCa(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil options returns error", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, nil)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("returns valid root CA cert", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)
		require.NotNil(t, pair)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert)
	})

	t.Run("root CA has correct subject fields", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Zurich", cert.Subject.Locality[0])
		require.Equal(t, "TestRootOrg", cert.Subject.Organization[0])
		require.Equal(t, "TestRootCA", cert.Subject.CommonName)
	})

	t.Run("root CA is a CA certificate", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.True(t, cert.IsCA)
	})

	t.Run("root CA is self signed", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, cert.Subject.String(), cert.Issuer.String())
	})

	t.Run("root CA is detected as root CA", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isRootCa)
	})

	t.Run("root CA key matches cert", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		err = pair.CheckKeyMatchingCert()
		require.NoError(t, err)
	})

	t.Run("root CA can verify its own signature", func(t *testing.T) {
		pair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, cert, cert)
		require.NoError(t, err)
		require.True(t, isSigned)
	})

	t.Run("two root CAs have different serial numbers", func(t *testing.T) {
		pair1, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		pair2, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		cert1, err := pair1.GetX509Certificate()
		require.NoError(t, err)
		cert2, err := pair2.GetX509Certificate()
		require.NoError(t, err)

		require.NotEqual(t, 0, cert1.SerialNumber.Cmp(cert2.SerialNumber))
	})
}

func Test_CreateIntermediateCertificate(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil options returns error", func(t *testing.T) {
		pair, err := genericx509utils.CreateIntermediateCertificate(ctx, nil)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("returns valid intermediate cert", func(t *testing.T) {
		pair, err := genericx509utils.CreateIntermediateCertificate(ctx, getDefaultIntermediateOptions())
		require.NoError(t, err)
		require.NotNil(t, pair)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert)
	})

	t.Run("intermediate has correct subject fields", func(t *testing.T) {
		pair, err := genericx509utils.CreateIntermediateCertificate(ctx, getDefaultIntermediateOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Bern", cert.Subject.Locality[0])
		require.Equal(t, "TestIntOrg", cert.Subject.Organization[0])
		require.Equal(t, "TestIntermediateCA", cert.Subject.CommonName)
	})

	t.Run("intermediate is a CA certificate", func(t *testing.T) {
		pair, err := genericx509utils.CreateIntermediateCertificate(ctx, getDefaultIntermediateOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.True(t, cert.IsCA)
	})

	t.Run("intermediate key matches cert", func(t *testing.T) {
		pair, err := genericx509utils.CreateIntermediateCertificate(ctx, getDefaultIntermediateOptions())
		require.NoError(t, err)

		err = pair.CheckKeyMatchingCert()
		require.NoError(t, err)
	})
}

func Test_CreateSelfSignedCertificate(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil options returns error", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, nil)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("returns valid self signed cert", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, getDefaultSelfSignedOptions())
		require.NoError(t, err)
		require.NotNil(t, pair)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert)
	})

	t.Run("self signed has correct subject fields", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, getDefaultSelfSignedOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, "DE", cert.Subject.Country[0])
		require.Equal(t, "Berlin", cert.Subject.Locality[0])
		require.Equal(t, "SelfSignedOrg", cert.Subject.Organization[0])
		require.Equal(t, "selfsigned.example.com", cert.Subject.CommonName)
	})

	t.Run("self signed is self signed", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, getDefaultSelfSignedOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, cert.Subject.String(), cert.Issuer.String())
	})

	t.Run("self signed is not a CA", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, getDefaultSelfSignedOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.False(t, cert.IsCA)
	})

	t.Run("self signed key matches cert", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, getDefaultSelfSignedOptions())
		require.NoError(t, err)

		err = pair.CheckKeyMatchingCert()
		require.NoError(t, err)
	})

	t.Run("self signed has common name as SAN", func(t *testing.T) {
		pair, err := genericx509utils.CreateSelfSignedCertificate(ctx, getDefaultSelfSignedOptions())
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		require.Contains(t, cert.DNSNames, "selfsigned.example.com")
	})
}

func Test_CreateSignedIntermediateCertificate(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil options returns error", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		pair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, nil, rootPair)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("nil root CA returns error", func(t *testing.T) {
		pair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), nil)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("returns valid signed intermediate cert", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)
		require.NotNil(t, intPair)

		cert, err := intPair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert)
	})

	t.Run("signed intermediate has correct subject fields", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		cert, err := intPair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Bern", cert.Subject.Locality[0])
		require.Equal(t, "TestIntOrg", cert.Subject.Organization[0])
		require.Equal(t, "TestIntermediateCA", cert.Subject.CommonName)
	})

	t.Run("signed intermediate is a CA certificate", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		cert, err := intPair.GetX509Certificate()
		require.NoError(t, err)

		require.True(t, cert.IsCA)
	})

	t.Run("signed intermediate is not self signed", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		cert, err := intPair.GetX509Certificate()
		require.NoError(t, err)

		require.NotEqual(t, cert.Subject.String(), cert.Issuer.String())
	})

	t.Run("signed intermediate is detected as intermediate", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		cert, err := intPair.GetX509Certificate()
		require.NoError(t, err)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isIntermediate)
	})

	t.Run("signed intermediate is signed by root CA", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		rootCert, err := rootPair.GetX509Certificate()
		require.NoError(t, err)
		intCert, err := intPair.GetX509Certificate()
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, intCert, rootCert)
		require.NoError(t, err)
		require.True(t, isSigned)
	})

	t.Run("signed intermediate key matches cert", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		err = intPair.CheckKeyMatchingCert()
		require.NoError(t, err)
	})

	t.Run("signed intermediate issuer matches root CA subject", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		rootCert, err := rootPair.GetX509Certificate()
		require.NoError(t, err)
		intCert, err := intPair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, rootCert.Subject.String(), intCert.Issuer.String())
	})
}

func Test_CreateSignedEndEntityCertificate(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil options returns error", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		pair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, nil, rootPair)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("nil CA returns error", func(t *testing.T) {
		pair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), nil)
		require.Error(t, err)
		require.Nil(t, pair)
	})

	t.Run("returns valid signed end entity cert", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)
		require.NotNil(t, eePair)

		cert, err := eePair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert)
	})

	t.Run("end entity has correct subject fields", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		cert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Basel", cert.Subject.Locality[0])
		require.Equal(t, "TestEEOrg", cert.Subject.Organization[0])
		require.Equal(t, "server.example.com", cert.Subject.CommonName)
	})

	t.Run("end entity is not a CA certificate", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		cert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		require.False(t, cert.IsCA)
	})

	t.Run("end entity is detected as end entity", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		cert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isEndEntity)
	})

	t.Run("end entity is signed by CA", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		rootCert, err := rootPair.GetX509Certificate()
		require.NoError(t, err)
		eeCert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, eeCert, rootCert)
		require.NoError(t, err)
		require.True(t, isSigned)
	})

	t.Run("end entity key matches cert", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		err = eePair.CheckKeyMatchingCert()
		require.NoError(t, err)
	})

	t.Run("end entity issuer matches CA subject", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		rootCert, err := rootPair.GetX509Certificate()
		require.NoError(t, err)
		eeCert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		require.Equal(t, rootCert.Subject.String(), eeCert.Issuer.String())
	})

	t.Run("end entity has common name as SAN", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		cert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		require.Contains(t, cert.DNSNames, "server.example.com")
	})

	t.Run("end entity with additional SANs", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		options := &x509options.X509CreateCertificateOptions{
			CountryName:    "CH",
			Organization:   "SANOrg",
			CommonName:     "main.example.com",
			AdditionalSans: []string{"www.example.com", "api.example.com"},
			PrivateKeySize: 2048,
		}

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, options, rootPair)
		require.NoError(t, err)

		cert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		require.Contains(t, cert.DNSNames, "main.example.com")
		require.Contains(t, cert.DNSNames, "www.example.com")
		require.Contains(t, cert.DNSNames, "api.example.com")
	})

	t.Run("end entity is not signed by unrelated CA", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		unrelatedPair, err := genericx509utils.CreateRootCa(ctx, &x509options.X509CreateCertificateOptions{
			CountryName:    "FR",
			Organization:   "UnrelatedOrg",
			CommonName:     "UnrelatedRoot",
			PrivateKeySize: 2048,
		})
		require.NoError(t, err)

		unrelatedCert, err := unrelatedPair.GetX509Certificate()
		require.NoError(t, err)
		eeCert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, eeCert, unrelatedCert)
		require.NoError(t, err)
		require.False(t, isSigned)
	})

	t.Run("full chain root-intermediate-end entity is valid", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCa(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), intPair)
		require.NoError(t, err)

		rootCert, err := rootPair.GetX509Certificate()
		require.NoError(t, err)
		intCert, err := intPair.GetX509Certificate()
		require.NoError(t, err)
		eeCert, err := eePair.GetX509Certificate()
		require.NoError(t, err)

		isValidChain, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, intCert, rootCert})
		require.NoError(t, err)
		require.True(t, isValidChain)
	})
}

func Test_GeneratePrivateKey(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("returns non nil private key", func(t *testing.T) {
		key, err := genericx509utils.GeneratePrivateKey(ctx)
		require.NoError(t, err)
		require.NotNil(t, key)
	})

	t.Run("generated keys are unique", func(t *testing.T) {
		key1, err := genericx509utils.GeneratePrivateKey(ctx)
		require.NoError(t, err)

		key2, err := genericx509utils.GeneratePrivateKey(ctx)
		require.NoError(t, err)

		key1WithEqual, ok := key1.(interface {
			Equal(x crypto.PrivateKey) bool
		})
		require.True(t, ok)
		require.False(t, key1WithEqual.Equal(key2))
	})

	t.Run("generated key can extract public key", func(t *testing.T) {
		key, err := genericx509utils.GeneratePrivateKey(ctx)
		require.NoError(t, err)

		withPublic, ok := key.(interface{ Public() crypto.PublicKey })
		require.True(t, ok)

		pubKey := withPublic.Public()
		require.NotNil(t, pubKey)
	})
}
