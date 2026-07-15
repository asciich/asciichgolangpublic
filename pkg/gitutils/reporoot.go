package gitutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgit"
)

func GetRepositoryRootPathByPath(ctx context.Context, path string) (string, error) {
	return commandexecutorgit.GetRepositoryRootPathByPath(
		ctx,
		commandexecutorexecoo.Exec(),
		path,
	)
}
