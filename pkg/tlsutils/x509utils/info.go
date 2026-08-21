package x509utils

import (
	"crypto/x509"
	"math/big"
	"time"

	"crypto/x509/pkix"

	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func GetCommonName(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetCommonName(cert)
}

func GetSerialNumberAsString(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSerialNumberAsString(cert)
}

func GetSerialNumberAsBigInt(cert *x509.Certificate) (*big.Int, error) {
	return genericx509utils.GetSerialNumberAsBigInt(cert)
}

func GetSerialNumberAsHexColonSeparated(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSerialNumberAsHexColonSeparated(cert)
}

func GetCountryName(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetCountryName(cert)
}

func GetAdditionalSans(cert *x509.Certificate) ([]string, error) {
	return genericx509utils.GetAdditionalSans(cert)
}

func GetAdditionalSansOrEmptySliceIfUnset(cert *x509.Certificate) ([]string, error) {
	return genericx509utils.GetAdditionalSansOrEmptySliceIfUnset(cert)
}

func GetLocality(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetLocality(cert)
}

func GetOrganization(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetOrganization(cert)
}

func GetValidityDuration(cert *x509.Certificate) (*time.Duration, error) {
	return genericx509utils.GetValidityDuration(cert)
}

func GetValidityDurationAsString(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetValidityDurationAsString(cert)
}

func GetSubjectAndSerialString(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSubjectAndSerialString(cert)
}

func GetSubjectAsPkixName(cert *x509.Certificate) (pkix.Name, error) {
	return genericx509utils.GetSubjectAsPkixName(cert)
}

func GetSubjectCountryName(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSubjectCountryName(cert)
}

func GetSubjectLocalityName(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSubjectLocalityName(cert)
}

func GetSubjectOrganizationName(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSubjectOrganizationName(cert)
}

func GetSubjectStringForOpenssl(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetSubjectStringForOpenssl(cert)
}

func GetInfoSring(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetInfoString(cert)
}

func GetInfoString(cert *x509.Certificate) (string, error) {
	return genericx509utils.GetInfoString(cert)
}
