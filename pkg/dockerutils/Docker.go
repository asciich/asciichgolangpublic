package dockerutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/hosts"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetDockerContainerOnHost(host hosts.Host, containerName string) (dockerContainer containerinterfaces.Container, err error) {
	if host == nil {
		return nil, tracederrors.TracedErrorNil("host")
	}

	if containerName == "" {
		return nil, tracederrors.TracedErrorEmptyString("containerName")
	}

	docker, err := GetDockerOnHost(host)
	if err != nil {
		return nil, err
	}

	return docker.GetContainerByName(containerName)
}

func GetDockerOnHost(host hosts.Host) (docker dockerinterfaces.Docker, err error) {
	if host == nil {
		return nil, tracederrors.TracedErrorNil("host")
	}

	return commandexecutordocker.GetCommandExecutorDockerOnHost(host)
}

func GetDockerOnLocalHost() (dockerinterfaces.Docker, error) {
	return nativedocker.NewDocker(), nil
}

func ListContainerNames(ctx context.Context) ([]string, error) {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return nil, err
	}

	return docker.ListContainerNames(ctx)
}
