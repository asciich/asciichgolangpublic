package containeroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateSingleFileArchiveOptions struct {
	// Path to the local file to add:
	SourceFilePath string

	// Path inside the container image where the source file is added:
	PathInImage string

	// New name of the new Image after adding the file:
	NewImageNameAndTag string

	// FileMode to set.
	// Eg. to set it to 0644 you can use: pointerutils.ToInt64Pointer(0644) during initialization.
	Mode *int64

	// Architexture, usually "amd64" or "arm" (32bit) or "arm64"
	Architecture string
}

func (a *CreateSingleFileArchiveOptions) GetSourceFilePath() (string, error) {
	if a.SourceFilePath == "" {
		return "", tracederrors.TracedError("SourceFilePath not set")
	}

	return a.SourceFilePath, nil
}

func (a *CreateSingleFileArchiveOptions) GetPathInImage() (string, error) {
	if a.PathInImage == "" {
		return "", tracederrors.TracedError("PathInImage not set")
	}

	return a.PathInImage, nil
}

func (a *CreateSingleFileArchiveOptions) GetNewImageNameAndTag() (string, error) {
	if a.NewImageNameAndTag == "" {
		return "", tracederrors.TracedError("NewImageNameAndTag not set")
	}

	return a.NewImageNameAndTag, nil
}

func (a *CreateSingleFileArchiveOptions) GetMode() (int64, error) {
	if a.Mode == nil {
		return 0, tracederrors.TracedErrorf("Mode not set")
	}

	return *a.Mode, nil
}

func (a *CreateSingleFileArchiveOptions) GetArchitecture() (string, error) {
	if a.Architecture == "" {
		return "", tracederrors.TracedError("Architecture not set")
	}

	return a.Architecture, nil
}
