package asciichgolangpublic

import "github.com/asciich/asciichgolangpublic/files"

// An artifact handler is used to download or update artifacts.
// While artifacts could be some compiled binaries, docker images, vm images...
type ArtifactHandler interface {
	DownloadAndValidateArtifact(downloadOptions *ArtifactDownloadOptions) (downloadedArtifact files.File, err error)
	MustDownloadAndValidateArtifact(downloadOptions *ArtifactDownloadOptions) (downloadedArtifact files.File)
	GetLatestArtifactVersionAsString(artifactName string, verbose bool) (latestVersion string, err error)
	IsHandlingArtifactByName(artifactName string) (isHandlingArtifactByName bool, err error)
	UploadBinaryArtifact(uploadOptions *UploadArtifactOptions) (err error)
}
