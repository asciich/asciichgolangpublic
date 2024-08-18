package asciichgolangpublic

import (
	"strings"
)

type X509CreateCertificateOptions struct {
	UseTemporaryDirectory bool

	// Certificate Attributes
	CommonName     string // the CN field
	CountryName    string // the C field
	Locality       string // the L field
	AdditionalSans []string

	KeyOutputFilePath         string
	CertificateOutputFilePath string

	IntermediateCertificateInGopass *GopassSecretOptions

	OverwriteExistingCertificateInGopass bool
	Verbose                              bool
}

func NewX509CreateCertificateOptions() (x *X509CreateCertificateOptions) {
	return new(X509CreateCertificateOptions)
}

func (o *X509CreateCertificateOptions) GetCertificateOutputFilePath() (certOutputPath string, err error) {
	if len(o.CertificateOutputFilePath) <= 0 {
		return "", TracedError("CertificateOutputFilePath not set")
	}

	return o.CertificateOutputFilePath, nil
}

func (o *X509CreateCertificateOptions) GetCommonName() (commonName string, err error) {
	commonName = strings.TrimSpace(o.CommonName)
	if len(commonName) <= 0 {
		return "", TracedError("commonName not set")
	}

	return commonName, nil
}

func (o *X509CreateCertificateOptions) GetCountryName() (countryName string, err error) {
	countryName = strings.TrimSpace(o.CountryName)
	if len(countryName) <= 0 {
		return "", TracedError("countryName not set")
	}

	return countryName, nil
}

func (o *X509CreateCertificateOptions) GetDeepCopy() (copy *X509CreateCertificateOptions) {
	copy = new(X509CreateCertificateOptions)

	*copy = *o

	return copy
}

