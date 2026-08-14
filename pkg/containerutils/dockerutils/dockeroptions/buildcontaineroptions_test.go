package dockeroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
)

func Test_BuildContainerOptions_Validation(t *testing.T) {
	// Test case 1: Empty image name
	options := &dockeroptions.BuildContainerOptions{
		ImageNameAndTag:   "",
		DockerfileContent: "FROM alpine:latest",
	}
	err := options.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "ImageNameAndTag is required")

	// Test case 2: Neither DockerfilePath nor DockerfileContent set
	options = &dockeroptions.BuildContainerOptions{
		ImageNameAndTag: "test:latest",
	}
	err = options.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Either DockerfilePath or DockerfileContent must be set")

	// Test case 3: Both DockerfilePath and DockerfileContent set
	options = &dockeroptions.BuildContainerOptions{
		ImageNameAndTag:   "test:latest",
		DockerfilePath:    "/some/path/Dockerfile",
		DockerfileContent: "FROM alpine:latest",
	}
	err = options.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Only one of DockerfilePath or DockerfileContent can be set")

	// Test case 4: Valid options with DockerfilePath
	options = &dockeroptions.BuildContainerOptions{
		ImageNameAndTag: "test:latest",
		DockerfilePath:  "/some/path/Dockerfile",
	}
	err = options.Validate()
	require.NoError(t, err)

	// Test case 5: Valid options with DockerfileContent
	options = &dockeroptions.BuildContainerOptions{
		ImageNameAndTag:   "test:latest",
		DockerfileContent: "FROM alpine:latest",
	}
	err = options.Validate()
	require.NoError(t, err)
}

func Test_BuildContainerOptions_Setters(t *testing.T) {
	options := dockeroptions.NewBuildContainerOptions()

	// Test SetImageNameAndTag
	err := options.SetImageNameAndTag("test-image:latest")
	require.NoError(t, err)
	require.Equal(t, "test-image:latest", options.ImageNameAndTag)

	// Test SetImageNameAndTag with empty string
	err = options.SetImageNameAndTag("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "imageNameAndTag is empty string")

	// Test SetDockerfilePath
	err = options.SetDockerfilePath("/path/to/Dockerfile")
	require.NoError(t, err)
	require.Equal(t, "/path/to/Dockerfile", options.DockerfilePath)

	// Test SetDockerfilePath with empty string
	err = options.SetDockerfilePath("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "dockerfilePath is empty string")

	// Test SetDockerfileContent
	err = options.SetDockerfileContent("FROM alpine:latest")
	require.NoError(t, err)
	require.Equal(t, "FROM alpine:latest", options.DockerfileContent)

	// Test SetDockerfileContent with empty string
	err = options.SetDockerfileContent("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "dockerfileContent is empty string")

	// Test SetBuildContextPath
	err = options.SetBuildContextPath("/path/to/context")
	require.NoError(t, err)
	require.Equal(t, "/path/to/context", options.BuildContextPath)

	// Test SetAdditionalBuildArgs
	buildArgs := map[string]string{"ARG1": "value1", "ARG2": "value2"}
	err = options.SetAdditionalBuildArgs(buildArgs)
	require.NoError(t, err)
	require.Equal(t, buildArgs, options.AdditionalBuildArgs)

	// Test SetNoCache
	err = options.SetNoCache(true)
	require.NoError(t, err)
	require.True(t, options.NoCache)

	err = options.SetNoCache(false)
	require.NoError(t, err)
	require.False(t, options.NoCache)

	// Test SetPullParentImages
	err = options.SetPullParentImages(true)
	require.NoError(t, err)
	require.True(t, options.PullParentImages)

	err = options.SetPullParentImages(false)
	require.NoError(t, err)
	require.False(t, options.PullParentImages)
}
