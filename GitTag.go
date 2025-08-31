package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GitTag interface {
	GetHash(ctx context.Context) (hash string, err error)
	GetName() (name string, err error)
	GetGitRepository() (repo GitRepository, err error)
	IsVersionTag() (isVersionTag bool, err error)
	SetName(name string) (err error)

	// These function can be implemented by embedding the GitTagBase struct:
	GetVersion() (version versionutils.Version, err error)
}
