package x509utils

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"time"

	"crypto/x509"
	"encoding/pem"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// deprecated: TODO remove this one since obsolete....
type X509Certificate struct {
	nativeX509Certificate *x509.Certificate
}

func GetX509CertificateFromFile(certFile filesinterfaces.File) (cert *X509Certificate, err error) {
	if certFile == nil {
		return nil, err
	}

	path, err := certFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	cert, err = GetX509CertificateFromFilePath(path)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func GetX509CertificateFromFilePath(certFilePath string) (cert *X509Certificate, err error) {
	certFilePath = strings.TrimSpace(certFilePath)
	if len(certFilePath) <= 0 {
		return nil, tracederrors.TracedError("certPAth is empty string")
	}

	cert = NewX509Certificate()
	err = cert.LoadFromFilePath(certFilePath)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func MustGetX509CertificateFromFile(certFile filesinterfaces.File) (cert *X509Certificate) {
	cert, err := GetX509CertificateFromFile(certFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cert
}

func MustGetX509CertificateFromFilePath(certFilePath string) (cert *X509Certificate) {
	cert, err := GetX509CertificateFromFilePath(certFilePath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cert
}

func NewX509Certificate() (cert *X509Certificate) {
	return new(X509Certificate)
}

// TODO replace by EncodeCertificateAsPEMString
func (c *X509Certificate) GetAsPemBytes() (pemBytes []byte, err error) {
	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	var pemCert = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: nativeCert.Raw,
	}
	err = pem.Encode(io.Writer(&buf), pemCert)
	if err != nil {
		return nil, err
	}

	pemBytes = buf.Bytes()
	const minLen = 50
	if len(pemBytes) < minLen {
		return nil, tracederrors.TracedErrorf("pemBytes has less than '%v' bytes which is not enough for a pem certificate", minLen)
	}

	return pemBytes, nil
}

func (c *X509Certificate) GetAsPemString() (pemString string, err error) {
	pemBytes, err := c.GetAsPemBytes()
	if err != nil {
		return "", err
	}

	return string(pemBytes), nil
}

func (c *X509Certificate) GetIssuerString() (issuerString string, err error) {
	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return "", err
	}

	issuerString = nativeCert.Issuer.String()

	return issuerString, nil
}

func (c *X509Certificate) GetNativeCertificate() (nativeCertificate *x509.Certificate, err error) {
	if c.nativeX509Certificate == nil {
		return nil, tracederrors.TracedError("native certificate not set. Is the certificate Loaded?")
	}

	return c.nativeX509Certificate, nil
}

func (c *X509Certificate) GetSubjectString() (subject string, err error) {
	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return "", err
	}

	subject = nativeCert.Subject.String()
	return subject, nil
}

func (c *X509Certificate) GetVersion() (version int, err error) {
	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return -1, err
	}

	version = nativeCert.Version

	return version, nil
}

func (c *X509Certificate) IsIntermediateCertificate() (isIntermediateCertificate bool, err error) {
	isV1, err := c.IsV1()
	if err != nil {
		return false, err
	}

	if isV1 {
		return false, tracederrors.TracedError("v1 certificates are not supported as intermediate certificate.")
	}

	subjectString, err := c.GetSubjectString()
	if err != nil {
		return false, err
	}

	issuerString, err := c.GetIssuerString()
	if err != nil {
		return false, err
	}

	subjectsEqual := subjectString == issuerString
	if subjectsEqual {
		return false, nil
	}

	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return false, err
	}

	isIntermediateCertificate = nativeCert.IsCA

	return isIntermediateCertificate, nil
}

func (c *X509Certificate) IsRootCa(verbose bool) (isRootCa bool, err error) {
	isV1, err := c.IsV1()
	if err != nil {
		return false, err
	}

	if isV1 {
		if verbose {
			logging.LogWarn("v1 certificates are not supported as root ca.")
		}
		return false, err
	}

	subjectString, err := c.GetSubjectString()
	if err != nil {
		return false, err
	}

	issuerString, err := c.GetIssuerString()
	if err != nil {
		return false, err
	}

	subjectsEqual := subjectString == issuerString
	if !subjectsEqual {
		return false, nil
	}

	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return false, err
	}

	isRootCa = nativeCert.IsCA

	return isRootCa, nil
}

