package kubernetesparameteroptions

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/slicesutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/parameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type RunCommandOptions struct {
	RunCommandOptions        parameteroptions.RunCommandOptions
	Namespace                string
	Image                    string
	PodName                  string
	Command                  []string
	DeleteAlreadyExistingPod bool
}

func (r *RunCommandOptions) GetNamespaceName() (string, error) {
	if r.Namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	return r.Namespace, nil
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
