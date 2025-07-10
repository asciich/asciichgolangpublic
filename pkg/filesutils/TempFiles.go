package filesutils

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func CreateNamedTemporaryFile(ctx context.Context, fileName string) (string, error) {
	if fileName == "" {
		return "", tracederrors.TracedErrorEmptyString("fileName")
	}

	osFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}
	defer osFile.Close()

	path := osFile.Name()
	logging.LogChangedByCtxf(ctx, "Created temporary file '%s'", path)

	return path, nil
}

func CreateTemporaryFile(ctx context.Context) (string, error) {
	return CreateNamedTemporaryFile(ctx, "tempfile")
}