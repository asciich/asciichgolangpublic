package asciichgolangpublic


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
		return "", TracedErrorf("ArtifactName not set")
	}

	return u.ArtifactName, nil
}

func (u *UploadArtifactOptions) GetAuthOptions() (authOptions []AuthenticationOption, err error) {
	if u.AuthOptions == nil {
		return nil, TracedErrorf("AuthOptions not set")
	}

	if len(u.AuthOptions) <= 0 {
		return nil, TracedErrorf("AuthOptions has no elements")
	}

	return u.AuthOptions, nil
}

func (u *UploadArtifactOptions) GetBinaryPath() (binaryPath string, err error) {
	if u.BinaryPath == "" {
		return "", TracedErrorf("BinaryPath not set")
	}

	return u.BinaryPath, nil
}

func (u *UploadArtifactOptions) GetSignaturePath() (signaturePath string, err error) {
	if u.SignaturePath == "" {
		return "", TracedErrorf("SignaturePath not set")
	}

	return u.SignaturePath, nil
}

func (u *UploadArtifactOptions) GetSoftwareVersionString() (softwareVersionString string, err error) {
	if u.SoftwareVersionString == "" {
		return "", TracedErrorf("SoftwareVersionString not set")
	}

	return u.SoftwareVersionString, nil
}

func (u *UploadArtifactOptions) GetVerbose() (verbose bool) {

	return u.Verbose
}

func (u *UploadArtifactOptions) MustGetArtifactName() (artifactName string) {
	artifactName, err := u.GetArtifactName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return artifactName
}

func (u *UploadArtifactOptions) MustGetAuthOptions() (authOptions []AuthenticationOption) {
	authOptions, err := u.GetAuthOptions()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authOptions
}

func (u *UploadArtifactOptions) MustGetBinaryPath() (binaryPath string) {
	binaryPath, err := u.GetBinaryPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return binaryPath
}

func (u *UploadArtifactOptions) MustGetSignaturePath() (signaturePath string) {
	signaturePath, err := u.GetSignaturePath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return signaturePath
}

func (u *UploadArtifactOptions) MustGetSoftwareVersionString() (softwareVersionString string) {
	softwareVersionString, err := u.GetSoftwareVersionString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return softwareVersionString
}

func (u *UploadArtifactOptions) MustSetArtifactName(artifactName string) {
	err := u.SetArtifactName(artifactName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetAuthOptions(authOptions []AuthenticationOption) {
	err := u.SetAuthOptions(authOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetBinaryPath(binaryPath string) {
	err := u.SetBinaryPath(binaryPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetSignaturePath(signaturePath string) {
	err := u.SetSignaturePath(signaturePath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) MustSetSoftwareVersionString(softwareVersionString string) {
	err := u.SetSoftwareVersionString(softwareVersionString)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UploadArtifactOptions) SetArtifactName(artifactName string) (err error) {
	if artifactName == "" {
		return TracedErrorf("artifactName is empty string")
	}

	u.ArtifactName = artifactName

	return nil
}

func (u *UploadArtifactOptions) SetAuthOptions(authOptions []AuthenticationOption) (err error) {
	if authOptions == nil {
		return TracedErrorf("authOptions is nil")
	}

	if len(authOptions) <= 0 {
		return TracedErrorf("authOptions has no elements")
	}

	u.AuthOptions = authOptions

	return nil
}

func (u *UploadArtifactOptions) SetBinaryPath(binaryPath string) (err error) {
	if binaryPath == "" {
		return TracedErrorf("binaryPath is empty string")
	}

	u.BinaryPath = binaryPath

	return nil
}

func (u *UploadArtifactOptions) SetSignaturePath(signaturePath string) (err error) {
	if signaturePath == "" {
		return TracedErrorf("signaturePath is empty string")
	}

	u.SignaturePath = signaturePath

	return nil
}

func (u *UploadArtifactOptions) SetSoftwareVersionString(softwareVersionString string) (err error) {
	if softwareVersionString == "" {
		return TracedErrorf("softwareVersionString is empty string")
	}

	u.SoftwareVersionString = softwareVersionString

	return nil
}

func (u *UploadArtifactOptions) SetVerbose(verbose bool) {
	u.Verbose = verbose
}
