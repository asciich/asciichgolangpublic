package nativefiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func WriteString(ctx context.Context, pathToWrite string, content string) error {
	if pathToWrite == "" {
		return tracederrors.TracedErrorEmptyString("pathToWrite")
	}

	err := os.WriteFile(pathToWrite, []byte(content), 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to write to file '%s': %w", pathToWrite, err)
	}

	logging.LogChangedByCtxf(ctx, "Wrote content to file '%s'.", pathToWrite)

	return nil
}

func WriteBytes(ctx context.Context, pathToWrite string, content []byte) error {
	if pathToWrite == "" {
		return tracederrors.TracedErrorEmptyString("pathToWrite")
	}

	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	err := os.WriteFile(pathToWrite, content, 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to write to file '%s': %w", pathToWrite, err)
	}

	logging.LogChangedByCtxf(ctx, "Wrote content to file '%s'.", pathToWrite)

	return nil
}

func ReadAsString(ctx context.Context, pathToRead string) (string, error) {
	if pathToRead == "" {
		return "", tracederrors.TracedErrorEmptyString("pathToRead")
	}

	content, err := os.ReadFile(pathToRead)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to read file '%s': %w", pathToRead, err)
	}

	logging.LogInfoByCtxf(ctx, "Read content of file '%s'.", pathToRead)

	return string(content), nil
}

func ReadAsBytes(ctx context.Context, pathToRead string) ([]byte, error) {
	if pathToRead == "" {
		return nil, tracederrors.TracedErrorEmptyString("pathToRead")
	}

	content, err := os.ReadFile(pathToRead)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to read file '%s': %w", pathToRead, err)
	}

	logging.LogInfoByCtxf(ctx, "Read content of file '%s'.", pathToRead)

	return content, nil
}
