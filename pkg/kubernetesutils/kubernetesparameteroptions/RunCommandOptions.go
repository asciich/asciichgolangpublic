package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type RunCommandOptions struct {
	RunCommandOptions        parameteroptions.RunCommandOptions
	Namespace                string
	Image                    string
	PodName                  string
	ContainerName            string
	Command                  []string
	DeleteAlreadyExistingPod bool

	// Wait until pod is in "running" state
	WaitForPodRunning bool
}

func (r *RunCommandOptions) GetNamespaceName() (string, error) {
	if r.Namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	return r.Namespace, nil
}

func (r *RunCommandOptions) GetContainerName() (string, error) {
	if r.ContainerName == "" {
		// If the container name is not explicitly defined the same name as for the pod is used:
		return r.GetPodName()
	}

	return r.ContainerName, nil
}

func (r *RunCommandOptions) GetPodName() (string, error) {
	if r.PodName == "" {
		return "", tracederrors.TracedError("PodName not set")
	}

	return r.PodName, nil
}

func (r *RunCommandOptions) GetImageName() (string, error) {
	if r.Image == "" {
		return "", tracederrors.TracedError("ImageName not set")
	}

	return r.Image, nil
}

func (r *RunCommandOptions) GetCommand() ([]string, error) {
	if len(r.Command) <= 0 {
		return nil, tracederrors.TracedError("Command not set")
	}

	return slicesutils.GetDeepCopyOfStringsSlice(r.Command), nil
}
