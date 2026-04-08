package containeroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type AddFileToImageArchiveOptions struct {
	// Path to the local file to add:
	SourceFilePath string

	// Path inside the container image where the source file is added:
	PathInImage string

	// New name of the new Image after adding the file:
	NewImageNameAndTag string

	// Overwrite source archive with the updated image archive:
	OverwriteSourceArchive bool

	// FileMode to set.
	// Eg. to set it to 0644 you can use: pointerutils.ToInt64Pointer(0644) during initialization.
	Mode *int64
}

func (a *AddFileToImageArchiveOptions) GetSourceFilePath() (string, error) {
	if a.SourceFilePath == "" {
		return "", tracederrors.TracedError("SourceFilePath not set")
	}

	return a.SourceFilePath, nil
}

func (a *AddFileToImageArchiveOptions) GetPathInImage() (string, error) {
	if a.PathInImage == "" {
		return "", tracederrors.TracedError("PathInImage not set")
	}

	return a.PathInImage, nil
}

func (a *AddFileToImageArchiveOptions) GetNewImageNameAndTag() (string, error) {
	if a.NewImageNameAndTag == "" {
		return "", tracederrors.TracedError("NewImageNameAndTag not set")
	}

	return a.NewImageNameAndTag, nil
}

func (a *AddFileToImageArchiveOptions) GetMode() (int64, error) {
	if a.Mode == nil {
		return 0, tracederrors.TracedErrorf("Mode not set")
	}

	return *a.Mode, nil
}
