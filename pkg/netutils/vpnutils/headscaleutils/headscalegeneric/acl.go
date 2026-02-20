package headscalegeneric

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func WriteAclAllOpenAsTemporaryFile(ctx context.Context) (string, error) {
	logging.LogInfoByCtxf(ctx, "Write minimal headscale config as temporary file started.")

	temporaryFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "{}\n")
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Write minimal headscale config as temporary file '%s' finished.", temporaryFilePath)

	return temporaryFilePath, nil
}
