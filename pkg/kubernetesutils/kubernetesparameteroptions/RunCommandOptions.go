package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type RunCommandOptions struct {
	RunCommandOptions                  parameteroptions.RunCommandOptions
	Namespace                          string
	Image                              string
	PodName                            string
	ReplicaSetName                     string
	ContainerName                      string
	Command                            []string
	DeleteAlreadyExistingPod           bool
	DeleteAlreadyExistingReplicaSet    bool

	// Wait until pod is in "running" state
	WaitForPodRunning bool

	// Wait until ReplicaSet has all replicas available
	WaitForReplicaSetAvailable bool

	// Number of replicas for ReplicaSet (default: 1)
	Replicas int32

	StdinBytes []byte
}

func (r *RunCommandOptions) GetNamespaceName() (string, error) {
	if r.Namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	return r.Namespace, nil
}

func (r *RunCommandOptions) GetContainerName() (string, error) {
	if r.ContainerName == "" {
		// If the container name is not explicitly defined, use the pod or ReplicaSet name:
		if r.PodName != "" {
			return r.GetPodName()
		}
		if r.ReplicaSetName != "" {
			return r.GetReplicaSetName()
		}
		return "", tracederrors.TracedError("ContainerName not set and no PodName or ReplicaSetName available")
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

func (r *RunCommandOptions) IsStinDataAvailable() bool {
	return len(r.StdinBytes) > 0 
}

func (r *RunCommandOptions) GetReplicaSetName() (string, error) {
	if r.ReplicaSetName == "" {
		return "", tracederrors.TracedError("ReplicaSetName not set")
	}

	return r.ReplicaSetName, nil
}

func (r *RunCommandOptions) GetReplicas() int32 {
	if r.Replicas <= 0 {
		return 1
	}
	return r.Replicas
}