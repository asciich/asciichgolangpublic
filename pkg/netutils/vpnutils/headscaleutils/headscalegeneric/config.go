package headscalegeneric

import (
	"context"
	_ "embed"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

//go:embed files/minimalconfig.yaml
var minimalConfig string

func GetMinimalDockerConfig() string {
	return minimalConfig
}

func WriteMinimalConfigAsTemporaryFile(ctx context.Context) (string, error) {
	logging.LogInfoByCtxf(ctx, "Write minimal headscale config as temporary file started.")

	temporaryFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, GetMinimalDockerConfig())
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Write minimal headscale config as temporary file '%s' finished.", temporaryFilePath)

	return temporaryFilePath, nil
}
