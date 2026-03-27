package commandexecutordocker

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Image struct {
	docker *CommandExecutorDocker
	name                  string
}

func (i *Image) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	i.name = name

	return nil
}

func (i *Image) GetName() (string, error) {
	if i.name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return i.name, nil
}

func (i *Image) SetDocker(commandExecutorDocker *CommandExecutorDocker) error {
	if commandExecutorDocker == nil {
		return tracederrors.TracedErrorNil("commandExecutorDocker")
	}

	i.docker = commandExecutorDocker

	return nil
}

func (i *Image) GetDocker() (*CommandExecutorDocker, error) {
	if i.docker == nil {
		return nil, tracederrors.TracedError("commandExecutorDocker not set")
	}

	return i.docker, nil
}

func (i *Image) Exists(ctx context.Context) (bool, error) {
	name, err := i.GetName()
	if err != nil {
		return false, err
	}

	docker, err := i.GetDocker()
	if err != nil {
		return false, err
	}

	return docker.ImageExists(ctx, name)
}
