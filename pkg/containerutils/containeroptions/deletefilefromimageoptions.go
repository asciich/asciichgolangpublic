package containeroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type DeleteFileFromImageOptions struct {
	// The path of the file to delete:
	PathInImage string

	// New name of the new Image after adding the file:
	NewImageNameAndTag string

	// Overwrite source archive with the updated image archive:
	OverwriteSourceArchive bool
}

func (d *DeleteFileFromImageOptions) GetPathInImage() (string, error) {
	if d.PathInImage == "" {
		return "", tracederrors.TracedError("PathInImage not set")
	}

	return d.PathInImage, nil
}

func (d *DeleteFileFromImageOptions) GetNewImageNameAndTag() (string, error) {
	if d.NewImageNameAndTag == "" {
		return "", tracederrors.TracedError("NewImageNameAndTag not set")
	}

	return d.NewImageNameAndTag, nil
}
