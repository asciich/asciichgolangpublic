package aiderutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const DEFAULT_CONTAINER_NAME = "aider-local"
const DEFAULT_CONTAINER_TAG = "latest"
const DEFAULT_CONTAINER_NAME_AND_TAG = DEFAULT_CONTAINER_NAME + ":" + DEFAULT_CONTAINER_TAG

func BuildAiderDockerContainer(ctx context.Context) error {
	const tag = "latest"

	logging.LogInfoByCtxf(ctx, "Build local aider docker container '%s' started.", DEFAULT_CONTAINER_NAME_AND_TAG)

	tempDir, err := tempfiles.CreateTempDir(ctx)
	if err != nil {
		return err
	}
	defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

	content := `
FROM python:3.12-slim

RUN apt-get update && apt-get install -y git && \
    pip install --no-cache-dir --upgrade pip setuptools wheel && \
    pip install --no-cache-dir aider-chat

env HOME /aider-home
RUN mkdir -p $HOME && chmod 777 $HOME

WORKDIR /app

# Initialize Git to avoid warnings about an uninitialized repository
RUN git config --global user.name "aider" && \
    git config --global user.email "aider@asciich.ch"
`
	err = nativefiles.WriteString(ctx, filepath.Join(tempDir, "Dockerfile"), content)
	if err != nil {
		return err
	}

	_, err = commandexecutorexec.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"bash", "-c", fmt.Sprintf("cd '%s' && docker build . -t '%s' 2>&1", tempDir, DEFAULT_CONTAINER_NAME_AND_TAG)},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Build local aider docker container '%s' finished.", DEFAULT_CONTAINER_NAME_AND_TAG)

	return nil
}

func GetRunCommand(noninteractive bool) ([]string, error) {
	cmd := []string{
		"docker",
		"run",
		"--rm",
	}

	if noninteractive {
		cmd = append(cmd, "-t")
	} else {
		cmd = append(cmd, "-it")
	}

	uid := os.Getuid()
	gid := os.Getgid()

	workingDirectory, err := osutils.GetCurrentWorkingDirectoryAsString()
	if err != nil {
		return nil, err
	}

	extend := []string{
		"-e",
		"OLLAMA_API_BASE=http://host.docker.internal:11434",
		"-v",
		workingDirectory + ":/app",
		"-w",
		"/app",
		"--add-host=host.docker.internal:host-gateway",
		"--user",
		fmt.Sprintf("%d:%d", uid, gid),
		"aider-local:latest",
		"aider",
		"--model",
		"ollama_chat/qwen2.5-coder:7b",
		"--no-auto-commits",
		"--no-show-release-notes",
	}

	cmd = append(cmd, extend...)

	return cmd, nil
}

func RunAider(ctx context.Context, prompt string, files []string) error {
	if prompt == "" {
		return tracederrors.TracedErrorEmptyString("prompt")
	}

	logging.LogInfoByCtxf(ctx, "Run aider noninteractively started.")

	cmd, err := GetRunCommand(true)
	if err != nil {
		return err
	}
	cmd = append(cmd, "--yes", "--message", prompt)
	if len(files) > 0 {
		cmd = append(cmd, files...)
	}

	_, err = commandexecutorexec.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run aider noninteractively finished.")

	return nil
}
