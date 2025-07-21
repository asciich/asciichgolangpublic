package artifactparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type ArtifactDownloadOptions struct {
	ArtifactName      string
	OutputPath        string
	VersionToDownload string
	OverwriteExisting bool
	Verbose           bool
}

func NewArtifactDownloadOptions() (a *ArtifactDownloadOptions) {
	return new(ArtifactDownloadOptions)
}

func NewAsciichArtifactDownloadOptions() (a *ArtifactDownloadOptions) {
	return new(ArtifactDownloadOptions)
}

func (a *ArtifactDownloadOptions) GetArtifactName() (artifactName string, err error) {
	if a.ArtifactName == "" {
		return "", tracederrors.TracedErrorf("ArtifactName not set")
	}

	return a.ArtifactName, nil
}

func (a *ArtifactDownloadOptions) GetOutputPath() (outputPath string, err error) {
	if a.OutputPath == "" {
		return "", tracederrors.TracedErrorf("OutputPath not set")
	}

	return a.OutputPath, nil
}

func (a *ArtifactDownloadOptions) GetOverwriteExisting() (overwriteExisting bool, err error) {

	return a.OverwriteExisting, nil
}

func (a *ArtifactDownloadOptions) GetVerbose() (verbose bool, err error) {

	return a.Verbose, nil
}

func (a *ArtifactDownloadOptions) GetVersionToDownload() (versionToDownload string, err error) {
	if a.VersionToDownload == "" {
		return "", tracederrors.TracedErrorf("VersionToDownload not set")
	}

	return a.VersionToDownload, nil
}

func (a *ArtifactDownloadOptions) IsOutputPathSet() (isSet bool) {
	return a.OutputPath != ""
}

func (a *ArtifactDownloadOptions) IsVersionToDownloadSet() (isSet bool) {
	return a.VersionToDownload != ""
}

func (a *ArtifactDownloadOptions) SetArtifactName(artifactName string) (err error) {
	if artifactName == "" {
		return tracederrors.TracedErrorf("artifactName is empty string")
	}

	a.ArtifactName = artifactName

	return nil
}

func (a *ArtifactDownloadOptions) SetOutputPath(outputPath string) (err error) {
	if outputPath == "" {
		return tracederrors.TracedErrorf("outputPath is empty string")
	}

	a.OutputPath = outputPath

	return nil
}

func (a *ArtifactDownloadOptions) SetOverwriteExisting(overwriteExisting bool) (err error) {
	a.OverwriteExisting = overwriteExisting

	return nil
}

func (a *ArtifactDownloadOptions) SetVerbose(verbose bool) (err error) {
	a.Verbose = verbose

	return nil
}

func (a *ArtifactDownloadOptions) SetVersionToDownload(versionToDownload string) (err error) {
	if versionToDownload == "" {
		return tracederrors.TracedErrorf("versionToDownload is empty string")
	}

	a.VersionToDownload = versionToDownload

	return nil
}
