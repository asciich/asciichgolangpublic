package copilotutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const defaultContainerImage = "ghcr.io/henrybravo/docker-sandbox-run-copilot"

type RunCopilotCliContainerOptions struct {
	WorkspacePath  string
	SSHAgent       bool
	SystemPrompt   string
	Prompt         string
	ContainerImage string
	GitUserName    string
	GitUserEmail   string
	GitConfigPath  string
	SshConfigPath  string
}

func NewRunCopilotCliContainerOptions() (r *RunCopilotCliContainerOptions) {
	return new(RunCopilotCliContainerOptions)
}

func (r *RunCopilotCliContainerOptions) GetWorkspacePath() (workspacePath string, err error) {
	if r.WorkspacePath == "" {
		return "", tracederrors.TracedErrorf("WorkspacePath not set")
	}

	return r.WorkspacePath, nil
}

func (r *RunCopilotCliContainerOptions) SetWorkspacePath(workspacePath string) (err error) {
	if workspacePath == "" {
		return tracederrors.TracedErrorf("workspacePath is empty string")
	}

	r.WorkspacePath = workspacePath

	return nil
}

func (r *RunCopilotCliContainerOptions) GetSSHAgent() (sshAgent bool) {
	return r.SSHAgent
}

func (r *RunCopilotCliContainerOptions) SetSSHAgent(sshAgent bool) {
	r.SSHAgent = sshAgent
}

func (r *RunCopilotCliContainerOptions) GetSystemPrompt() (systemPrompt string, err error) {
	if r.SystemPrompt == "" {
		return "", tracederrors.TracedErrorf("SystemPrompt not set")
	}

	return r.SystemPrompt, nil
}

func (r *RunCopilotCliContainerOptions) SetSystemPrompt(systemPrompt string) (err error) {
	if systemPrompt == "" {
		return tracederrors.TracedErrorf("systemPrompt is empty string")
	}

	r.SystemPrompt = systemPrompt

	return nil
}

func (r *RunCopilotCliContainerOptions) GetPrompt() (prompt string, err error) {
	if r.Prompt == "" {
		return "", tracederrors.TracedErrorf("Prompt not set")
	}

	return r.Prompt, nil
}

func (r *RunCopilotCliContainerOptions) SetPrompt(prompt string) (err error) {
	if prompt == "" {
		return tracederrors.TracedErrorf("prompt is empty string")
	}

	r.Prompt = prompt

	return nil
}

func (r *RunCopilotCliContainerOptions) GetContainerImage() (containerImage string, err error) {
	if r.ContainerImage == "" {
		return defaultContainerImage, nil
	}

	return r.ContainerImage, nil
}

func (r *RunCopilotCliContainerOptions) SetContainerImage(containerImage string) (err error) {
	if containerImage == "" {
		return tracederrors.TracedErrorf("containerImage is empty string")
	}

	r.ContainerImage = containerImage

	return nil
}

func (r *RunCopilotCliContainerOptions) GetGitUserName() (gitUserName string, err error) {
	if r.GitUserName == "" {
		return "", tracederrors.TracedErrorf("GitUserName not set")
	}

	return r.GitUserName, nil
}

func (r *RunCopilotCliContainerOptions) SetGitUserName(gitUserName string) (err error) {
	if gitUserName == "" {
		return tracederrors.TracedErrorf("gitUserName is empty string")
	}

	r.GitUserName = gitUserName

	return nil
}

func (r *RunCopilotCliContainerOptions) GetGitUserEmail() (gitUserEmail string, err error) {
	if r.GitUserEmail == "" {
		return "", tracederrors.TracedErrorf("GitUserEmail not set")
	}

	return r.GitUserEmail, nil
}

func (r *RunCopilotCliContainerOptions) SetGitUserEmail(gitUserEmail string) (err error) {
	if gitUserEmail == "" {
		return tracederrors.TracedErrorf("gitUserEmail is empty string")
	}

	r.GitUserEmail = gitUserEmail

	return nil
}

func (r *RunCopilotCliContainerOptions) GetGitConfigPath() (gitConfigPath string, err error) {
	if r.GitConfigPath == "" {
		return "", tracederrors.TracedErrorf("GitConfigPath not set")
	}

	return r.GitConfigPath, nil
}

func (r *RunCopilotCliContainerOptions) SetGitConfigPath(gitConfigPath string) (err error) {
	if gitConfigPath == "" {
		return tracederrors.TracedErrorf("gitConfigPath is empty string")
	}

	r.GitConfigPath = gitConfigPath

	return nil
}

func (r *RunCopilotCliContainerOptions) GetSshConfigPath() (sshConfigPath string, err error) {
	if r.SshConfigPath == "" {
		return "", tracederrors.TracedErrorf("SshConfigPath not set")
	}

	return r.SshConfigPath, nil
}

func (r *RunCopilotCliContainerOptions) SetSshConfigPath(sshConfigPath string) (err error) {
	if sshConfigPath == "" {
		return tracederrors.TracedErrorf("sshConfigPath is empty string")
	}

	r.SshConfigPath = sshConfigPath

	return nil
}
