package x509utils

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type X509SignCertificateOptions struct {
	CertFileUsedForSigning filesinterfaces.File
	KeyFileUsedForSigning  filesinterfaces.File
	KeyFileToSign          filesinterfaces.File
	OutputCertificateFile  filesinterfaces.File
	SigningRequestFile     filesinterfaces.File
	CommonName             string
	CountryName            string
	Locality               string
	Verbose                bool
}

func NewX509SignCertificateOptions() (deepCopy *X509SignCertificateOptions) {
	return new(X509SignCertificateOptions)
}

func (o *X509SignCertificateOptions) GetCertFileUsedForSigning() (keyFileForSigning filesinterfaces.File, err error) {
	if o.CertFileUsedForSigning == nil {
		return nil, tracederrors.TracedError("CertFileUsedForSigning not set")
	}

	return o.CertFileUsedForSigning, nil
}

func (o *X509SignCertificateOptions) GetCommonName() (commonName string, err error) {
	if len(o.CommonName) <= 0 {
		return "", tracederrors.TracedError("CommonName is not set")
	}

	return o.CommonName, nil
}

func (o *X509SignCertificateOptions) GetCountryName() (countryName string, err error) {
	if len(o.CountryName) <= 0 {
		return "", tracederrors.TracedError("CountryName is not set")
	}

	return o.CountryName, nil
}

func (o *X509SignCertificateOptions) GetDeepCopy() (deepCopy *X509SignCertificateOptions) {
	deepCopy = NewX509SignCertificateOptions()

	*deepCopy = *o
	if o.CertFileUsedForSigning != nil {
		deepCopy.CertFileUsedForSigning = o.CertFileUsedForSigning.GetDeepCopy()
	}

	if o.KeyFileUsedForSigning != nil {
		deepCopy.KeyFileUsedForSigning = o.KeyFileUsedForSigning.GetDeepCopy()
	}

	if o.KeyFileToSign != nil {
		deepCopy.KeyFileToSign = o.KeyFileToSign.GetDeepCopy()
	}

	if o.OutputCertificateFile != nil {
		deepCopy.OutputCertificateFile = o.OutputCertificateFile.GetDeepCopy()
	}

	if o.SigningRequestFile != nil {
		deepCopy.SigningRequestFile = o.SigningRequestFile.GetDeepCopy()
	}

	return deepCopy
}

func (o *X509SignCertificateOptions) GetKeyFileToSign() (keyFileForSigning filesinterfaces.File, err error) {
	if o.KeyFileToSign == nil {
		return nil, tracederrors.TracedError("KeyFileToSign not set")
	}

	return o.KeyFileToSign, nil
}

func (o *X509SignCertificateOptions) GetKeyFileToSignPath() (keyFileForSigningPath string, err error) {
	keyFile, err := o.GetKeyFileToSign()
	if err != nil {
		return "", err
	}

	keyFileForSigningPath, err = keyFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return keyFileForSigningPath, nil
}

func (o *X509SignCertificateOptions) GetKeyFileUsedForSigning() (keyFileForSigning filesinterfaces.File, err error) {
	if o.KeyFileUsedForSigning == nil {
		return nil, tracederrors.TracedError("KeyFileUsedForSigning not set")
	}

	return o.KeyFileUsedForSigning, nil
}

func (o *X509SignCertificateOptions) GetLocality() (locality string, err error) {
	if len(o.Locality) <= 0 {
		return "", tracederrors.TracedError("Locality is not set")
	}

	return o.CommonName, nil
}

func (o *X509SignCertificateOptions) GetOutputCertificateFile() (keyFileForSigning filesinterfaces.File, err error) {
	if o.OutputCertificateFile == nil {
		return nil, tracederrors.TracedError("OutputCertificateFile not set")
	}

	return o.OutputCertificateFile, nil
}

func (o *X509SignCertificateOptions) GetSigningRequestFile() (signingRequestFile filesinterfaces.File, err error) {
	if o.SigningRequestFile == nil {
		return nil, tracederrors.TracedError("SigningRequestFile is not set")
	}

	return o.SigningRequestFile, nil
}

func (o *X509SignCertificateOptions) GetSigningRequestFilePath() (signingRequestFilePath string, err error) {
	signingRequestFile, err := o.GetSigningRequestFile()
	if err != nil {
		return "", err
	}

	signingRequestFilePath, err = signingRequestFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return signingRequestFilePath, nil
}

func (o *X509SignCertificateOptions) GetSubjectToSign() (subjectToSign string, err error) {
	countryName, err := o.GetCountryName()
	if err != nil {
		return "", err
	}

	locality, err := o.GetLocality()
	if err != nil {
		return "", err
	}

	subjectToSign = fmt.Sprintf("/C=%s/L=%s", countryName, locality)
	return subjectToSign, nil
}

func (x *X509SignCertificateOptions) GetVerbose() (verbose bool, err error) {

	return x.Verbose, nil
}

func (x *X509SignCertificateOptions) SetCertFileUsedForSigning(certFileUsedForSigning filesinterfaces.File) (err error) {
	if certFileUsedForSigning == nil {
		return tracederrors.TracedErrorf("certFileUsedForSigning is nil")
	}

	x.CertFileUsedForSigning = certFileUsedForSigning

	return nil
}

func (x *X509SignCertificateOptions) SetCommonName(commonName string) (err error) {
	if commonName == "" {
		return tracederrors.TracedErrorf("commonName is empty string")
	}

	x.CommonName = commonName

	return nil
}

func (x *X509SignCertificateOptions) SetCountryName(countryName string) (err error) {
	if countryName == "" {
		return tracederrors.TracedErrorf("countryName is empty string")
	}

	x.CountryName = countryName

	return nil
}

func (x *X509SignCertificateOptions) SetKeyFileToSign(keyFileToSign filesinterfaces.File) (err error) {
	if keyFileToSign == nil {
		return tracederrors.TracedErrorf("keyFileToSign is nil")
	}

	x.KeyFileToSign = keyFileToSign

	return nil
}

func (x *X509SignCertificateOptions) SetKeyFileUsedForSigning(keyFileUsedForSigning filesinterfaces.File) (err error) {
	if keyFileUsedForSigning == nil {
		return tracederrors.TracedErrorf("keyFileUsedForSigning is nil")
	}

	x.KeyFileUsedForSigning = keyFileUsedForSigning

	return nil
}

func (x *X509SignCertificateOptions) SetLocality(locality string) (err error) {
	if locality == "" {
		return tracederrors.TracedErrorf("locality is empty string")
	}

	x.Locality = locality

	return nil
}

func (x *X509SignCertificateOptions) SetOutputCertificateFile(outputCertificateFile filesinterfaces.File) (err error) {
	if outputCertificateFile == nil {
		return tracederrors.TracedErrorf("outputCertificateFile is nil")
	}

	x.OutputCertificateFile = outputCertificateFile

	return nil
}

func (x *X509SignCertificateOptions) SetSigningRequestFile(signingRequestFile filesinterfaces.File) (err error) {
	if signingRequestFile == nil {
		return tracederrors.TracedErrorf("signingRequestFile is nil")
	}

	x.SigningRequestFile = signingRequestFile

	return nil
}

func (x *X509SignCertificateOptions) SetVerbose(verbose bool) (err error) {
	x.Verbose = verbose

	return nil
}