func (o *X509CreateCertificateOptions) GetIntermediateCertificateGopassCredential() (certificate *GopassCredential, err error) {
	if o.IntermediateCertificateInGopass == nil {
		return nil, TracedError("IntermediateCertificateKeyInGopass not set")
	}

	optionsToUse := o.IntermediateCertificateInGopass.GetDeepCopy()
	optionsToUse.SecretBasename = "intermediateCertificate.crt"

	certificate, err = Gopass().GetCredential(optionsToUse)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func (o *X509CreateCertificateOptions) GetIntermediateCertificateKeyGopassCredential() (key *GopassCredential, err error) {
	if o.IntermediateCertificateInGopass == nil {
		return nil, TracedError("IntermediateCertificateKeyInGopass not set")
	}

	optionsToUse := o.IntermediateCertificateInGopass.GetDeepCopy()
	optionsToUse.SecretBasename = "intermediateCertificate.key"

	key, err = Gopass().GetCredential(optionsToUse)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (o *X509CreateCertificateOptions) GetKeyOutputFilePath() (keyOutputPath string, err error) {
	if len(o.KeyOutputFilePath) <= 0 {
		return "", TracedError("KeyOutputFilePath not set")
	}

	return o.KeyOutputFilePath, nil
}

func (o *X509CreateCertificateOptions) GetLocallity() (locality string, err error) {
	locality = strings.TrimSpace(o.Locality)
	if len(locality) <= 0 {
		return "", TracedError("locality not set")
	}

	return locality, nil
}

func (o *X509CreateCertificateOptions) GetSubjectStringForOpenssl() (subjectString string, err error) {
	subjectString = ""

	commonName, err := o.GetCommonName()
	if err != nil {
		return "", err
	}

	subjectString += "/CN=" + commonName

	countryName, err := o.GetCountryName()
	if err != nil {
		return "", err
	}

	subjectString += "/C=" + countryName

	locality, err := o.GetLocallity()
	if err != nil {
		return "", err
	}

	subjectString += "/L=" + locality
	if err != nil {
		return "", err
	}

	return subjectString, nil
}

func (o *X509CreateCertificateOptions) GetUseTemporaryDirectory() (UseTemporaryDirectory bool) {
	return o.UseTemporaryDirectory
}

func (o *X509CreateCertificateOptions) IsCertificateOutputFilePathSet() (isSet bool) {
	return len(o.CertificateOutputFilePath) > 0
}

func (x *X509CreateCertificateOptions) GetAdditionalSans() (additionalSans []string, err error) {
	if x.AdditionalSans == nil {
		return nil, TracedErrorf("AdditionalSans not set")
	}

	if len(x.AdditionalSans) <= 0 {
		return nil, TracedErrorf("AdditionalSans has no elements")
	}

	return x.AdditionalSans, nil
}

func (x *X509CreateCertificateOptions) GetIntermediateCertificateInGopass() (intermediateCertificateInGopass *GopassSecretOptions, err error) {
	if x.IntermediateCertificateInGopass == nil {
		return nil, TracedErrorf("IntermediateCertificateInGopass not set")
	}

	return x.IntermediateCertificateInGopass, nil
}

func (x *X509CreateCertificateOptions) GetLocality() (locality string, err error) {
	if x.Locality == "" {
		return "", TracedErrorf("Locality not set")
	}

	return x.Locality, nil
}

func (x *X509CreateCertificateOptions) GetOverwriteExistingCertificateInGopass() (overwriteExistingCertificateInGopass bool, err error) {

	return x.OverwriteExistingCertificateInGopass, nil
}

func (x *X509CreateCertificateOptions) GetVerbose() (verbose bool, err error) {

	return x.Verbose, nil
}

func (x *X509CreateCertificateOptions) MustGetAdditionalSans() (additionalSans []string) {
	additionalSans, err := x.GetAdditionalSans()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return additionalSans
}

func (x *X509CreateCertificateOptions) MustGetCertificateOutputFilePath() (certOutputPath string) {
	certOutputPath, err := x.GetCertificateOutputFilePath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return certOutputPath
}

func (x *X509CreateCertificateOptions) MustGetCommonName() (commonName string) {
	commonName, err := x.GetCommonName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commonName
}

func (x *X509CreateCertificateOptions) MustGetCountryName() (countryName string) {
	countryName, err := x.GetCountryName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return countryName
}

func (x *X509CreateCertificateOptions) MustGetIntermediateCertificateGopassCredential() (certificate *GopassCredential) {
	certificate, err := x.GetIntermediateCertificateGopassCredential()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return certificate
}

func (x *X509CreateCertificateOptions) MustGetIntermediateCertificateInGopass() (intermediateCertificateInGopass *GopassSecretOptions) {
	intermediateCertificateInGopass, err := x.GetIntermediateCertificateInGopass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return intermediateCertificateInGopass
}

func (x *X509CreateCertificateOptions) MustGetIntermediateCertificateKeyGopassCredential() (key *GopassCredential) {
	key, err := x.GetIntermediateCertificateKeyGopassCredential()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return key
}

func (x *X509CreateCertificateOptions) MustGetKeyOutputFilePath() (keyOutputPath string) {
	keyOutputPath, err := x.GetKeyOutputFilePath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyOutputPath
}

func (x *X509CreateCertificateOptions) MustGetLocality() (locality string) {
	locality, err := x.GetLocality()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return locality
}

func (x *X509CreateCertificateOptions) MustGetLocallity() (locality string) {
	locality, err := x.GetLocallity()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return locality
}

func (x *X509CreateCertificateOptions) MustGetOverwriteExistingCertificateInGopass() (overwriteExistingCertificateInGopass bool) {
	overwriteExistingCertificateInGopass, err := x.GetOverwriteExistingCertificateInGopass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return overwriteExistingCertificateInGopass
}

func (x *X509CreateCertificateOptions) MustGetSubjectStringForOpenssl() (subjectString string) {
	subjectString, err := x.GetSubjectStringForOpenssl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subjectString
}

func (x *X509CreateCertificateOptions) MustGetVerbose() (verbose bool) {
	verbose, err := x.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (x *X509CreateCertificateOptions) MustSetAdditionalSans(additionalSans []string) {
	err := x.SetAdditionalSans(additionalSans)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetCertificateOutputFilePath(certificateOutputFilePath string) {
	err := x.SetCertificateOutputFilePath(certificateOutputFilePath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetCommonName(commonName string) {
	err := x.SetCommonName(commonName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetCountryName(countryName string) {
	err := x.SetCountryName(countryName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetIntermediateCertificateInGopass(intermediateCertificateInGopass *GopassSecretOptions) {
	err := x.SetIntermediateCertificateInGopass(intermediateCertificateInGopass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetKeyOutputFilePath(keyOutputFilePath string) {
	err := x.SetKeyOutputFilePath(keyOutputFilePath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetLocality(locality string) {
	err := x.SetLocality(locality)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetOverwriteExistingCertificateInGopass(overwriteExistingCertificateInGopass bool) {
	err := x.SetOverwriteExistingCertificateInGopass(overwriteExistingCertificateInGopass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetUseTemporaryDirectory(useTemporaryDirectory bool) {
	err := x.SetUseTemporaryDirectory(useTemporaryDirectory)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetVerbose(verbose bool) {
	err := x.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) SetAdditionalSans(additionalSans []string) (err error) {
	if additionalSans == nil {
		return TracedErrorf("additionalSans is nil")
	}

	if len(additionalSans) <= 0 {
		return TracedErrorf("additionalSans has no elements")
	}

	x.AdditionalSans = additionalSans

	return nil
}

func (x *X509CreateCertificateOptions) SetCertificateOutputFilePath(certificateOutputFilePath string) (err error) {
	if certificateOutputFilePath == "" {
		return TracedErrorf("certificateOutputFilePath is empty string")
	}

	x.CertificateOutputFilePath = certificateOutputFilePath

	return nil
}

func (x *X509CreateCertificateOptions) SetCommonName(commonName string) (err error) {
	if commonName == "" {
		return TracedErrorf("commonName is empty string")
	}

	x.CommonName = commonName

	return nil
}

func (x *X509CreateCertificateOptions) SetCountryName(countryName string) (err error) {
	if countryName == "" {
		return TracedErrorf("countryName is empty string")
	}

	x.CountryName = countryName

	return nil
}

func (x *X509CreateCertificateOptions) SetIntermediateCertificateInGopass(intermediateCertificateInGopass *GopassSecretOptions) (err error) {
	if intermediateCertificateInGopass == nil {
		return TracedErrorf("intermediateCertificateInGopass is nil")
	}

	x.IntermediateCertificateInGopass = intermediateCertificateInGopass

	return nil
}

func (x *X509CreateCertificateOptions) SetKeyOutputFilePath(keyOutputFilePath string) (err error) {
	if keyOutputFilePath == "" {
		return TracedErrorf("keyOutputFilePath is empty string")
	}

	x.KeyOutputFilePath = keyOutputFilePath

	return nil
}

func (x *X509CreateCertificateOptions) SetLocality(locality string) (err error) {
	if locality == "" {
		return TracedErrorf("locality is empty string")
	}

	x.Locality = locality

	return nil
}

func (x *X509CreateCertificateOptions) SetOverwriteExistingCertificateInGopass(overwriteExistingCertificateInGopass bool) (err error) {
	x.OverwriteExistingCertificateInGopass = overwriteExistingCertificateInGopass

	return nil
}

func (x *X509CreateCertificateOptions) SetUseTemporaryDirectory(useTemporaryDirectory bool) (err error) {
	x.UseTemporaryDirectory = useTemporaryDirectory

	return nil
}

func (x *X509CreateCertificateOptions) SetVerbose(verbose bool) (err error) {
	x.Verbose = verbose

	return nil
}
