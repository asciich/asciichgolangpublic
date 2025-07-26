package tempfiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
