package genericx509utils_test

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func generateCertWithSubject(t *testing.T, subject string, days int) string {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.pem")
	certPath := filepath.Join(tmpDir, "cert.pem")

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", fmt.Sprintf("%d", days),
		"-nodes",
		"-subj", subject,
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl command failed: %s", string(output))

	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)

	return string(certPEM)
}

func generateCertWithSANs(t *testing.T, subject string, sans []string) string {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.pem")
	certPath := filepath.Join(tmpDir, "cert.pem")
	extPath := filepath.Join(tmpDir, "ext.cnf")

	sanEntries := ""
	for i, san := range sans {
		if i > 0 {
			sanEntries += ","
		}
		sanEntries += "DNS:" + san
	}

	extContent := fmt.Sprintf("[v3_req]\nsubjectAltName=%s\n", sanEntries)
	err := os.WriteFile(extPath, []byte(extContent), 0644)
	require.NoError(t, err)

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "1",
		"-nodes",
		"-subj", subject,
		"-extensions", "v3_req",
		"-config", extPath,
	)

	// We need a full config for -config, build one
	fullConfig := fmt.Sprintf(
		"[req]\ndistinguished_name=dn\nx509_extensions=v3_req\n[dn]\n[v3_req]\nsubjectAltName=%s\n",
		sanEntries,
	)
	err = os.WriteFile(extPath, []byte(fullConfig), 0644)
	require.NoError(t, err)

	cmd = exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "1",
		"-nodes",
		"-subj", subject,
		"-config", extPath,
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl command failed: %s", string(output))

	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)

	return string(certPEM)
}

func Test_GetCommonName(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		cn, err := genericx509utils.GetCommonName(nil)
		require.Error(t, err)
		require.Empty(t, cn)
	})

	t.Run("returns common name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=my-server.example.com", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		cn, err := genericx509utils.GetCommonName(cert)
		require.NoError(t, err)
		require.Equal(t, "my-server.example.com", cn)
	})

	t.Run("empty common name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		cn, err := genericx509utils.GetCommonName(cert)
		require.NoError(t, err)
		require.Empty(t, cn)
	})
}

func Test_GetSerialNumberAsString(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		serial, err := genericx509utils.GetSerialNumberAsString(nil)
		require.Error(t, err)
		require.Empty(t, serial)
	})

	t.Run("returns non empty serial number string", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		serial, err := genericx509utils.GetSerialNumberAsString(cert)
		require.NoError(t, err)
		require.NotEmpty(t, serial)
	})

	t.Run("serial number is a valid decimal number", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		serial, err := genericx509utils.GetSerialNumberAsString(cert)
		require.NoError(t, err)

		parsed, ok := new(big.Int).SetString(serial, 10)
		require.True(t, ok, "serial number '%s' should be a valid decimal number", serial)
		require.True(t, parsed.Sign() > 0)
	})
}

func Test_GetSerialNumberAsBigInt(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		serial, err := genericx509utils.GetSerialNumberAsBigInt(nil)
		require.Error(t, err)
		require.Nil(t, serial)
	})

	t.Run("returns non nil big int", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		serial, err := genericx509utils.GetSerialNumberAsBigInt(cert)
		require.NoError(t, err)
		require.NotNil(t, serial)
		require.True(t, serial.Sign() > 0)
	})

	t.Run("consistent with string representation", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		serialBigInt, err := genericx509utils.GetSerialNumberAsBigInt(cert)
		require.NoError(t, err)

		serialString, err := genericx509utils.GetSerialNumberAsString(cert)
		require.NoError(t, err)

		require.Equal(t, serialBigInt.String(), serialString)
	})
}

func Test_GetSerialNumberAsHexColonSeparated(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		serial, err := genericx509utils.GetSerialNumberAsHexColonSeparated(nil)
		require.Error(t, err)
		require.Empty(t, serial)
	})

	t.Run("returns non empty hex string", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		serial, err := genericx509utils.GetSerialNumberAsHexColonSeparated(cert)
		require.NoError(t, err)
		require.NotEmpty(t, serial)
		require.Contains(t, serial, ":")
	})
}

func Test_GetCountryName(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		country, err := genericx509utils.GetCountryName(nil)
		require.Error(t, err)
		require.Empty(t, country)
	})

	t.Run("returns country name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		country, err := genericx509utils.GetCountryName(cert)
		require.NoError(t, err)
		require.Equal(t, "CH", country)
	})

	t.Run("empty country", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		country, err := genericx509utils.GetCountryName(cert)
		require.NoError(t, err)
		require.Empty(t, country)
	})
}

