package dockerutils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_BuildContainerImage_FromDockerfilePath(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExecutorDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const imageName = "test-build-dockerfile:latest"

				// Clean up any existing image
				err := docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)

				// Create a temporary directory with a simple Dockerfile
				tempDir, err := os.MkdirTemp("", "docker-build-test-")
				require.NoError(t, err)
				defer os.RemoveAll(tempDir)

				// Create a simple alpine-based Dockerfile
				dockerfilePath := filepath.Join(tempDir, "Dockerfile")
				dockerfileContent := `FROM alpine:latest
RUN echo "Building from Dockerfile path"
CMD ["echo", "Hello from built image"]`

				err = nativefiles.WriteString(ctx, dockerfilePath, dockerfileContent)
				require.NoError(t, err)

				// Build the image
				options := &dockeroptions.BuildContainerOptions{
					ImageNameAndTag:  imageName,
					DockerfilePath:   dockerfilePath,
					BuildContextPath: tempDir,
				}

				image, err := docker.BuildImage(ctx, options)

				require.NoError(t, err)
				require.NotNil(t, image)

				// Verify image exists
				name, err := image.GetName()
				require.NoError(t, err)
				require.Equal(t, imageName, name)

				exists, err := docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.True(t, exists)

				// Clean up
				err = docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)
			},
		)
	}
}

func Test_BuildContainerImage_FromDockerfileContent(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExecutorDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const imageName = "test-build-content:latest"

				// Clean up any existing image
				err := docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)

				// Create a simple alpine-based Dockerfile content
				dockerfileContent := `FROM alpine:latest
RUN echo "Building from Dockerfile content"
CMD ["echo", "Hello from inline Dockerfile"]`

				// Build the image from content
				options := &dockeroptions.BuildContainerOptions{
					ImageNameAndTag:   imageName,
					DockerfileContent: dockerfileContent,
				}

				image, err := docker.BuildImage(ctx, options)

				require.NoError(t, err)
				require.NotNil(t, image)

				// Verify image exists
				name, err := image.GetName()
				require.NoError(t, err)
				require.Equal(t, imageName, name)

				exists, err := docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.True(t, exists)

				// Clean up
				err = docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)
			},
		)
	}
}

func Test_BuildContainerImage_WithBuildArgs(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExecutorDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const imageName = "test-build-args:latest"

				// Clean up any existing image
				err := docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)

				// Create a simple alpine-based Dockerfile that uses build args
				dockerfileContent := `FROM alpine:latest
ARG BUILD_VERSION=1.0
RUN echo "Build version: ${BUILD_VERSION}"
CMD ["echo", "Build args test"]`

				// Build the image with build arguments
				options := &dockeroptions.BuildContainerOptions{
					ImageNameAndTag:     imageName,
					DockerfileContent:   dockerfileContent,
					AdditionalBuildArgs: map[string]string{"BUILD_VERSION": "2.0"},
				}

				image, err := docker.BuildImage(ctx, options)

				require.NoError(t, err)
				require.NotNil(t, image)

				// Verify image exists
				name, err := image.GetName()
				require.NoError(t, err)
				require.Equal(t, imageName, name)

				exists, err := docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.True(t, exists)

				// Clean up
				err = docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)
			},
		)
	}
}

func Test_BuildContainerImage_WithNoCache(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExecutorDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const imageName = "test-build-nocache:latest"

				// Clean up any existing image
				err := docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)

				// Create a simple alpine-based Dockerfile
				dockerfileContent := `FROM alpine:latest
RUN echo "Building with no cache"
CMD ["echo", "No cache test"]`

				// Build the image with no-cache flag
				options := &dockeroptions.BuildContainerOptions{
					ImageNameAndTag:   imageName,
					DockerfileContent: dockerfileContent,
					NoCache:           true,
				}

				image, err := docker.BuildImage(ctx, options)

				require.NoError(t, err)
				require.NotNil(t, image)

				// Verify image exists
				name, err := image.GetName()
				require.NoError(t, err)
				require.Equal(t, imageName, name)

				exists, err := docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.True(t, exists)

				// Clean up
				err = docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
				require.NoError(t, err)
			},
		)
	}
}

func Test_BuildContainerImage_Validation(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExecutorDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				// Test case 1: Empty image name
				options := &dockeroptions.BuildContainerOptions{
					ImageNameAndTag:   "",
					DockerfileContent: "FROM alpine:latest",
				}
				_, err := docker.BuildImage(ctx, options)
				require.Error(t, err)
				require.Contains(t, err.Error(), "ImageNameAndTag is required")

				// Test case 2: Neither DockerfilePath nor DockerfileContent set
				options = &dockeroptions.BuildContainerOptions{
					ImageNameAndTag: "test:latest",
				}
				_, err = docker.BuildImage(ctx, options)
				require.Error(t, err)
				require.Contains(t, err.Error(), "Either DockerfilePath or DockerfileContent must be set")

				// Test case 3: Both DockerfilePath and DockerfileContent set
				options = &dockeroptions.BuildContainerOptions{
					ImageNameAndTag:   "test:latest",
					DockerfilePath:    "/some/path/Dockerfile",
					DockerfileContent: "FROM alpine:latest",
				}
				_, err = docker.BuildImage(ctx, options)
				require.Error(t, err)
				require.Contains(t, err.Error(), "Only one of DockerfilePath or DockerfileContent can be set")

				// Test case 4: Nil options
				_, err = docker.BuildImage(ctx, nil)
				require.Error(t, err)
				require.Contains(t, err.Error(), "nil")
			},
		)
	}
}
