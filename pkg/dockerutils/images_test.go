package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_docker_images(t *testing.T) {
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

				const imageName = "ubuntu:23.04"

				docker := getDockerImplementationByName(tt.implementationName)
				
				err := docker.RemoveImage(ctx, imageName)
				require.NoError(t, err)
				exists, err := docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.False(t, exists)

				image, err := docker.PullImage(ctx, imageName)
				require.NoError(t, err)
				require.NotNil(t, image)

				name, err := image.GetName()
				require.NoError(t, err)
				require.EqualValues(t, imageName, name)
				
				exists, err = docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.True(t, exists)
				
				exists, err = image.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				err = docker.RemoveImage(ctx, imageName)
				require.NoError(t, err)

				exists, err = docker.ImageExists(ctx, imageName)
				require.NoError(t, err)
				require.False(t, exists)
				
				exists, err = image.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
			},
		)
	}
}
