package genericx509utils

import (
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"
	"time"

	"crypto/x509/pkix"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/bigintutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetCommonName(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	return cert.Subject.CommonName, nil
}

func GetSerialNumberAsString(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	serial := cert.SerialNumber
	if serial == nil {
		return "", tracederrors.TracedError("unable to get serial number from x509 certificate. SerialNumber is nil")
	}

	serialNumber := serial.String()
	if serialNumber == "" {
		return "", tracederrors.TracedError("serial number is empty string after evaluation")
	}

	return serialNumber, nil
}

func GetSerialNumberAsBigInt(cert *x509.Certificate) (*big.Int, error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	serial := cert.SerialNumber
	if serial == nil {
		return nil, tracederrors.TracedError("unable to get serial number from x509 certificate. SerialNumber is nil")
	}

	return serial, nil
}

func GetSerialNumberAsHexColonSeparated(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	serial := cert.SerialNumber
	if serial == nil {
		return "", tracederrors.TracedError("unable to get serial number from x509 certificate. SerialNumber is nil")
	}

	return bigintutils.ToHexStringColonSeparated(serial)
}

func GetCountryName(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	country := cert.Subject.Country

	nCountries := len(country)
	if nCountries == 0 {
		return "", nil
	}

	if nCountries == 1 {
		return country[0], nil
	}

	return "", tracederrors.TracedErrorf(
		"Not implemented for nCountries != 1. Got '%d' countries: '%v'",
		nCountries,
		country,
	)
}

func GetAdditionalSans(cert *x509.Certificate) ([]string, error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	sans := cert.DNSNames
	if len(sans) == 0 {
		return nil, tracederrors.TracedError("no additional SANs found in certificate")
	}

	return sans, nil
}

func GetAdditionalSansOrEmptySliceIfUnset(cert *x509.Certificate) ([]string, error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	sans := cert.DNSNames
	if sans == nil {
		return []string{}, nil
	}

	return sans, nil
}

func GetLocality(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	locality := cert.Subject.Locality

	nLocalities := len(locality)
	if nLocalities == 0 {
		return "", nil
	}

	if nLocalities == 1 {
		return locality[0], nil
	}

	return "", tracederrors.TracedErrorf(
		"Not implemented for nLocalities != 1. Got '%d' localities: '%v'",
		nLocalities,
		locality,
	)
}

func GetOrganization(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	organization := cert.Subject.Organization

	nOrganizations := len(organization)
	if nOrganizations == 0 {
		return "", nil
	}

	if nOrganizations == 1 {
		return organization[0], nil
	}

	return "", tracederrors.TracedErrorf(
		"Not implemented for nOrganizations != 1. Got '%d' organizations: '%v'",
		nOrganizations,
		organization,
	)
}

func GetValidityDuration(cert *x509.Certificate) (*time.Duration, error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	diff := cert.NotAfter.Sub(cert.NotBefore)

	return &diff, nil
}

func GetValidityDurationAsString(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	duration, err := GetValidityDuration(cert)
	if err != nil {
		return "", err
	}

	return duration.String(), nil
}

func GetSubjectAndSerialString(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	serial, err := GetSerialNumberAsHexColonSeparated(cert)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s serial: %s",
		cert.Subject,
		serial,
	), nil
}

func GetSubjectAsPkixName(cert *x509.Certificate) (pkix.Name, error) {
	if cert == nil {
		return pkix.Name{}, tracederrors.TracedErrorNil("cert")
	}

	return cert.Subject, nil
}

func GetSubjectCountryName(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	return GetCountryName(cert)
}

func GetSubjectLocalityName(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	return GetLocality(cert)
}

func GetSubjectOrganizationName(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	return GetOrganization(cert)
}

func GetSubjectStringForOpenssl(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	var parts []string

	country := cert.Subject.Country
	if len(country) > 0 && country[0] != "" {
		parts = append(parts, fmt.Sprintf("/C=%s", country[0]))
	}

	locality := cert.Subject.Locality
	if len(locality) > 0 && locality[0] != "" {
		parts = append(parts, fmt.Sprintf("/L=%s", locality[0]))
	}

	organization := cert.Subject.Organization
	if len(organization) > 0 && organization[0] != "" {
		parts = append(parts, fmt.Sprintf("/O=%s", organization[0]))
	}

	commonName := cert.Subject.CommonName
	if commonName != "" {
		parts = append(parts, fmt.Sprintf("/CN=%s", commonName))
	}

	if len(parts) == 0 {
		return "", tracederrors.TracedError("certificate subject has no fields to format for openssl")
	}

	return strings.Join(parts, ""), nil
}

func GetInfoString(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	commonName, err := GetCommonName(cert)
	if err != nil {
		return "", err
	}

	serial, err := GetSerialNumberAsHexColonSeparated(cert)
	if err != nil {
		return "", err
	}

	notBefore := cert.NotBefore.Format(time.RFC1123)
	notAfter := cert.NotAfter.Format(time.RFC1123)

	infoString := fmt.Sprintf(
		"CN: %s, Serial Number: %s, Not Before: %s, Not After: %s",
		commonName,
		serial,
		notBefore,
		notAfter,
	)

	return infoString, nil
}
