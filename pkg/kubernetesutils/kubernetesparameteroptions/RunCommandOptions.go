package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KubernetesRunCommandOptions struct {
	RunCommandOptions               *parameteroptions.RunCommandOptions
	Image                           string
	PodName                         string
	ReplicaSetName                  string
	DeploymentName                  string
	ContainerName                   string
	DeleteAlreadyExistingPod        bool
	DeleteAlreadyExistingReplicaSet bool
	DeleteAlreadyExistingDeployment bool

	// Wait until pod is in "running" state
	WaitForPodRunning bool

	// Wait until ReplicaSet has all replicas available
	WaitForReplicaSetAvailable bool

	// Wait until Deployment has all replicas available
	WaitForDeploymentAvailable bool

	// Number of replicas for ReplicaSet/Deployment (default: 1)
	Replicas int32

	StdinBytes []byte

	// SecretEnvVars maps environment variable names to secret names and keys
	// Format: map[envVarName]SecretEnvVarSource
	SecretEnvVars map[string]SecretEnvVarSource

	// SecretMounts defines secrets to mount as files in the container
	// Format: map[mountPath]SecretMountSource
	SecretMounts map[string]SecretMountSource
}

// SecretEnvVarSource defines the source of a secret for an environment variable
type SecretEnvVarSource struct {
	SecretName string
	SecretKey  string
}

// SecretMountSource defines the source of a secret to mount as files
type SecretMountSource struct {
	SecretName string
}

func (r *KubernetesRunCommandOptions) GetContainerName() (string, error) {
	if r.ContainerName == "" {
		// If the container name is not explicitly defined, use the pod, ReplicaSet, or Deployment name:
		if r.PodName != "" {
			return r.GetPodName()
		}
		if r.ReplicaSetName != "" {
			return r.GetReplicaSetName()
		}
		if r.DeploymentName != "" {
			return r.GetDeploymentName()
		}
		return "", tracederrors.TracedError("ContainerName not set and no PodName, ReplicaSetName, or DeploymentName available")
	}

	return r.ContainerName, nil
}

func (r *KubernetesRunCommandOptions) GetPodName() (string, error) {
	if r.PodName == "" {
		return "", tracederrors.TracedError("PodName not set")
	}

	return r.PodName, nil
}

func (r *KubernetesRunCommandOptions) GetImageName() (string, error) {
	if r.Image == "" {
		return "", tracederrors.TracedError("ImageName not set")
	}

	return r.Image, nil
}

func (r *KubernetesRunCommandOptions) GetRunCommandOptions() (*parameteroptions.RunCommandOptions, error) {
	if r.RunCommandOptions == nil {
		return nil, tracederrors.TracedError("RunCommandOptions not set")
	}

	return r.RunCommandOptions, nil
}

func (r *KubernetesRunCommandOptions) GetCommand() ([]string, error) {
	runCommandOptions, err := r.GetRunCommandOptions()
	if err != nil {
		return nil, err
	}

	got, err := runCommandOptions.GetCommand()
	if err != nil {
		return nil, err
	}

	return slicesutils.GetDeepCopyOfStringsSlice(got), nil
}

func (r *KubernetesRunCommandOptions) IsStinDataAvailable() bool {
	return len(r.StdinBytes) > 0
}

func (r *KubernetesRunCommandOptions) GetReplicaSetName() (string, error) {
	if r.ReplicaSetName == "" {
		return "", tracederrors.TracedError("ReplicaSetName not set")
	}

	return r.ReplicaSetName, nil
}

func (r *KubernetesRunCommandOptions) GetReplicas() int32 {
	if r.Replicas <= 0 {
		return 1
	}
	return r.Replicas
}

func (r *KubernetesRunCommandOptions) GetDeploymentName() (string, error) {
	if r.DeploymentName == "" {
		return "", tracederrors.TracedError("DeploymentName not set")
	}

	return r.DeploymentName, nil
}
