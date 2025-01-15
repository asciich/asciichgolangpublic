package asciichgolangpublic

import (
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type X509CertificateFile struct {
	File
}

func GetX509CertificateFileFromFile(input File) (x509CertificateFile *X509CertificateFile, err error) {
	if input == nil {
		return nil, tracederrors.TracedErrorNil("input")
	}

	fileToAdd := input.GetDeepCopy()

	x509CertificateFile = NewX509CertificateFile()

	x509CertificateFile.File = fileToAdd

	return x509CertificateFile, nil
}

func GetX509CertificateFileFromPath(inputPath string) (x509CertificateFile *X509CertificateFile, err error) {
	if inputPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("inputPath")
	}

	inputFile, err := GetLocalFileByPath(inputPath)
	if err != nil {
		return nil, err
	}

	x509CertificateFile, err = GetX509CertificateFileFromFile(inputFile)
	if err != nil {
		return nil, err
	}

	return x509CertificateFile, nil
}

func MustGetX509CertificateFileFromFile(input File) (x509CertificateFile *X509CertificateFile) {
	x509CertificateFile, err := GetX509CertificateFileFromFile(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return x509CertificateFile
}

func MustGetX509CertificateFileFromPath(inputPath string) (x509CertificateFile *X509CertificateFile) {
	x509CertificateFile, err := GetX509CertificateFileFromPath(inputPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return x509CertificateFile
}

func NewX509CertificateFile() (x *X509CertificateFile) {
	return new(X509CertificateFile)
}

func (x *X509CertificateFile) GetAsX509Certificate() (cert *X509Certificate, err error) {
	contentBytes, err := x.ReadAsBytes()
	if err != nil {
		return nil, err
	}

	cert = NewX509Certificate()
	err = cert.LoadFromBytes(contentBytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (x *X509CertificateFile) IsX509Certificate(verbose bool) (isX509Certificate bool, err error) {
	exists, err := x.Exists(verbose)
	if err != nil {
		return false, err
	}

	pathString, err := x.GetLocalPath()
	if err != nil {
		return false, err
	}

	if !exists {
		return false, tracederrors.TracedErrorf("file '%v' does not exist", pathString)
	}

	checkCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf(
			"openssl x509 -in '%v' -text &>/dev/null && echo yes || echo no",
			pathString,
		),
	}

	stdout, err := Bash().RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: checkCommand,
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout == "yes" {
		return true, nil
	}

	if stdout == "no" {
		return false, nil
	}

	return false, tracederrors.TracedErrorf(
		"Unable to check if '%v' contains a X509 certificate. Unexpected stdout: '%v'",
		pathString,
		stdout,
	)
}

func (x *X509CertificateFile) IsX509CertificateSignedByCertificateFile(signingCertificateFile File, verbose bool) (isSignedBy bool, err error) {
	if signingCertificateFile == nil {
		return false, tracederrors.TracedErrorNil("signingCertificateFile")
	}

	isSignedBy, err = X509Certificates().IsCertificateFileSignedByCertificateFile(x, signingCertificateFile, verbose)
	if err != nil {
		return false, err
	}

	return isSignedBy, nil
}

func (x *X509CertificateFile) IsX509IntermediateCertificate() (isIntermediateCertificate bool, err error) {
	cert, err := x.GetAsX509Certificate()
	if err != nil {
		return false, err
	}

	isIntermediateCertificate, err = cert.IsIntermediateCertificate()
	if err != nil {
		return false, err
	}

	return isIntermediateCertificate, nil
}

func (x *X509CertificateFile) IsX509RootCertificate(verbose bool) (isX509Certificate bool, err error) {
	cert, err := x.GetAsX509Certificate()
	if err != nil {
		return false, err
	}

	isX509Certificate, err = cert.IsRootCa(verbose)
	if err != nil {
		return false, err
	}

	return isX509Certificate, nil
}

func (x *X509CertificateFile) IsX509v3() (isX509v3 bool, err error) {
	cert, err := x.GetAsX509Certificate()
	if err != nil {
		return false, err
	}

	isX509v3, err = cert.IsV3()
	if err != nil {
		return false, err
	}

	return isX509v3, nil
}

func (x *X509CertificateFile) MustGetAsX509Certificate() (cert *X509Certificate) {
	cert, err := x.GetAsX509Certificate()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cert
}

func (x *X509CertificateFile) MustIsExpired(verbose bool) (isExpired bool, err error) {
	localPath, err := x.GetLocalPath()
	if err != nil {
		return false, err
	}

	certificate, err := x.GetAsX509Certificate()
	if err != nil {
		return false, err
	}

	isExpired, err = certificate.IsExpired()
	if err != nil {
		return false, err
	}

	if verbose {
		if isExpired {
			logging.LogInfof("X509Certificate in '%s' is expired.", localPath)
		} else {
			logging.LogInfof("X509Certificate in '%s' is NOT expired.", localPath)
		}
	}

	return isExpired, nil
}

func (x *X509CertificateFile) MustIsX509Certificate(verbose bool) (isX509Certificate bool) {
	isX509Certificate, err := x.IsX509Certificate(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isX509Certificate
}

func (x *X509CertificateFile) MustIsX509CertificateSignedByCertificateFile(signingCertificateFile File, verbose bool) (isSignedBy bool) {
	isSignedBy, err := x.IsX509CertificateSignedByCertificateFile(signingCertificateFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isSignedBy
}

func (x *X509CertificateFile) MustIsX509IntermediateCertificate() (isIntermediateCertificate bool) {
	isIntermediateCertificate, err := x.IsX509IntermediateCertificate()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isIntermediateCertificate
}

func (x *X509CertificateFile) MustIsX509RootCertificate(verbose bool) (isX509Certificate bool) {
	isX509Certificate, err := x.IsX509RootCertificate(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isX509Certificate
}

func (x *X509CertificateFile) MustIsX509v3() (isX509v3 bool) {
	isX509v3, err := x.IsX509v3()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isX509v3
}
