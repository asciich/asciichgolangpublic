package containeroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type AddFileToImageOptions struct {
	// Path to the local file to add:
	SourceFilePath string

	// Path inside the container image where the source file is added:
	PathInImage string

	// New name of the new Image after adding the file:
	NewImageNameAndTag string

	// Overwrite source archive with the updated image archive:
	OverwriteSourceArchive bool
}

func (a *AddFileToImageOptions) GetSourceFilePath() (string, error) {
	if a.SourceFilePath == "" {
		return "", tracederrors.TracedError("SourceFilePath not set")
	}

	return a.SourceFilePath, nil
}

func (a *AddFileToImageOptions) GetPathInImage() (string, error) {
	if a.PathInImage == "" {
		return "", tracederrors.TracedError("PathInImage not set")
	}

	return a.PathInImage, nil
}

func (a *AddFileToImageOptions) GetNewImageNameAndTag() (string, error) {
	if a.NewImageNameAndTag == "" {
		return "", tracederrors.TracedError("NewImageNameAndTag not set")
	}

	return a.NewImageNameAndTag, nil
}
