package dockerutils

import (
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/containers"
	"github.com/asciich/asciichgolangpublic/hosts"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type Docker interface {
	GetContainerByName(name string) (container containers.Container, err error)
	GetHostDescription() (hostDescription string, err error)
}

func GetDockerContainerOnHost(host hosts.Host, containerName string) (dockerContainer containers.Container, err error) {
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

func GetDockerOnHost(host hosts.Host) (docker Docker, err error) {
	if host == nil {
		return nil, tracederrors.TracedErrorNil("host")
	}

	return GetCommandExecutorDockerOnHost(host)
}

func MustGetDockerContainerOnHost(host hosts.Host, containerName string) (dockerContainer containers.Container) {
	dockerContainer, err := GetDockerContainerOnHost(host, containerName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dockerContainer
}

func MustGetDockerOnHost(host hosts.Host) (docker Docker) {
	docker, err := GetDockerOnHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func GetDockerOnLocalHost() (Docker, error) {
	return GetCommandExecutorDocker(commandexecutor.Bash())
}

