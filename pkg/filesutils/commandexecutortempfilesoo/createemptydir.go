package commandexecutortempfilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
)

func CreateEmptyTemporaryDirectory(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (filesinterfaces.Directory, error) {
	path, err := commandexecutortempfile.CreateEmptyTemporaryDirectory(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	return commandexecutorfileoo.NewDirectory(commandExecutor, path)
}

func CreateLocalEmptyTemporaryDirectory(ctx context.Context) (filesinterfaces.Directory, error) {
	return CreateEmptyTemporaryDirectory(ctx, commandexecutorexecoo.Exec())
}