func (c *X509Certificate) IsSignedByCertificateFile(signingCertificate filesinterfaces.File, verbose bool) (isSignedBy bool, err error) {
	if signingCertificate == nil {
		return false, tracederrors.TracedError("signingCertificate is nil")
	}

	rootCert, err := GetX509CertificateFromFile(signingCertificate)
	if err != nil {
		return false, err
	}

	rootCertPemBytes, err := rootCert.GetAsPemBytes()
	if err != nil {
		return false, err
	}

	nativeCert, err := c.GetNativeCertificate()
	if err != nil {
		return false, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCertPemBytes)
	if !ok {
		return false, tracederrors.TracedError("Unable to add cert to root pool")
	}

	verifyOptions := x509.VerifyOptions{
		Roots: roots,
	}

	_, err = nativeCert.Verify(verifyOptions)
	if err != nil {
		if errors.As(err, &x509.UnknownAuthorityError{}) {
			isSignedBy = false
		} else {
			return false, tracederrors.TracedErrorf("Unable to verify certificate: '%v'", err)
		}
	} else {
		isSignedBy = true
	}

	if verbose {
		rootCertPath, err := signingCertificate.GetLocalPath()
		if err != nil {
			return false, err
		}

		rootCertSubject, err := rootCert.GetSubjectString()
		if err != nil {
			return false, err
		}

		toCheckSubject, err := c.GetSubjectString()
		if err != nil {
			return false, err
		}

		if isSignedBy {
			logging.LogInfof("The certificate '%v' is signed by '%v' read from file '%v'.", toCheckSubject, rootCertSubject, rootCertPath)
		} else {
			logging.LogInfof("The certificate '%v' is NOT signed by '%v' read from file '%v'.", toCheckSubject, rootCertSubject, rootCertPath)
		}
	}

	return isSignedBy, nil
}

func (c *X509Certificate) IsV1() (isV1 bool, err error) {
	version, err := c.GetVersion()
	if err != nil {
		return false, err
	}

	isV1 = version == 1

	return isV1, nil
}

func (c *X509Certificate) IsV3() (isV3 bool, err error) {
	version, err := c.GetVersion()
	if err != nil {
		return false, err
	}

	isV3 = version == 3

	return isV3, nil
}

func (c *X509Certificate) LoadFromBytes(certBytes []byte) (err error) {
	if len(certBytes) <= 0 {
		return tracederrors.TracedError("certBytes has len 0")
	}

	block, _ := pem.Decode(certBytes)
	var blockBytes []byte = nil
	if block == nil {
		return tracederrors.TracedError("Failed to parse certificate PEM")
	} else {
		if block.Bytes == nil {
			return tracederrors.TracedError("Decode returned block.Bytes as nil")
		} else {
			blockBytes = block.Bytes
		}
	}

	var nativeCert *x509.Certificate = nil
	if blockBytes == nil {
		return tracederrors.TracedError("blockBytes is nil after evaluation")
	} else {
		nativeCert, err = x509.ParseCertificate(blockBytes)
		if err != nil {
			return tracederrors.TracedError("failed to parse certificate: " + err.Error())
		}
	}

	if nativeCert == nil {
		return tracederrors.TracedError("nativeCert is nil after evaluation")
	} else {
		c.nativeX509Certificate = nativeCert
	}

	return nil
}

func (c *X509Certificate) LoadFromFile(loadFile filesinterfaces.File) (err error) {
	if loadFile == nil {
		return tracederrors.TracedError("loadFile is nil")
	}

	contentString, err := loadFile.ReadAsString()
	if err != nil {
		return err
	}

	err = c.LoadFromString(contentString)
	if err != nil {
		return err
	}

	return nil
}

func (c *X509Certificate) LoadFromFilePath(loadPath string) (err error) {
	loadPath = strings.TrimSpace(loadPath)
	if len(loadPath) <= 0 {
		return tracederrors.TracedError("loadPath is empty string")
	}

	loadFile, err := files.GetLocalFileByPath(loadPath)
	if err != nil {
		return err
	}

	err = c.LoadFromFile(loadFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *X509Certificate) LoadFromString(certString string) (err error) {
	if len(certString) <= 0 {
		return tracederrors.TracedError("certString is empty string")
	}

	err = c.LoadFromBytes([]byte(certString))
	if err != nil {
		return err
	}

	return nil
}

func (c *X509Certificate) WritePemToFile(outputFile filesinterfaces.File, verbose bool) (err error) {
	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	pemBytes, err := c.GetAsPemBytes()
	if err != nil {
		return err
	}

	err = outputFile.WriteBytes(pemBytes, verbose)
	if err != nil {
		return err
	}

	path, err := outputFile.GetLocalPath()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf("Wrote certificate as PEM into '%s'", path)
	}

	return nil
}

