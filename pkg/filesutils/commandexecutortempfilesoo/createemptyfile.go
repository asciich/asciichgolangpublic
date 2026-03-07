package commandexecutortempfilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
)

func CreateEmptyTemporaryFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (filesinterfaces.File, error) {
	path, err := commandexecutortempfile.CreateEmptyTemporaryFile(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	return commandexecutorfileoo.New(commandExecutor, path)
}

func CreateLocalEmptyTemporaryFile(ctx context.Context) (filesinterfaces.File, error) {
	return CreateEmptyTemporaryFile(ctx, commandexecutorexecoo.Exec())
}
