package filesutils

import (
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetAbsolutePath(path string) (string, error) {
	got, err := filepath.Abs(path)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get absolute path: %w", err)
	}

	return got, nil
}
