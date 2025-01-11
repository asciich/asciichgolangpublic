package docker

import (
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/containers"
	"github.com/asciich/asciichgolangpublic/hosts"
)

type Docker interface {
	GetContainerByName(name string) (container containers.Container, err error)
	GetHostDescription() (hostDescription string, err error)
	MustGetContainerByName(name string) (container containers.Container)
	MustGetHostDescription() (hostDescription string)
}

func GetDockerContainerOnHost(host hosts.Host, containerName string) (dockerContainer containers.Container, err error) {
	if host == nil {
		return nil, asciichgolangpublic.TracedErrorNil("host")
	}

	if containerName == "" {
		return nil, asciichgolangpublic.TracedErrorEmptyString("containerName")
	}

	docker, err := GetDockerOnHost(host)
	if err != nil {
		return nil, err
	}

	return docker.GetContainerByName(containerName)
}

func GetDockerOnHost(host hosts.Host) (docker Docker, err error) {
	if host == nil {
		return nil, asciichgolangpublic.TracedErrorNil("host")
	}

	return GetCommandExecutorDockerOnHost(host)
}

func MustGetDockerContainerOnHost(host hosts.Host, containerName string) (dockerContainer containers.Container) {
	dockerContainer, err := GetDockerContainerOnHost(host, containerName)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return dockerContainer
}

func MustGetDockerOnHost(host hosts.Host) (docker Docker) {
	docker, err := GetDockerOnHost(host)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return docker
}