func Test_GetAdditionalSans(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		sans, err := genericx509utils.GetAdditionalSans(nil)
		require.Error(t, err)
		require.Nil(t, sans)
	})

	t.Run("cert without SANs returns error", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		sans, err := genericx509utils.GetAdditionalSans(cert)
		require.Error(t, err)
		require.Nil(t, sans)
	})

	t.Run("cert with SANs returns them", func(t *testing.T) {
		expectedSans := []string{"www.example.com", "api.example.com"}
		certPEM := generateCertWithSANs(t, "/C=CH/O=TestOrg/CN=example.com", expectedSans)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		sans, err := genericx509utils.GetAdditionalSans(cert)
		require.NoError(t, err)
		require.ElementsMatch(t, expectedSans, sans)
	})
}

func Test_GetAdditionalSansOrEmptySliceIfUnset(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		sans, err := genericx509utils.GetAdditionalSansOrEmptySliceIfUnset(nil)
		require.Error(t, err)
		require.Nil(t, sans)
	})

	t.Run("cert without SANs returns empty slice", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		sans, err := genericx509utils.GetAdditionalSansOrEmptySliceIfUnset(cert)
		require.NoError(t, err)
		require.NotNil(t, sans)
		require.Empty(t, sans)
	})

	t.Run("cert with SANs returns them", func(t *testing.T) {
		expectedSans := []string{"mail.example.com"}
		certPEM := generateCertWithSANs(t, "/C=CH/O=TestOrg/CN=example.com", expectedSans)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		sans, err := genericx509utils.GetAdditionalSansOrEmptySliceIfUnset(cert)
		require.NoError(t, err)
		require.ElementsMatch(t, expectedSans, sans)
	})
}

func Test_GetLocality(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		locality, err := genericx509utils.GetLocality(nil)
		require.Error(t, err)
		require.Empty(t, locality)
	})

	t.Run("returns locality", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/L=Zurich/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		locality, err := genericx509utils.GetLocality(cert)
		require.NoError(t, err)
		require.Equal(t, "Zurich", locality)
	})

	t.Run("empty locality", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		locality, err := genericx509utils.GetLocality(cert)
		require.NoError(t, err)
		require.Empty(t, locality)
	})
}

func Test_GetOrganization(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		org, err := genericx509utils.GetOrganization(nil)
		require.Error(t, err)
		require.Empty(t, org)
	})

	t.Run("returns organization", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=MyCompany/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		org, err := genericx509utils.GetOrganization(cert)
		require.NoError(t, err)
		require.Equal(t, "MyCompany", org)
	})

	t.Run("empty organization", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		org, err := genericx509utils.GetOrganization(cert)
		require.NoError(t, err)
		require.Empty(t, org)
	})
}

func Test_GetValidityDuration(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		duration, err := genericx509utils.GetValidityDuration(nil)
		require.Error(t, err)
		require.Nil(t, duration)
	})

	t.Run("returns positive duration", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 30)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		duration, err := genericx509utils.GetValidityDuration(cert)
		require.NoError(t, err)
		require.NotNil(t, duration)
		require.True(t, *duration > 0)
	})

	t.Run("duration matches expected days approximately", func(t *testing.T) {
		days := 365
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", days)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		duration, err := genericx509utils.GetValidityDuration(cert)
		require.NoError(t, err)
		require.NotNil(t, duration)

		expectedDuration := time.Duration(days) * 24 * time.Hour
		tolerance := 2 * time.Hour
		require.InDelta(t, float64(expectedDuration), float64(*duration), float64(tolerance))
	})
}

func Test_GetValidityDurationAsString(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		durationStr, err := genericx509utils.GetValidityDurationAsString(nil)
		require.Error(t, err)
		require.Empty(t, durationStr)
	})

	t.Run("returns non empty string", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		durationStr, err := genericx509utils.GetValidityDurationAsString(cert)
		require.NoError(t, err)
		require.NotEmpty(t, durationStr)
		require.Contains(t, durationStr, "h")
	})
}

func Test_GetSubjectAndSerialString(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		result, err := genericx509utils.GetSubjectAndSerialString(nil)
		require.Error(t, err)
		require.Empty(t, result)
	})

	t.Run("returns string with subject and serial", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		result, err := genericx509utils.GetSubjectAndSerialString(cert)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Contains(t, result, "serial:")
		require.Contains(t, result, ":")
	})
}

