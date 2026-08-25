package gitutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
)

func GetRepositoryRootPathByPath(ctx context.Context, path string) (string, error) {
	return nativegit.GetRepositoryRootPathByPath(
		ctx,
		path,
	)
}
