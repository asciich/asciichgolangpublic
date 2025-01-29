package x509utils

import (
	"crypto/x509/pkix"
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type X509CreateCertificateOptions struct {
	UseTemporaryDirectory bool

	// Certificate Attributes
	CommonName     string // the CN field
	CountryName    string // the C field
	Organization   string
	Locality       string // the L field
	AdditionalSans []string

	KeyOutputFilePath         string
	CertificateOutputFilePath string

	IntermediateCertificateInGopass *parameteroptions.GopassSecretOptions // TODO move to Gopass

	OverwriteExistingCertificateInGopass bool
	Verbose                              bool
}

func NewX509CreateCertificateOptions() (x *X509CreateCertificateOptions) {
	return new(X509CreateCertificateOptions)
}

func (o *X509CreateCertificateOptions) GetSubjectAsPkixName() (subject *pkix.Name, err error) {
	countryName, err := o.GetCountryName()
	if err != nil {
		return nil, err
	}

	locality, err := o.GetLocality()
	if err != nil {
		return nil, err
	}

	organization, err := o.GetOrganization()
	if err != nil {
		return nil, err
	}

	subject = &pkix.Name{
		Organization: []string{organization},
		Country:      []string{countryName},
		Province:     []string{""},
		Locality:     []string{locality},
	}

	return subject, nil
}

func (o *X509CreateCertificateOptions) GetOrganization() (organization string, err error) {
	if o.Organization == "" {
		return "", tracederrors.TracedError("Organization not set")
	}

	return o.Organization, nil
}

func (o *X509CreateCertificateOptions) GetCertificateOutputFilePath() (certOutputPath string, err error) {
	if len(o.CertificateOutputFilePath) <= 0 {
		return "", tracederrors.TracedError("CertificateOutputFilePath not set")
	}

	return o.CertificateOutputFilePath, nil
}

func (o *X509CreateCertificateOptions) GetCommonName() (commonName string, err error) {
	commonName = strings.TrimSpace(o.CommonName)
	if len(commonName) <= 0 {
		return "", tracederrors.TracedError("commonName not set")
	}

	return commonName, nil
}

func (o *X509CreateCertificateOptions) GetCountryName() (countryName string, err error) {
	countryName = strings.TrimSpace(o.CountryName)
	if len(countryName) <= 0 {
		return "", tracederrors.TracedError("countryName not set")
	}

	return countryName, nil
}

func (o *X509CreateCertificateOptions) GetDeepCopy() (copy *X509CreateCertificateOptions) {
	copy = new(X509CreateCertificateOptions)

	*copy = *o

	return copy
}

/* TODO move to gopass
func (o *X509CreateCertificateOptions) GetIntermediateCertificateGopassCredential() (certificate *GopassCredential, err error) {
	if o.IntermediateCertificateInGopass == nil {
		return nil, tracederrors.TracedError("IntermediateCertificateKeyInGopass not set")
	}

	optionsToUse := o.IntermediateCertificateInGopass.GetDeepCopy()
	optionsToUse.SecretBasename = "intermediateCertificate.crt"

	certificate, err = Gopass().GetCredential(optionsToUse)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}
*/

/* TODO move to gopass
func (o *X509CreateCertificateOptions) GetIntermediateCertificateKeyGopassCredential() (key *GopassCredential, err error) {
	if o.IntermediateCertificateInGopass == nil {
		return nil, tracederrors.TracedError("IntermediateCertificateKeyInGopass not set")
	}

	optionsToUse := o.IntermediateCertificateInGopass.GetDeepCopy()
	optionsToUse.SecretBasename = "intermediateCertificate.key"

	key, err = Gopass().GetCredential(optionsToUse)
	if err != nil {
		return nil, err
	}

	return key, nil
}
*/

func (o *X509CreateCertificateOptions) GetKeyOutputFilePath() (keyOutputPath string, err error) {
	if len(o.KeyOutputFilePath) <= 0 {
		return "", tracederrors.TracedError("KeyOutputFilePath not set")
	}

	return o.KeyOutputFilePath, nil
}

func (o *X509CreateCertificateOptions) GetLocallity() (locality string, err error) {
	locality = strings.TrimSpace(o.Locality)
	if len(locality) <= 0 {
		return "", tracederrors.TracedError("locality not set")
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
		return nil, tracederrors.TracedErrorf("AdditionalSans not set")
	}

	if len(x.AdditionalSans) <= 0 {
		return nil, tracederrors.TracedErrorf("AdditionalSans has no elements")
	}

	return x.AdditionalSans, nil
}

/*
func (x *X509CreateCertificateOptions) GetIntermediateCertificateInGopass() (intermediateCertificateInGopass *GopassSecretOptions, err error) {
	if x.IntermediateCertificateInGopass == nil {
		return nil, tracederrors.TracedErrorf("IntermediateCertificateInGopass not set")
	}

	return x.IntermediateCertificateInGopass, nil
}
*/

func (x *X509CreateCertificateOptions) GetLocality() (locality string, err error) {
	if x.Locality == "" {
		return "", tracederrors.TracedErrorf("Locality not set")
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
		logging.LogGoErrorFatal(err)
	}

	return additionalSans
}

func (x *X509CreateCertificateOptions) MustGetCertificateOutputFilePath() (certOutputPath string) {
	certOutputPath, err := x.GetCertificateOutputFilePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return certOutputPath
}

func (x *X509CreateCertificateOptions) MustGetCommonName() (commonName string) {
	commonName, err := x.GetCommonName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commonName
}

func (x *X509CreateCertificateOptions) MustGetCountryName() (countryName string) {
	countryName, err := x.GetCountryName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return countryName
}

/*
func (x *X509CreateCertificateOptions) MustGetIntermediateCertificateGopassCredential() (certificate *GopassCredential) {
	certificate, err := x.GetIntermediateCertificateGopassCredential()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return certificate
}
*/

/*
func (x *X509CreateCertificateOptions) MustGetIntermediateCertificateInGopass() (intermediateCertificateInGopass *GopassSecretOptions) {
	intermediateCertificateInGopass, err := x.GetIntermediateCertificateInGopass()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return intermediateCertificateInGopass
}
*/

/*
func (x *X509CreateCertificateOptions) MustGetIntermediateCertificateKeyGopassCredential() (key *GopassCredential) {
	key, err := x.GetIntermediateCertificateKeyGopassCredential()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return key
}
*/

func (x *X509CreateCertificateOptions) MustGetKeyOutputFilePath() (keyOutputPath string) {
	keyOutputPath, err := x.GetKeyOutputFilePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyOutputPath
}

func (x *X509CreateCertificateOptions) MustGetLocality() (locality string) {
	locality, err := x.GetLocality()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return locality
}

func (x *X509CreateCertificateOptions) MustGetLocallity() (locality string) {
	locality, err := x.GetLocallity()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return locality
}

func (x *X509CreateCertificateOptions) MustGetOverwriteExistingCertificateInGopass() (overwriteExistingCertificateInGopass bool) {
	overwriteExistingCertificateInGopass, err := x.GetOverwriteExistingCertificateInGopass()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return overwriteExistingCertificateInGopass
}

func (x *X509CreateCertificateOptions) MustGetSubjectStringForOpenssl() (subjectString string) {
	subjectString, err := x.GetSubjectStringForOpenssl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subjectString
}

func (x *X509CreateCertificateOptions) MustGetVerbose() (verbose bool) {
	verbose, err := x.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (x *X509CreateCertificateOptions) MustSetAdditionalSans(additionalSans []string) {
	err := x.SetAdditionalSans(additionalSans)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetCertificateOutputFilePath(certificateOutputFilePath string) {
	err := x.SetCertificateOutputFilePath(certificateOutputFilePath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetCommonName(commonName string) {
	err := x.SetCommonName(commonName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetCountryName(countryName string) {
	err := x.SetCountryName(countryName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

/*
func (x *X509CreateCertificateOptions) MustSetIntermediateCertificateInGopass(intermediateCertificateInGopass *GopassSecretOptions) {
	err := x.SetIntermediateCertificateInGopass(intermediateCertificateInGopass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}
*/

func (x *X509CreateCertificateOptions) MustSetKeyOutputFilePath(keyOutputFilePath string) {
	err := x.SetKeyOutputFilePath(keyOutputFilePath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetLocality(locality string) {
	err := x.SetLocality(locality)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetOverwriteExistingCertificateInGopass(overwriteExistingCertificateInGopass bool) {
	err := x.SetOverwriteExistingCertificateInGopass(overwriteExistingCertificateInGopass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetUseTemporaryDirectory(useTemporaryDirectory bool) {
	err := x.SetUseTemporaryDirectory(useTemporaryDirectory)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) MustSetVerbose(verbose bool) {
	err := x.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (x *X509CreateCertificateOptions) SetAdditionalSans(additionalSans []string) (err error) {
	if additionalSans == nil {
		return tracederrors.TracedErrorf("additionalSans is nil")
	}

	if len(additionalSans) <= 0 {
		return tracederrors.TracedErrorf("additionalSans has no elements")
	}

	x.AdditionalSans = additionalSans

	return nil
}

func (x *X509CreateCertificateOptions) SetCertificateOutputFilePath(certificateOutputFilePath string) (err error) {
	if certificateOutputFilePath == "" {
		return tracederrors.TracedErrorf("certificateOutputFilePath is empty string")
	}

	x.CertificateOutputFilePath = certificateOutputFilePath

	return nil
}

func (x *X509CreateCertificateOptions) SetCommonName(commonName string) (err error) {
	if commonName == "" {
		return tracederrors.TracedErrorf("commonName is empty string")
	}

	x.CommonName = commonName

	return nil
}

func (x *X509CreateCertificateOptions) SetCountryName(countryName string) (err error) {
	if countryName == "" {
		return tracederrors.TracedErrorf("countryName is empty string")
	}

	x.CountryName = countryName

	return nil
}

/*
func (x *X509CreateCertificateOptions) SetIntermediateCertificateInGopass(intermediateCertificateInGopass *GopassSecretOptions) (err error) {
	if intermediateCertificateInGopass == nil {
		return tracederrors.TracedErrorf("intermediateCertificateInGopass is nil")
	}

	x.IntermediateCertificateInGopass = intermediateCertificateInGopass

	return nil
}*/

func (x *X509CreateCertificateOptions) SetKeyOutputFilePath(keyOutputFilePath string) (err error) {
	if keyOutputFilePath == "" {
		return tracederrors.TracedErrorf("keyOutputFilePath is empty string")
	}

	x.KeyOutputFilePath = keyOutputFilePath

	return nil
}

func (x *X509CreateCertificateOptions) SetLocality(locality string) (err error) {
	if locality == "" {
		return tracederrors.TracedErrorf("locality is empty string")
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