func (c *X509Certificate) WritePemToFilePath(filePath string, verbose bool) (err error) {
	outFile, err := files.GetLocalFileByPath(filePath)
	if err != nil {
		return err
	}

	err = c.WritePemToFile(outFile, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (x *X509Certificate) GetExpiryDate() (expiryDate *time.Time, err error) {
	nativeCert, err := x.GetNativeCertificate()
	if err != nil {
		return nil, err
	}

	expiryDate = new(time.Time)
	*expiryDate = nativeCert.NotAfter

	return expiryDate, nil
}

func (x *X509Certificate) GetNativeX509Certificate() (nativeX509Certificate *x509.Certificate, err error) {
	if x.nativeX509Certificate == nil {
		return nil, tracederrors.TracedErrorf("nativeX509Certificate not set")
	}

	return x.nativeX509Certificate, nil
}

func (x *X509Certificate) IsExpired() (isExpired bool, err error) {
	expiryDate, err := x.GetExpiryDate()
	if err != nil {
		return false, err
	}

	isExpired = time.Now().After(*expiryDate)

	return isExpired, nil
}

func (x *X509Certificate) MustGetAsPemBytes() (pemBytes []byte) {
	pemBytes, err := x.GetAsPemBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return pemBytes
}

func (x *X509Certificate) MustGetAsPemString() (pemString string) {
	pemString, err := x.GetAsPemString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return pemString
}

func (x *X509Certificate) MustGetExpiryDate() (expiryDate *time.Time) {
	expiryDate, err := x.GetExpiryDate()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return expiryDate
}

func (x *X509Certificate) MustGetIssuerString() (issuerString string) {
	issuerString, err := x.GetIssuerString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return issuerString
}

func (x *X509Certificate) MustGetNativeCertificate() (nativeCertificate *x509.Certificate) {
	nativeCertificate, err := x.GetNativeCertificate()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeCertificate
}

func (x *X509Certificate) MustGetNativeX509Certificate() (nativeX509Certificate *x509.Certificate) {
	nativeX509Certificate, err := x.GetNativeX509Certificate()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeX509Certificate
}

func (x *X509Certificate) MustGetSubjectString() (subject string) {
	subject, err := x.GetSubjectString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subject
}

func (x *X509Certificate) MustGetVersion() (version int) {
	version, err := x.GetVersion()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func (x *X509Certificate) MustIsExpired() (isExpired bool) {
	isExpired, err := x.IsExpired()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isExpired
}

func (x *X509Certificate) MustIsIntermediateCertificate() (isIntermediateCertificate bool) {
	isIntermediateCertificate, err := x.IsIntermediateCertificate()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isIntermediateCertificate
}

func (x *X509Certificate) MustIsRootCa(verbose bool) (isRootCa bool) {
	isRootCa, err := x.IsRootCa(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRootCa
}

func (x *X509Certificate) MustIsSignedByCertificateFile(signingCertificate filesinterfaces.File, verbose bool) (isSignedBy bool) {
	isSignedBy, err := x.IsSignedByCertificateFile(signingCertificate, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isSignedBy
}

func (x *X509Certificate) MustIsV1() (isV1 bool) {
	isV1, err := x.IsV1()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isV1
}

func (x *X509Certificate) MustIsV3() (isV3 bool) {
	isV3, err := x.IsV3()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isV3
}

func (x *X509Certificate) MustLoadFromBytes(certBytes []byte) {
	err := x.LoadFromBytes(certBytes)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) MustLoadFromFile(loadFile filesinterfaces.File) {
	err := x.LoadFromFile(loadFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) MustLoadFromFilePath(loadPath string) {
	err := x.LoadFromFilePath(loadPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) MustLoadFromString(certString string) {
	err := x.LoadFromString(certString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) MustSetNativeX509Certificate(nativeX509Certificate *x509.Certificate) {
	err := x.SetNativeX509Certificate(nativeX509Certificate)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) MustWritePemToFile(outputFile filesinterfaces.File, verbose bool) {
	err := x.WritePemToFile(outputFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) MustWritePemToFilePath(filePath string, verbose bool) {
	err := x.WritePemToFilePath(filePath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509Certificate) SetNativeX509Certificate(nativeX509Certificate *x509.Certificate) (err error) {
	if nativeX509Certificate == nil {
		return tracederrors.TracedErrorf("nativeX509Certificate is nil")
	}

	x.nativeX509Certificate = nativeX509Certificate

	return nil
}
