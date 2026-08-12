package copilotutils

import (
	"fmt"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunCopilotCliContainer(options *RunCopilotCliContainerOptions) (ret string, err error) {
	if options == nil {
		return "", tracederrors.TracedErrorf("options is nil")
	}

	prompt, err := options.GetPrompt()
	if err != nil {
		return "", err
	}

	containerImage, err := options.GetContainerImage()
	if err != nil {
		return "", err
	}

	copilotGithubToken, err := getCopilotGithubToken()
	if err != nil {
		return "", err
	}

	ctx := contextutils.ContextVerbose()

	docker, err := dockerutils.GetDockerOnLocalHost()
	if err != nil {
		return "", err
	}

	runOptions := &dockeroptions.DockerRunContainerOptions{
		ImageName: containerImage,
		Command:   []string{"copilot", prompt},
		AdditionalEnvVars: map[string]string{
			"COPILOT_GITHUB_TOKEN": copilotGithubToken,
		},
	}

	if options.SystemPrompt != "" {
		runOptions.AdditionalEnvVars["SYSTEM_PROMPT"] = options.SystemPrompt
	}

	if options.GitUserName != "" {
		runOptions.AdditionalEnvVars["GIT_USER_NAME"] = options.GitUserName
	}

	if options.GitUserEmail != "" {
		runOptions.AdditionalEnvVars["GIT_USER_EMAIL"] = options.GitUserEmail
	}

	if options.WorkspacePath != "" {
		runOptions.Mounts = append(
			runOptions.Mounts,
			fmt.Sprintf("%s:/workspace", options.WorkspacePath),
		)
	}

	if options.GitConfigPath != "" {
		runOptions.Mounts = append(
			runOptions.Mounts,
			fmt.Sprintf("%s:/root/.gitconfig:ro", options.GitConfigPath),
		)
	}

	if options.SshConfigPath != "" {
		runOptions.Mounts = append(
			runOptions.Mounts,
			fmt.Sprintf("%s:/root/.ssh/config:ro", options.SshConfigPath),
		)
	}

	if options.GetSSHAgent() {
		sshAuthSock, err := getSSHAuthSock()
		if err != nil {
			return "", err
		}

		runOptions.AdditionalEnvVars["SSH_AUTH_SOCK"] = sshAuthSock
		runOptions.Mounts = append(
			runOptions.Mounts,
			fmt.Sprintf("%s:%s", sshAuthSock, sshAuthSock),
		)
	}

	container, err := docker.RunContainer(ctx, runOptions)
	if err != nil {
		return "", err
	}
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	ret, err = container.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"copilot", prompt},
		},
	)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func getCopilotGithubToken() (copilotGithubToken string, err error) {
	copilotGithubToken = os.Getenv("COPILOT_GITHUB_TOKEN")
	if copilotGithubToken == "" {
		return "", tracederrors.TracedErrorf("COPILOT_GITHUB_TOKEN is not set")
	}

	return copilotGithubToken, nil
}

func getSSHAuthSock() (sshAuthSock string, err error) {
	sshAuthSock = os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		return "", tracederrors.TracedErrorf("SSH_AUTH_SOCK is not set")
	}

	return sshAuthSock, nil
}
