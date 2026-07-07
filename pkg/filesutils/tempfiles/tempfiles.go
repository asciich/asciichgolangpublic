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

func CreateTemporaryFileFromContentString(ctx context.Context, content string) (string, error) {
	return CreateTemporaryFileFromContentBytes(ctx, []byte(content))
}

func CreateTemporaryFileFromContentBytes(ctx context.Context, content []byte) (string, error) {
	if content == nil {
		return "", tracederrors.TracedErrorEmptyString("content")
	}

	tempFile, err := CreateTemporaryFile(ctx)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(tempFile, content, 0644)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to write content in temporary file '%s': %w", tempFile, err)
	}

	return tempFile, err
}
