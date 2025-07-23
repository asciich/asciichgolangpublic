package dockerutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/hosts"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/containerutils/containerinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/dockerutils/dockerinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
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

	return GetCommandExecutorDockerOnHost(host)
}

func MustGetDockerContainerOnHost(host hosts.Host, containerName string) (dockerContainer containerinterfaces.Container) {
	dockerContainer, err := GetDockerContainerOnHost(host, containerName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dockerContainer
}

func MustGetDockerOnHost(host hosts.Host) (docker dockerinterfaces.Docker) {
	docker, err := GetDockerOnHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func GetDockerOnLocalHost() (dockerinterfaces.Docker, error) {
	return GetCommandExecutorDocker(commandexecutor.Bash())
}

func ListContainerNames(ctx context.Context) ([]string, error) {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return nil, err
	}

	return docker.ListContainerNames(ctx)
}
