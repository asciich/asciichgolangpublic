package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type UploadArtifactOptions struct {
	ArtifactName          string
	BinaryPath            string
	SignaturePath         string
	SoftwareVersionString string
	Verbose               bool
	AuthOptions           []AuthenticationOption
}

func NewUploadArtifactOptions() (u *UploadArtifactOptions) {
	return new(UploadArtifactOptions)
}

func NewUploadartifactOptions() (u *UploadArtifactOptions) {
	return new(UploadArtifactOptions)
}

func (u *UploadArtifactOptions) GetArtifactName() (artifactName string, err error) {
	if u.ArtifactName == "" {
		return "", errors.TracedErrorf("ArtifactName not set")
	}

	return u.ArtifactName, nil
}

func (u *UploadArtifactOptions) GetAuthOptions() (authOptions []AuthenticationOption, err error) {
	if u.AuthOptions == nil {
		return nil, errors.TracedErrorf("AuthOptions not set")
	}

	if len(u.AuthOptions) <= 0 {
		return nil, errors.TracedErrorf("AuthOptions has no elements")
	}

	return u.AuthOptions, nil
}

func (u *UploadArtifactOptions) GetBinaryPath() (binaryPath string, err error) {
	if u.BinaryPath == "" {
		return "", errors.TracedErrorf("BinaryPath not set")
	}

	return u.BinaryPath, nil
}

func (u *UploadArtifactOptions) GetSignaturePath() (signaturePath string, err error) {
	if u.SignaturePath == "" {
		return "", errors.TracedErrorf("SignaturePath not set")
	}

	return u.SignaturePath, nil
}

func (u *UploadArtifactOptions) GetSoftwareVersionString() (softwareVersionString string, err error) {
	if u.SoftwareVersionString == "" {
		return "", errors.TracedErrorf("SoftwareVersionString not set")
	}

	return u.SoftwareVersionString, nil
}

func (u *UploadArtifactOptions) GetVerbose() (verbose bool) {

	return u.Verbose
}

func (u *UploadArtifactOptions) MustGetArtifactName() (artifactName string) {
	artifactName, err := u.GetArtifactName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return artifactName
}

func (u *UploadArtifactOptions) MustGetAuthOptions() (authOptions []AuthenticationOption) {
	authOptions, err := u.GetAuthOptions()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authOptions
}

func (u *UploadArtifactOptions) MustGetBinaryPath() (binaryPath string) {
	binaryPath, err := u.GetBinaryPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return binaryPath
}

func (u *UploadArtifactOptions) MustGetSignaturePath() (signaturePath string) {
	signaturePath, err := u.GetSignaturePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return signaturePath
}

func (u *UploadArtifactOptions) MustGetSoftwareVersionString() (softwareVersionString string) {
	softwareVersionString, err := u.GetSoftwareVersionString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return softwareVersionString
}

func (u *UploadArtifactOptions) MustSetArtifactName(artifactName string) {
	err := u.SetArtifactName(artifactName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetAuthOptions(authOptions []AuthenticationOption) {
	err := u.SetAuthOptions(authOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetBinaryPath(binaryPath string) {
	err := u.SetBinaryPath(binaryPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetSignaturePath(signaturePath string) {
	err := u.SetSignaturePath(signaturePath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetSoftwareVersionString(softwareVersionString string) {
	err := u.SetSoftwareVersionString(softwareVersionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) SetArtifactName(artifactName string) (err error) {
	if artifactName == "" {
		return errors.TracedErrorf("artifactName is empty string")
	}

	u.ArtifactName = artifactName

	return nil
}

func (u *UploadArtifactOptions) SetAuthOptions(authOptions []AuthenticationOption) (err error) {
	if authOptions == nil {
		return errors.TracedErrorf("authOptions is nil")
	}

	if len(authOptions) <= 0 {
		return errors.TracedErrorf("authOptions has no elements")
	}

	u.AuthOptions = authOptions

	return nil
}

func (u *UploadArtifactOptions) SetBinaryPath(binaryPath string) (err error) {
	if binaryPath == "" {
		return errors.TracedErrorf("binaryPath is empty string")
	}

	u.BinaryPath = binaryPath

	return nil
}

func (u *UploadArtifactOptions) SetSignaturePath(signaturePath string) (err error) {
	if signaturePath == "" {
		return errors.TracedErrorf("signaturePath is empty string")
	}

	u.SignaturePath = signaturePath

	return nil
}

func (u *UploadArtifactOptions) SetSoftwareVersionString(softwareVersionString string) (err error) {
	if softwareVersionString == "" {
		return errors.TracedErrorf("softwareVersionString is empty string")
	}

	u.SoftwareVersionString = softwareVersionString

	return nil
}

func (u *UploadArtifactOptions) SetVerbose(verbose bool) {
	u.Verbose = verbose
}
