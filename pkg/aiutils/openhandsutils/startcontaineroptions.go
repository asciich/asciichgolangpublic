package openhandsutils

import (
	"strconv"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type StartContainerOptions struct {
	ContainerName string

	Port int

	// By default it binds only to 127.0.0.1 to be reachable from the local machine.
	// If set to true other machines in the network can reach openhands as well.
	ReachableByOtherMachines bool

	WorkspacePath string
}

func (o *StartContainerOptions) GetContainerName() (string, error) {
	if o.ContainerName == "" {
		return "", tracederrors.TracedError("ContainerName not set")
	}

	return o.ContainerName, nil
}

func (o *StartContainerOptions) GetPort() (int, error) {
	if o.Port <= 0 {
		return 0, tracederrors.TracedErrorf("Port '%d' is not set or is invalid.", o.Port)
	}

	return o.Port, nil
}

func (o *StartContainerOptions) GetReachableByOtherMachines() bool {
	return o.ReachableByOtherMachines
}

func (o *StartContainerOptions) GetBindAddress() (string, error) {
	port, err := o.GetPort()
	if err != nil {
		return "", err
	}

	prefix := "127.0.0.1"
	if o.ReachableByOtherMachines {
		prefix = "0.0.0.0"
	}

	return prefix + ":" + strconv.Itoa(port), nil
}

func (o *StartContainerOptions) GetWorkspacePath() (string, error) {
	if o.WorkspacePath == "" {
		return "", tracederrors.TracedError("WorkspacePath not set")
	}

	got, err := filesutils.GetAbsolutePath(o.WorkspacePath)
	if err != nil {
		return "", err
	}

	return got, nil
}
