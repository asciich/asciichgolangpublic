package artifacthandler

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/artifactparameteroptions"
)

// An artifact handler is used to download or update artifacts.
// While artifacts could be some compiled binaries, docker images, vm images...
type ArtifactHandler interface {
	DownloadAndValidateArtifact(ctx context.Context, downloadOptions *artifactparameteroptions.ArtifactDownloadOptions) (downloadedArtifactPath string, err error)
	GetLatestArtifactVersionAsString(artifactName string, verbose bool) (latestVersion string, err error)
	IsHandlingArtifactByName(artifactName string) (isHandlingArtifactByName bool, err error)
	UploadBinaryArtifact(uploadOptions *artifactparameteroptions.UploadArtifactOptions) (err error)
}
