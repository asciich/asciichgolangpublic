package ollamautils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// Start ollama in a docker container in CPU only/ no GPU mode.
func RunCpuOnly(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Run ollama in cpu only mode started.")

	_, err := commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"docker", "run", "-d", "-v", "ollama:/root/.ollama", "-p", "11434:11434", "--name", "ollama", "ollama/ollama"},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ollama in cpu only mode finished.")

	return nil
}