func Test_GetSubjectAsPkixName(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		_, err := genericx509utils.GetSubjectAsPkixName(nil)
		require.Error(t, err)
	})

	t.Run("returns correct pkix name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/L=Bern/O=PkixOrg/CN=PkixCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		subject, err := genericx509utils.GetSubjectAsPkixName(cert)
		require.NoError(t, err)
		require.Equal(t, "PkixCert", subject.CommonName)
		require.Equal(t, []string{"CH"}, subject.Country)
		require.Equal(t, []string{"Bern"}, subject.Locality)
		require.Equal(t, []string{"PkixOrg"}, subject.Organization)
	})
}

func Test_GetSubjectCountryName(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		country, err := genericx509utils.GetSubjectCountryName(nil)
		require.Error(t, err)
		require.Empty(t, country)
	})

	t.Run("returns country name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=DE/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		country, err := genericx509utils.GetSubjectCountryName(cert)
		require.NoError(t, err)
		require.Equal(t, "DE", country)
	})

	t.Run("consistent with GetCountryName", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=AT/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		countryFromSubject, err := genericx509utils.GetSubjectCountryName(cert)
		require.NoError(t, err)

		countryFromGet, err := genericx509utils.GetCountryName(cert)
		require.NoError(t, err)

		require.Equal(t, countryFromGet, countryFromSubject)
	})
}

func Test_GetSubjectLocalityName(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		locality, err := genericx509utils.GetSubjectLocalityName(nil)
		require.Error(t, err)
		require.Empty(t, locality)
	})

	t.Run("returns locality name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/L=Geneva/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		locality, err := genericx509utils.GetSubjectLocalityName(cert)
		require.NoError(t, err)
		require.Equal(t, "Geneva", locality)
	})

	t.Run("consistent with GetLocality", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/L=Basel/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		localityFromSubject, err := genericx509utils.GetSubjectLocalityName(cert)
		require.NoError(t, err)

		localityFromGet, err := genericx509utils.GetLocality(cert)
		require.NoError(t, err)

		require.Equal(t, localityFromGet, localityFromSubject)
	})
}

func Test_GetSubjectOrganizationName(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		org, err := genericx509utils.GetSubjectOrganizationName(nil)
		require.Error(t, err)
		require.Empty(t, org)
	})

	t.Run("returns organization name", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=SwissCorp/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		org, err := genericx509utils.GetSubjectOrganizationName(cert)
		require.NoError(t, err)
		require.Equal(t, "SwissCorp", org)
	})

	t.Run("consistent with GetOrganization", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/O=ConsOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		orgFromSubject, err := genericx509utils.GetSubjectOrganizationName(cert)
		require.NoError(t, err)

		orgFromGet, err := genericx509utils.GetOrganization(cert)
		require.NoError(t, err)

		require.Equal(t, orgFromGet, orgFromSubject)
	})
}

func Test_GetSubjectStringForOpenssl(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		result, err := genericx509utils.GetSubjectStringForOpenssl(nil)
		require.Error(t, err)
		require.Empty(t, result)
	})

	t.Run("full subject", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=CH/L=Zurich/O=TestOrg/CN=TestCert", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		result, err := genericx509utils.GetSubjectStringForOpenssl(cert)
		require.NoError(t, err)
		require.Contains(t, result, "/C=CH")
		require.Contains(t, result, "/L=Zurich")
		require.Contains(t, result, "/O=TestOrg")
		require.Contains(t, result, "/CN=TestCert")
	})

	t.Run("partial subject only CN", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/CN=OnlyCN", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		result, err := genericx509utils.GetSubjectStringForOpenssl(cert)
		require.NoError(t, err)
		require.Equal(t, "/CN=OnlyCN", result)
	})

	t.Run("partial subject country and org", func(t *testing.T) {
		certPEM := generateCertWithSubject(t, "/C=FR/O=FrenchOrg", 1)
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		result, err := genericx509utils.GetSubjectStringForOpenssl(cert)
		require.NoError(t, err)
		require.Contains(t, result, "/C=FR")
		require.Contains(t, result, "/O=FrenchOrg")
		require.NotContains(t, result, "/L=")
		require.NotContains(t, result, "/CN=")
	})
}
