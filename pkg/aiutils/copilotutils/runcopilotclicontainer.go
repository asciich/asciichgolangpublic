package copilotutils

import (
	"fmt"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const defaultAgentHome = "/home/agent"

func RunCopilotCliContainer(options *RunCopilotCliContainerOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorf("options is nil")
	}

	prompt, err := options.GetPrompt()
	if err != nil {
		return nil, err
	}

	containerImage, err := options.GetContainerImage()
	if err != nil {
		return nil, err
	}

	copilotGithubToken, err := getCopilotGithubToken()
	if err != nil {
		return nil, err
	}

	containerName, err := generateContainerName()
	if err != nil {
		return nil, err
	}

	ctx := contextutils.ContextVerbose()

	logging.LogInfoByCtxf(ctx, "Run copilot CLI container '%s' using image '%s' started.", containerName, containerImage)

	docker, err := dockerutils.GetDockerOnLocalHost()
	if err != nil {
		return nil, err
	}

	runOptions := &dockeroptions.DockerRunContainerOptions{
		ImageName: containerImage,
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
			fmt.Sprintf("%s:%s/.gitconfig:ro", options.GitConfigPath, defaultAgentHome),
		)
	}

	if options.SshConfigPath != "" {
		runOptions.Mounts = append(
			runOptions.Mounts,
			fmt.Sprintf("%s:%s/.ssh/config:ro", options.SshConfigPath, defaultAgentHome),
		)
	}

	if options.GetSSHAgent() {
		sshAuthSock, err := getSSHAuthSock()
		if err != nil {
			return nil, err
		}

		runOptions.AdditionalEnvVars["SSH_AUTH_SOCK"] = sshAuthSock
		runOptions.Mounts = append(
			runOptions.Mounts,
			fmt.Sprintf("%s:%s", sshAuthSock, sshAuthSock),
		)
	}

	err = runOptions.SetName(containerName)
	if err != nil {
		return nil, err
	}

	runOptions.Command = []string{"copilot", "--prompt", prompt}
	runOptions.KeepStoppedContainer = true
	runOptions.EntryPoint = []string{}

	ret, err := docker.RunCommandInTemporaryContainer(ctx, runOptions)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run copilot CLI container '%s' using image '%s' finished.", containerName, containerImage)

	return ret, nil
}

func generateContainerName() (containerName string, err error) {
	randomSuffix, err := randomgenerator.GetRandomString(8)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("copilot-cli-%s", randomSuffix), nil
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
