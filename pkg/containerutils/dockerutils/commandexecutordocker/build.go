package commandexecutordocker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// BuildImage builds a Docker image using the docker build CLI command.
func (c *CommandExecutorDocker) BuildImage(ctx context.Context, options *dockeroptions.BuildContainerOptions) (containerinterfaces.Image, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Building Docker image '%s' started.", options.ImageNameAndTag)

	// Prepare build command
	buildCommand := []string{"docker", "build"}

	// Add tag
	buildCommand = append(buildCommand, "-t", options.ImageNameAndTag)

	// Add no-cache flag if set
	if options.NoCache {
		buildCommand = append(buildCommand, "--no-cache")
	}

	// Add pull flag if set
	if options.PullParentImages {
		buildCommand = append(buildCommand, "--pull")
	}

	// Add build args
	for key, value := range options.AdditionalBuildArgs {
		buildCommand = append(buildCommand, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	// Handle Dockerfile source
	var buildContextPath string
	var cleanupFunc func() error

	if options.DockerfilePath != "" {
		// Build from Dockerfile path
		dockerfilePath := options.DockerfilePath
		buildContextPath = options.BuildContextPath

		// If build context path not specified, use the directory containing the Dockerfile
		if buildContextPath == "" {
			buildContextPath = filepath.Dir(dockerfilePath)
		}

		// If Dockerfile is not in the build context root, we need to handle it
		dockerfileInContext := filepath.Join(buildContextPath, "Dockerfile")
		if dockerfilePath != dockerfileInContext {
			// Create a temporary directory and copy Dockerfile there
			tempDir, err := tempfiles.CreateTempDir(ctx)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to create temp directory: %w", err)
			}

			// Copy Dockerfile to temp directory
			content, err := readDockerfileContent(dockerfilePath)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to read Dockerfile: %w", err)
			}

			tempDockerfilePath := filepath.Join(tempDir, "Dockerfile")
			if err := nativefiles.WriteString(ctx, tempDockerfilePath, content); err != nil {
				return nil, tracederrors.TracedErrorf("Failed to write Dockerfile to temp: %w", err)
			}

			buildContextPath = tempDir
			cleanupFunc = func() error {
				return nativefiles.Delete(ctx, tempDir, nil)
			}
		}

		buildCommand = append(buildCommand, "-f", dockerfilePath)
	} else if options.DockerfileContent != "" {
		// Build from Dockerfile content
		tempDir, err := tempfiles.CreateTempDir(ctx)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to create temp directory: %w", err)
		}

		// Write Dockerfile content
		dockerfilePath := filepath.Join(tempDir, "Dockerfile")
		if err := nativefiles.WriteString(ctx, dockerfilePath, options.DockerfileContent); err != nil {
			return nil, tracederrors.TracedErrorf("Failed to write Dockerfile: %w", err)
		}

		buildContextPath = tempDir
		buildCommand = append(buildCommand, "-f", dockerfilePath)

		cleanupFunc = func() error {
			return nativefiles.Delete(ctx, tempDir, nil)
		}
	}

	// Add build context path
	if buildContextPath == "" {
		buildContextPath = "."
	}
	buildCommand = append(buildCommand, buildContextPath)

	logging.LogInfoByCtxf(ctx, "Running docker build command: %v", buildCommand)

	// Execute build command
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: buildCommand,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to build image '%s': %w", options.ImageNameAndTag, err)
	}

	// Cleanup temporary directory if created
	if cleanupFunc != nil {
		if cleanupErr := cleanupFunc(); cleanupErr != nil {
			logging.LogErrorByCtxf(ctx, "Failed to cleanup temp directory: %v", cleanupErr)
		}
	}

	logging.LogChangedByCtxf(ctx, "Built Docker image '%s'.", options.ImageNameAndTag)

	return c.GetImageByName(options.ImageNameAndTag)
}

func readDockerfileContent(dockerfilePath string) (string, error) {
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to open Dockerfile: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to read Dockerfile: %w", err)
	}

	return string(content), nil
}
