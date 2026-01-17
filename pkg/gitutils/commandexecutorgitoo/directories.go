package commandexecutorgitoo

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) GetDirectoryByPath(pathToSubDir ...string) (subDir filesinterfaces.Directory, err error) {
	if len(pathToSubDir) <= 0 {
		return nil, tracederrors.TracedError("pathToSubdir has no elements")
	}

	return g.GetSubDirectory(pathToSubDir...)
}
