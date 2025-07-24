package artifactparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type UploadArtifactOptions struct {
	ArtifactName          string
	BinaryPath            string
	SignaturePath         string
	SoftwareVersionString string
	Verbose               bool
	AuthOptions           []authenticationoptions.AuthenticationOption
}

func NewUploadArtifactOptions() (u *UploadArtifactOptions) {
	return new(UploadArtifactOptions)
}

func NewUploadartifactOptions() (u *UploadArtifactOptions) {
	return new(UploadArtifactOptions)
}

func (u *UploadArtifactOptions) GetArtifactName() (artifactName string, err error) {
	if u.ArtifactName == "" {
		return "", tracederrors.TracedErrorf("ArtifactName not set")
	}

	return u.ArtifactName, nil
}

func (u *UploadArtifactOptions) GetAuthOptions() (authOptions []authenticationoptions.AuthenticationOption, err error) {
	if u.AuthOptions == nil {
		return nil, tracederrors.TracedErrorf("AuthOptions not set")
	}

	if len(u.AuthOptions) <= 0 {
		return nil, tracederrors.TracedErrorf("AuthOptions has no elements")
	}

	return u.AuthOptions, nil
}

func (u *UploadArtifactOptions) GetBinaryPath() (binaryPath string, err error) {
	if u.BinaryPath == "" {
		return "", tracederrors.TracedErrorf("BinaryPath not set")
	}

	return u.BinaryPath, nil
}

func (u *UploadArtifactOptions) GetSignaturePath() (signaturePath string, err error) {
	if u.SignaturePath == "" {
		return "", tracederrors.TracedErrorf("SignaturePath not set")
	}

	return u.SignaturePath, nil
}

func (u *UploadArtifactOptions) GetSoftwareVersionString() (softwareVersionString string, err error) {
	if u.SoftwareVersionString == "" {
		return "", tracederrors.TracedErrorf("SoftwareVersionString not set")
	}

	return u.SoftwareVersionString, nil
}

func (u *UploadArtifactOptions) GetVerbose() (verbose bool) {

	return u.Verbose
}

func (u *UploadArtifactOptions) SetArtifactName(artifactName string) (err error) {
	if artifactName == "" {
		return tracederrors.TracedErrorf("artifactName is empty string")
	}

	u.ArtifactName = artifactName

	return nil
}

func (u *UploadArtifactOptions) SetAuthOptions(authOptions []authenticationoptions.AuthenticationOption) (err error) {
	if authOptions == nil {
		return tracederrors.TracedErrorf("authOptions is nil")
	}

	if len(authOptions) <= 0 {
		return tracederrors.TracedErrorf("authOptions has no elements")
	}

	u.AuthOptions = authOptions

	return nil
}

func (u *UploadArtifactOptions) SetBinaryPath(binaryPath string) (err error) {
	if binaryPath == "" {
		return tracederrors.TracedErrorf("binaryPath is empty string")
	}

	u.BinaryPath = binaryPath

	return nil
}

func (u *UploadArtifactOptions) SetSignaturePath(signaturePath string) (err error) {
	if signaturePath == "" {
		return tracederrors.TracedErrorf("signaturePath is empty string")
	}

	u.SignaturePath = signaturePath

	return nil
}

func (u *UploadArtifactOptions) SetSoftwareVersionString(softwareVersionString string) (err error) {
	if softwareVersionString == "" {
		return tracederrors.TracedErrorf("softwareVersionString is empty string")
	}

	u.SoftwareVersionString = softwareVersionString

	return nil
}

func (u *UploadArtifactOptions) SetVerbose(verbose bool) {
	u.Verbose = verbose
}
