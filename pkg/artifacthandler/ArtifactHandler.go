package artifacthandler

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/artifactparameteroptions"
)

// An artifact handler is used to download or update artifacts.
// While artifacts could be some compiled binaries, docker images, vm images...
type ArtifactHandler interface {
	DownloadAndValidateArtifact(ctx context.Context, downloadOptions *artifactparameteroptions.ArtifactDownloadOptions) (string, error)
	GetLatestArtifactVersionAsString(ctx context.Context, artifactName string) (string, error)
	IsHandlingArtifactByName(artifactName string) (bool, error)
	UploadBinaryArtifact(ctx context.Context, uploadOptions *artifactparameteroptions.UploadArtifactOptions) error
}
