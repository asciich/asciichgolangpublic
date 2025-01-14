package asciichgolangpublic

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type X509SignCertificateOptions struct {
	CertFileUsedForSigning File
	KeyFileUsedForSigning  File
	KeyFileToSign          File
	OutputCertificateFile  File
	SigningRequestFile     File
	CommonName             string
	CountryName            string
	Locality               string
	Verbose                bool
}

func NewX509SignCertificateOptions() (deepCopy *X509SignCertificateOptions) {
	return new(X509SignCertificateOptions)
}

func (o *X509SignCertificateOptions) GetCertFileUsedForSigning() (keyFileForSigning File, err error) {
	if o.CertFileUsedForSigning == nil {
		return nil, errors.TracedError("CertFileUsedForSigning not set")
	}

	return o.CertFileUsedForSigning, nil
}

func (o *X509SignCertificateOptions) GetCommonName() (commonName string, err error) {
	if len(o.CommonName) <= 0 {
		return "", errors.TracedError("CommonName is not set")
	}

	return o.CommonName, nil
}

func (o *X509SignCertificateOptions) GetCountryName() (countryName string, err error) {
	if len(o.CountryName) <= 0 {
		return "", errors.TracedError("CountryName is not set")
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

func (o *X509SignCertificateOptions) GetKeyFileToSign() (keyFileForSigning File, err error) {
	if o.KeyFileToSign == nil {
		return nil, errors.TracedError("KeyFileToSign not set")
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

func (o *X509SignCertificateOptions) GetKeyFileUsedForSigning() (keyFileForSigning File, err error) {
	if o.KeyFileUsedForSigning == nil {
		return nil, errors.TracedError("KeyFileUsedForSigning not set")
	}

	return o.KeyFileUsedForSigning, nil
}

func (o *X509SignCertificateOptions) GetLocality() (locality string, err error) {
	if len(o.Locality) <= 0 {
		return "", errors.TracedError("Locality is not set")
	}

	return o.CommonName, nil
}

func (o *X509SignCertificateOptions) GetOutputCertificateFile() (keyFileForSigning File, err error) {
	if o.OutputCertificateFile == nil {
		return nil, errors.TracedError("OutputCertificateFile not set")
	}

	return o.OutputCertificateFile, nil
}

func (o *X509SignCertificateOptions) GetSigningRequestFile() (signingRequestFile File, err error) {
	if o.SigningRequestFile == nil {
		return nil, errors.TracedError("SigningRequestFile is not set")
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

func (x *X509SignCertificateOptions) MustGetCertFileUsedForSigning() (keyFileForSigning File) {
	keyFileForSigning, err := x.GetCertFileUsedForSigning()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyFileForSigning
}

func (x *X509SignCertificateOptions) MustGetCommonName() (commonName string) {
	commonName, err := x.GetCommonName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commonName
}

func (x *X509SignCertificateOptions) MustGetCountryName() (countryName string) {
	countryName, err := x.GetCountryName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return countryName
}

func (x *X509SignCertificateOptions) MustGetKeyFileToSign() (keyFileForSigning File) {
	keyFileForSigning, err := x.GetKeyFileToSign()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyFileForSigning
}

func (x *X509SignCertificateOptions) MustGetKeyFileToSignPath() (keyFileForSigningPath string) {
	keyFileForSigningPath, err := x.GetKeyFileToSignPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyFileForSigningPath
}

func (x *X509SignCertificateOptions) MustGetKeyFileUsedForSigning() (keyFileForSigning File) {
	keyFileForSigning, err := x.GetKeyFileUsedForSigning()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyFileForSigning
}

func (x *X509SignCertificateOptions) MustGetLocality() (locality string) {
	locality, err := x.GetLocality()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return locality
}

func (x *X509SignCertificateOptions) MustGetOutputCertificateFile() (keyFileForSigning File) {
	keyFileForSigning, err := x.GetOutputCertificateFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyFileForSigning
}

func (x *X509SignCertificateOptions) MustGetSigningRequestFile() (signingRequestFile File) {
	signingRequestFile, err := x.GetSigningRequestFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return signingRequestFile
}

func (x *X509SignCertificateOptions) MustGetSigningRequestFilePath() (signingRequestFilePath string) {
	signingRequestFilePath, err := x.GetSigningRequestFilePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return signingRequestFilePath
}

func (x *X509SignCertificateOptions) MustGetSubjectToSign() (subjectToSign string) {
	subjectToSign, err := x.GetSubjectToSign()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subjectToSign
}

func (x *X509SignCertificateOptions) MustGetVerbose() (verbose bool) {
	verbose, err := x.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (x *X509SignCertificateOptions) MustSetCertFileUsedForSigning(certFileUsedForSigning File) {
	err := x.SetCertFileUsedForSigning(certFileUsedForSigning)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetCommonName(commonName string) {
	err := x.SetCommonName(commonName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetCountryName(countryName string) {
	err := x.SetCountryName(countryName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetKeyFileToSign(keyFileToSign File) {
	err := x.SetKeyFileToSign(keyFileToSign)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetKeyFileUsedForSigning(keyFileUsedForSigning File) {
	err := x.SetKeyFileUsedForSigning(keyFileUsedForSigning)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetLocality(locality string) {
	err := x.SetLocality(locality)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetOutputCertificateFile(outputCertificateFile File) {
	err := x.SetOutputCertificateFile(outputCertificateFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetSigningRequestFile(signingRequestFile File) {
	err := x.SetSigningRequestFile(signingRequestFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) MustSetVerbose(verbose bool) {
	err := x.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509SignCertificateOptions) SetCertFileUsedForSigning(certFileUsedForSigning File) (err error) {
	if certFileUsedForSigning == nil {
		return errors.TracedErrorf("certFileUsedForSigning is nil")
	}

	x.CertFileUsedForSigning = certFileUsedForSigning

	return nil
}

func (x *X509SignCertificateOptions) SetCommonName(commonName string) (err error) {
	if commonName == "" {
		return errors.TracedErrorf("commonName is empty string")
	}

	x.CommonName = commonName

	return nil
}

func (x *X509SignCertificateOptions) SetCountryName(countryName string) (err error) {
	if countryName == "" {
		return errors.TracedErrorf("countryName is empty string")
	}

	x.CountryName = countryName

	return nil
}

func (x *X509SignCertificateOptions) SetKeyFileToSign(keyFileToSign File) (err error) {
	if keyFileToSign == nil {
		return errors.TracedErrorf("keyFileToSign is nil")
	}

	x.KeyFileToSign = keyFileToSign

	return nil
}

func (x *X509SignCertificateOptions) SetKeyFileUsedForSigning(keyFileUsedForSigning File) (err error) {
	if keyFileUsedForSigning == nil {
		return errors.TracedErrorf("keyFileUsedForSigning is nil")
	}

	x.KeyFileUsedForSigning = keyFileUsedForSigning

	return nil
}

func (x *X509SignCertificateOptions) SetLocality(locality string) (err error) {
	if locality == "" {
		return errors.TracedErrorf("locality is empty string")
	}

	x.Locality = locality

	return nil
}

func (x *X509SignCertificateOptions) SetOutputCertificateFile(outputCertificateFile File) (err error) {
	if outputCertificateFile == nil {
		return errors.TracedErrorf("outputCertificateFile is nil")
	}

	x.OutputCertificateFile = outputCertificateFile

	return nil
}

func (x *X509SignCertificateOptions) SetSigningRequestFile(signingRequestFile File) (err error) {
	if signingRequestFile == nil {
		return errors.TracedErrorf("signingRequestFile is nil")
	}

	x.SigningRequestFile = signingRequestFile

	return nil
}

func (x *X509SignCertificateOptions) SetVerbose(verbose bool) (err error) {
	x.Verbose = verbose

	return nil
}
